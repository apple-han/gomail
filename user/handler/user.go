package handler

import (
	"context"
	"fmt"

	c "whisper/common"
	email "whisper/email/proto/email"
	e "whisper/pkg/error"
	"whisper/pkg/logging"
	user "whisper/user/proto/user"

	db "whisper/user/db"

	hystrix_go "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
)
// 1秒钟10次 redis 计数功能很强大  0-10  2  key expire 1秒  这么做有什么不好的地方300 300 - 900 1000 

// User 用户结构体
type User struct{}

var emailClient email.EmailService

// Init handle的初始化
func Init() {
	hystrix_go.DefaultVolumeThreshold = 5
	hystrix_go.DefaultErrorPercentThreshold = 20
	cl := hystrix.NewClientWrapper()(client.DefaultClient)
	emailClient = email.NewEmailService("go.micro.srv.email", cl)
}

// Create 创建用户信息
func (u *User) Create(ctx context.Context, req *user.CreateRequest, rsp *user.CreateResult) error {
	// id uuid
	id, err := c.GetUUID()
	logging.Info("id value is ", id)
	if err != nil {
		logging.Error("user.create", err.Error())
		return errors.InternalServerError("user.create", c.ErrorDetail(e.ERROR, e.GetMsg(e.ERROR)))
	}
	// 参数验证
	if err := req.Validate(); err != nil {
		logging.Error("user.create", err.Error())
		return errors.BadRequest("user.Create", c.ErrorDetail(
			e.INVALID_PARAMS, err.Error()))
	}
	req.User.Id = id
	// 用户的info, 存到数据库
	if err := db.Create(req); err != nil {
		if err := SendMail("user.create.SendMail", err.Error(), emailClient); err != nil {
			logging.Error("user.create.SendMail", err.Error())
		}
		logging.Error("user.create", err.Error())
		return errors.InternalServerError("user.create", c.ErrorDetail(e.ERROR, e.GetMsg(e.ERROR)))
	}
	rsp.Id = id
	return nil
}

// Delete 删除一个用户
func (u *User) Delete(ctx context.Context, req *user.DeleteRequest, rsp *user.DeleteResult) error {
	// 参数验证
	if err := req.Validate(); err != nil {
		logging.Error("user.Delete", err.Error())
		return errors.BadRequest("user.Delete", c.ErrorDetail(
			e.INVALID_PARAMS, err.Error()))
	}
	if err := db.Delete(req.Id); err != nil {
		logging.Error("user.Delete", err.Error())
		return errors.InternalServerError("user.Delete", c.ErrorDetail(e.ERROR, e.GetMsg(e.ERROR)))
	}
	return nil
}

// Update 更新一个用户
func (u *User) Update(ctx context.Context, req *user.UpdateRequest, rsp *user.UpdateResult) error {
	// 参数验证
	if err := req.Validate(); err != nil {
		logging.Error("user.Update", err.Error())
		return errors.BadRequest("user.Update", c.ErrorDetail(
			e.INVALID_PARAMS, err.Error()))
	}
	if err := db.Update(req); err != nil {
		logging.Error("user.Update", err.Error())
		return errors.InternalServerError("user.Update", c.ErrorDetail(e.ERROR, e.GetMsg(e.ERROR)))
	}
	return nil
}

// Read 读取一个用户
func (u *User) Read(ctx context.Context, req *user.ReadRequest, rsp *user.ReadResult) error {
	fmt.Println("eeeee")
	if err := req.Validate(); err != nil {
		logging.Error("user.Read", err.Error())
		return errors.BadRequest("user.Read", c.ErrorDetail(
			e.INVALID_PARAMS, err.Error()))
	}
	info, err := db.Read(req.Id)
	if err != nil {
		logging.Error("user.Read", err.Error())
		return errors.InternalServerError("user.Read", c.ErrorDetail(e.ERROR, e.GetMsg(e.ERROR)))
	}
	rsp.User = info
	return nil
}

// SendMail send email common  function
func SendMail(id, detail string, e email.EmailService) error {
	_, err := e.Error(context.Background(), &email.ErrorRequest{
		Id:     id,
		Detail: detail,
	})
	if err != nil {
		return err
	}
	return nil
}
