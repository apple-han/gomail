package db

var (
	schemas = []string{
		// 用户表
		`CREATE TABLE IF NOT EXISTS info (
			id VARCHAR(36) NOT NULL,
			username VARCHAR(128) NOT NULL COMMENT '用户名',
			secret VARCHAR(128) NOT NULL COMMENT '密码/密钥',
			phone VARCHAR(45) DEFAULT '' COMMENT '电话号码',
			role_name VARCHAR(45) COMMENT '角色的名称',
			created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY phone_UK (phone)
		);`,
	}
)