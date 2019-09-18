package common

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	c "whisper/pkg/configuration_center/client"

	"github.com/micro/go-plugins/config/source/grpc"
	uuid "github.com/satori/go.uuid"
)

// ErrorDetail 错误的细节
func ErrorDetail(code int, msg string) string {
	b, _ := json.Marshal(map[string]interface{}{
		"errorcode": code,
		"msg":       msg,
	})
	return string(b)
}

// GetUUID 获取uuid
func GetUUID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// InitCfg 初始化配置文件
func InitCfg() {
	source := grpc.NewSource(
		grpc.WithAddress("192.168.1.107:9600"),
		grpc.WithPath("whisper"),
	)
	c.Init(c.WithSource(source))
}

// VerifyEmailFormat 验证电子邮件的格式
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// Captcha 验证码
func Captcha() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}
