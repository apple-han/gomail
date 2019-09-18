package main

import (
	"whisper/user/db"
	"whisper/user/handler"

	"whisper/pkg/logging"

	c "whisper/common"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	user "whisper/user/proto/user"
)

func main() {
	c.InitCfg()

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) {
			// handel 初始化
			handler.Init()
			// db init
			if err := db.Init(); err != nil {
				logging.Fatal("user.db", err.Error())
				return
			}
			// logging init
			logging.Init()
		}),
	)

	// Register Handler
	user.RegisterUserHandler(service.Server(), new(handler.User))
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal("user.main", err.Error())
	}
}
