package db

import (
	"fmt"
	"time"

	c "whisper/pkg/configuration_center/client"

	user "whisper/user/proto/user"

	_ "github.com/go-sql-driver/mysql" // mysql drive
	"github.com/jmoiron/sqlx"
)

var (
	// 数据库连接
	db *sqlx.DB
)

// MysqlInfo mysqlInfo配置
type MysqlInfo struct {
	User Mysql `json: "user"`
}

// Mysql mysql配置结构体
type Mysql struct {
	Enabled           bool   `json: "enabled"`
	Url               string `json: "url"`
	MaxIdleConnection int    `json: "maxIdleConnection"`
	MaxOpenConnection int    `json: "maxOpenConnection"`
}

// Create 创建用户信息
func Create(user *user.CreateRequest) error {
	if _, err := db.NamedExec(`INSERT INTO info (id, username, phone, secret, role_name)
		VALUES (:id, :username, :phone, :secret, :role_name)`, user.User); err != nil {
		return err
	}
	return nil
}

// Read 读取用户信息
func Read(id string) (*user.Info, error) {
	info := &user.Info{}
	if err := db.Unsafe().Get(info, `SELECT IFNULL(role_name, " ") AS role_name FROM info WHERE id = ?`, id); err != nil {
		return nil, err
	}
	return info, nil
}

// Delete 删除一个用户
func Delete(id string) error {
	if _, err := db.Exec(`DELETE FROM info WHERE id = ?`, id); err != nil {
		return err
	}
	return nil
}

// Update 更新用户的名称
func Update(user *user.UpdateRequest) error {
	_, err := db.Exec(`UPDATE info SET username = ? WHERE id = ?`, user.Username, user.Id)
	if err != nil {
		return err
	}
	return err
}

// mapper roleName => role_name
func mapper(name string) string {
	var s []byte
	for i, r := range []byte(name) {
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
			if i != 0 {
				s = append(s, '_')
			}
		}
		s = append(s, r)
	}
	return string(s)
}

// Connect 初始化数据库连接池
func Connect(driver, uri string, maxOpen, maxIdel int) (db *sqlx.DB, err error) {
	db, err = sqlx.Connect(driver, uri+"?charset=utf8mb4&parseTime=true")
	if err != nil {
		return
	}
	// 配置连接池
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdel)
	db.SetConnMaxLifetime(2 * time.Hour)
	return
}

// Init 初始化
func Init() error {
	c := c.C()
	cfg := &MysqlInfo{}
	if err := c.App("db", cfg); err != nil {
		return err
	}
	fmt.Println("cfg.User.MaxOpenConnection---", cfg.User.MaxOpenConnection)

	var err error
	db, err = Connect("mysql", cfg.User.Url, cfg.User.MaxOpenConnection, cfg.User.MaxIdleConnection)
	if err != nil {
		return err
	}
	db.MapperFunc(mapper)
	sqlx.NameMapper = mapper

	for _, s := range schemas {
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}
