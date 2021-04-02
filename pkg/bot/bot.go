package bot

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/mavolin/disstate/v3/pkg/state"
	"go.uber.org/zap"

	"github.com/mavolin/stormy/pkg/config"
)

// Bot is an instance of the stormy bot.
type Bot struct {
	// Config is the config of the bot.
	Config *config.Config

	// State is the state.State the bot uses.
	State *state.State
}

// NewBot creates a new bot and registers all handler.s
func NewBot(c *config.Config) (*Bot, error) {
	s, err := state.NewWithIntents("Bot "+c.Token, gateway.IntentGuildMessages, gateway.IntentGuildMessageReactions)
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

func (b *Bot) Open() error {
	if err := b.State.Open(); err != nil {
		return err
	}

	return b.State.Gateway.UpdateStatus(gateway.UpdateStatusData{
		Activities: []discord.Activity{
			{
				Name: b.Config.ActivityName,
				Type: b.Config.ActivityType,
			},
		},
		Status: b.Config.Status,
	})
}

func errorHandler(err error) {
	zap.S().Error(err)
}

func panicHandler(err interface{}) {
	zap.S().Errorw("panic recovery", "err", err)
}
