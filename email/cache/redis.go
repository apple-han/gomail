package cache

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	c "whisper/pkg/configuration_center/client"
	"whisper/pkg/logging"

	"github.com/garyburd/redigo/redis"
)

var (
	// RedisHost redis host
	RedisHost = "127.0.0.1:6379"
	// RedisDB redis db selected
	RedisDB = 0
	// RedisAuth 远程密码
	RedisAuth = ""
	// RedisPool redis连接池
	RedisPool *redis.Pool
)

// Redis redis 配置
type Redis struct {
	RedisHost string `json:"redisHost"`
	RedisDB   int    `json:"redisDB"`
	RedisAuth string `json:"redisAuth"`
}

// CreateRedisPool 创建redis连接池
func CreateRedisPool(address string, db int, password string) (interface{}, error) {
	RedisPool = &redis.Pool{
		MaxIdle:     20,
		MaxActive:   10,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				fmt.Println("err value is .....", err.Error())
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	return nil, nil
}

// Set a key/value
func Set(key string, data interface{}, time interface{}) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()

	if reflect.TypeOf(data).Kind().String() == "struct" {
		data, err = json.Marshal(data)
		if err != nil {
			logging.Fatal("pkg.redis.set", err.Error())
			return err
		}
	}
	_, err = conn.Do("SET", key, data)
	if err != nil {
		logging.Fatal("pkg.redis.set", err.Error())
		return err
	}
	switch value := time.(type) {
	case int:
		_, err = conn.Do("EXPIRE", key, value)
		if err != nil {
			logging.Fatal("pkg.redis.set", err.Error())
			return err
		}
	}
	return nil
}

// Get get a key
func Get(key string) (interface{}, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	reply, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Init 初始化
func Init() (err error) {
	// 创建redis连接池
	var result = c.C()
	var cfg = &Redis{}
	var _ = result.App("redis", cfg)
	fmt.Println("cfg.RedisHost--->", cfg.RedisHost)
	c, err := CreateRedisPool(cfg.RedisHost, cfg.RedisDB, cfg.RedisAuth)
	if err != nil {
		fmt.Println("err value is", err)
	}
	fmt.Println(c)
	return nil
}
