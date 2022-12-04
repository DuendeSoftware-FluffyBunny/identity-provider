package main

import (
	"log"

	"github.com/voice0726/identity-provider/src/config"
	"github.com/voice0726/identity-provider/src/server"
	"go.uber.org/zap"
)

func main() {
	lg, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	conf := config.NewConfig()
	s := server.NewServer(conf, lg)
	s.Register()
	s.Start()
}
