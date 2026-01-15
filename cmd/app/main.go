package main

import (
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/most"
	"simple_brocker/internal/service/thread"
)

func main() {
	cfg := config.GetConfig()

	// fmt.Println(cfg)
	// fmt.Println(cfg.GetService("group1"))

	ioChan := most.New(cfg)
	t := thread.New(cfg, ioChan)
	t.TRun(ioChan)

	// s := server.New()
	// s.Run()
}
