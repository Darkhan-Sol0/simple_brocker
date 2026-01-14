package main

import (
	"simple_brocker/internal/config"
	"simple_brocker/internal/service/thread"
)

func main() {
	t := thread.New(config.GetConfig())
	t.AddThread()
	t.TRun()

	// cfg := config.GetConfig()

	// fmt.Println(cfg)
}
