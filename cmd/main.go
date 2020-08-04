package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/mavolin/stormy/pkg/bot"
	"github.com/mavolin/stormy/pkg/config"
)

var configPath = flag.String("c", "", "specify a custom config path")

func init() {
	flag.Parse()
}

func main() {
	var c *config.Config
	var err error

	if *configPath != "" {
		c, err = config.LoadFromPath(*configPath)
	} else {
		c, err = config.Load()
	}

	if err != nil {
		panic(err)
	}

	b, err := bot.NewBot(c)
	if err != nil {
		panic(err)
	}

	err = b.Open()
	if err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)

	<-s
}
