package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/mavolin/disstate/v3/pkg/state"
	"go.uber.org/zap"

	"github.com/mavolin/stormy/pkg/bot"
	"github.com/mavolin/stormy/pkg/config"
)

var configPathFlag = flag.String("c", "", "specify a custom config path")
var debugFlag = flag.Bool("debug", false, "run the bot in debug mode")

func init() {
	flag.Parse()
}

func main() {
	err := initLogger()
	if err != nil {
		panic(err)
	}

	c, err := loadConfig()
	if err != nil {
		panic(err)
	}

	err = startBot(c)
	if err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)

	<-s
}

func startBot(c *config.Config) error {
	b, err := bot.NewBot(c)
	if err != nil {
		return err
	}

	zap.S().Info("starting bot")

	b.State.MustAddHandlerOnce(func(_ *state.State, e *state.ReadyEvent) {
		zap.S().Infof("serving as %s on %d servers", e.User.Username, len(e.Guilds))
	})

	err = b.Open()
	if err != nil {
		return err
	}

	return nil
}

func initLogger() (err error) {
	var l *zap.Logger

	if *debugFlag {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}

	if err != nil {
		return
	}

	zap.ReplaceGlobals(l)
	return
}

func loadConfig() (*config.Config, error) {
	if *configPathFlag != "" {
		zap.S().Infow("reading config", "path", *configPathFlag)

		return config.LoadFromPath(*configPathFlag)
	}

	zap.S().Info("reading config")

	c, err := config.Load()
	zap.S().Debugw("config read", "config", c)

	return c, err
}
