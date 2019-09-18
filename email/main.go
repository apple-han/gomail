package main

import (
	"whisper/email/handler"
	"whisper/pkg/logging"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	c "whisper/common"
	cache "whisper/email/cache"
	email "whisper/email/proto/email"
)

func main() {
	c.InitCfg()
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.email"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) {
			// logging init
			logging.Init()
		}),
	)
	// Register Handler
	email.RegisterEmailHandler(service.Server(), new(handler.Email))

	// redis init
	cache.Init()
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
