package handler

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"whisper/pkg/logging"

	email "whisper/email/proto/email"

	c "whisper/common"
	w "whisper/pkg/error"

	cache "whisper/email/cache"
	cc "whisper/pkg/configuration_center/client"

	"github.com/garyburd/redigo/redis"
	"github.com/go-gomail/gomail"
	"github.com/micro/go-micro/errors"
)

// Email 邮件配置结构体
type Email struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Server    string `json:"server"`
	Port      int    `json:"port"`
	Address   string `json: "address"`
	ToAddress string `json: "toAddress"`
}

// Message 发送文本
func (e *Email) Message(ctx context.Context, req *email.MessageRequest, rsp *email.MessageResult) error {
	if err := req.Validate(); err != nil {
		logging.Error("email.Message", err.Error())
		return errors.BadRequest("email.Message", c.ErrorDetail(
			w.INVALID_PARAMS, err.Error()))
	}
	// 发邮件 发短信  Go 做并发编程的优势
	go func() {
		if err := MailSending(req.Content, "发送文本信息", req.To); err != nil {
			logging.Fatal("email.Message", err.Error())
			return
		}
	}()
	logging.Info("email.Message", "message码发送成功！")
	return nil
}

// Code 发送验证码
func (e *Email) Code(ctx context.Context, req *email.CodeRequest, rsp *email.CodeResult) error {
	if err := req.Validate(); err != nil {
		logging.Error("email.Code", err.Error())
		return errors.BadRequest("email.Code", c.ErrorDetail(
			w.INVALID_PARAMS, err.Error()))
	}
	go func() {
		if err := MailSending(c.Captcha(), "发送验证码", req.To); err != nil {
			logging.Fatal("email.Code", err.Error())
			return
		}
	}()
	logging.Info("email.Code", "code码发送成功！")
	return nil
}

// Error 错误信息的通知
func (e *Email) Error(ctx context.Context, req *email.ErrorRequest, rsp *email.ErrorResult) error {
	// 1. 同一个错误, 在一定得时间内 不应该重复发送
	// 2. 我希望前一天所有的bug,都应该是解决完的状态
	// expire  一天

	if err := DealErrorByRedis(req.Id, req.Detail); err != nil {
		logging.Fatal("email.Message", err.Error())
		return errors.InternalServerError("email.error", c.ErrorDetail(w.ERROR, w.GetMsg(w.ERROR)))
	}
	fmt.Println("到这里了吗")
	logging.Info("email.Error", "error邮件发送成功")
	return nil
}

// DealErrorByRedis  通过redis,报错异常的信息
func DealErrorByRedis(id, detail string) error {
	c := cache.RedisPool.Get()
	defer c.Close()
	key := "my:email:err"
	// 判断key isExist
	value, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return err
	}
	if value == false {
		_, err := c.Do("SADD", key, id+detail)
		if err != nil {
			return err
		}
		c.Do("EXPIRE", key, 3600*24+1) // 一天后过期
	} else {
		num, err := redis.Int(c.Do("SADD", key, id+detail))
		if err != nil {
			return err
		}
		if num == 0 {
			return nil
		}
	}
	// 发邮件
	MailSending(id+detail, "异常信息处理", "")
	return nil
}

// MailSending 发送邮件
func MailSending(content, subject, to string) error {
	// 获取配置参数
	var result = cc.C()
	var cfg = &Email{}
	var _ = result.App("email", cfg)
	if len(to) == 0 {
		fmt.Println(cfg.ToAddress)
		to = cfg.ToAddress
	}
	m := gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", cfg.Address, cfg.Username)
	// 收件人
	m.SetHeader("To", strings.Split(to, ",")[0])
	// 抄送
	m.SetHeader("Cc", strings.Split(to, ",")[1:]...)
	// 主题
	m.SetHeader("Subject", subject)
	// 内容
	m.SetBody("text/html", "<h4>CONTENT: "+content+"<h4>")

	// 发送邮件服务器、端口、发件人账号、发件人密码
	d := gomail.NewDialer(cfg.Server, cfg.Port, cfg.Address, cfg.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
