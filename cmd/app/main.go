package main

import (
	"simple_brocker/internal/server"
)

func main() {
	// cfg := config.GetConfig()
	// t := thread.New(cfg)
	// ioChan := most.New(cfg)
	// ioChan.MakeMost()
	// t.AddThread()
	// t.TRun(ioChan)

	s := server.New()
	s.Run()
}
