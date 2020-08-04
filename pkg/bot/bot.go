package bot

import (
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"go.uber.org/zap"

	"github.com/mavolin/stormy/pkg/config"
)

// Bot is an instance of the stormy bot.
type Bot struct {
	// Config is the config of the bot.
	Config config.Config

	// State is the state.State the bot uses.
	State *state.State
}

// NewBot creates a new bot and registers all handler.s
func NewBot(c config.Config) (*Bot, error) {
	s, err := state.NewWithIntents(c.Token, gateway.IntentGuildMessages, gateway.IntentGuildMessageReactions)
	if err != nil {
		return nil, err
	}

	s.ErrorHandler = errorHandler
	s.PanicHandler = panicHandler

	b := &Bot{
		Config: c,
		State:  s,
	}

	s.AutoAddHandlers(b) // add our handlers automatically

	return b, err
}

func errorHandler(err error) {
	zap.S().Error(err)
}

func panicHandler(err interface{}) {
	zap.S().Errorw("panic recovery", "err", err)
}
