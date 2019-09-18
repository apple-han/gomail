package go_micro_srv_email

import (
	"errors"
	c "whisper/common"
)

// Validate message的验证层
func (req *MessageRequest) Validate() error {
	if len(req.To) == 0 {
		return errors.New("收件人不能为空")
	}
	if c.VerifyEmailFormat(req.To) == false {
		return errors.New("邮箱格式不正确")
	}
	if len(req.Content) == 0 {
		return errors.New("信息不能为空")
	}

	if len(req.Content) > 280 {
		return errors.New("信息内容不能过长")
	}

	return nil
}

// Validate code的验证层
func (req *CodeRequest) Validate() error {
	if len(req.To) == 0 {
		return errors.New("收件人不能为空")
	}
	if c.VerifyEmailFormat(req.To) == false {
		return errors.New("邮箱格式不正确")
	}
	return nil
}