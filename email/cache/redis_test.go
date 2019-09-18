package cache

import (
	"fmt"
	"testing"
	c "whisper/common"

	"github.com/garyburd/redigo/redis"
)

func Test_Init(t *testing.T) {
	c.InitCfg()
	if err := Init(); err != nil {
		t.Errorf("but %s", err)
	}
}

func Test_Set(t *testing.T) {
	if err := Set("whisper:email", "hello", nil); err != nil {
		t.Errorf("but %s", err)
	}
}

func Test_Get(t *testing.T) {
	value, err := redis.String(Get("whisper:email"))

	fmt.Println("hello byte value is", []byte("hello"))
	fmt.Println("value-----", value)
	if err != nil {
		return
	}
	if "hello" != value {
		t.Errorf("excepted %s but %s", "hello", value)
	}
}
