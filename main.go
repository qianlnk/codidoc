package main

import (
	"github.com/qianlnk/codidoc/config"
	"github.com/qianlnk/codidoc/server"
	"github.com/qianlnk/log"
)

func main() {
	cfg, err := config.Load("./app.yaml")
	if err != nil {
		log.Errorf("load config err: %v", err)
		return
	}

	svr, err := server.New(cfg)
	if err != nil {
		log.Errorf("create server err: %v", err)
		return
	}

	svr.Run()

	select {}
}
