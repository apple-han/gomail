package go_micro_srv_user

import "errors"

// Validate Create的验证层
func (req *CreateRequest) Validate() error {
	if len(req.User.Username) == 0 {
		return errors.New("username is must not empty")
	}
	if len(req.User.Secret) == 0 {
		return errors.New("Secret is must not empty")
	}
	if len(req.User.Phone) == 0 {
		return errors.New("phone is must not empty")
	}
	return nil
}

// Validate 验证 DeleteRequest
func (req *DeleteRequest) Validate() error {
	if len(req.Id) == 0 {
		return errors.New("id is must not empty")
	}
	return nil
}

// Validate 验证 ReadRequest
func (req *ReadRequest) Validate() error {
	if len(req.Id) == 0 {
		return errors.New("id is must not empty")
	}
	return nil
}

// Validate 验证 UpdateRequest
func (req *UpdateRequest) Validate() error {
	if len(req.Username) == 0 {
		return errors.New("username is must not empty")
	}
	return nil
}