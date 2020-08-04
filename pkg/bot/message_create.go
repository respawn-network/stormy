package bot

import (
	"strings"

	"github.com/mavolin/dasync/pkg/dasync"
	"github.com/mavolin/disstate/pkg/state"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/mavolin/stormy/pkg/config"
)

func (b *Bot) HandleMessageCreate(s *state.State, e *state.MessageCreateEvent) error {
	for _, c := range b.Config.ChannelConfigs {
		if c.ChannelID == e.ChannelID {
			return b.handlePost(s, e, c)
		}
	}

	return nil
}

func (b *Bot) handlePost(s *state.State, e *state.MessageCreateEvent, c config.ChannelConfig) (err error) {
	zap.S().Infow("new message in watched channel",
		"message_id", e.Message.ID,
		"channel_id", e.ChannelID,
		"guild_id", e.GuildID)

	errCallbacks := make([]func() error, 0, len(c.AutoReactions)+len(c.RepostReactions))

	for _, r := range c.AutoReactions {
		rf := dasync.React(s, e.ChannelID, e.ID, r)
		errCallbacks = append(errCallbacks, rf)
	}

	for _, r := range c.ScanReactions {
		if strings.Contains(e.Content, r) {
			rf := dasync.React(s, e.ChannelID, e.ID, r)
			errCallbacks = append(errCallbacks, rf)
		}
	}

	for _, r := range c.RepostReactions {
		rf := dasync.React(s, e.ChannelID, e.ID, r.Reaction)
		errCallbacks = append(errCallbacks, rf)
	}

	// collect our errors
	for _, c := range errCallbacks {
		err = multierr.Append(err, c())
	}

	return
}
