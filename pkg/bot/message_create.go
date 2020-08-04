package bot

import (
	"strings"

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

	for _, r := range c.AutoReactions {
		err2 := s.React(e.ChannelID, e.ID, r)
		if err2 != nil {
			err = multierr.Append(err, err2)
		}
	}

	for _, r := range c.ScanReactions {
		if strings.Contains(e.Content, r) {
			err2 := s.React(e.ChannelID, e.ID, r)
			if err2 != nil {
				err = multierr.Append(err, err2)
			}
		}
	}

	for _, r := range c.RepostReactions {
		if r.AutoReact {
			err2 := s.React(e.ChannelID, e.ID, r.Reaction)
			if err2 != nil {
				err = multierr.Append(err, err2)
			}
		}
	}

	return
}
