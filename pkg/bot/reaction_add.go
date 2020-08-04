package bot

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/dasync/pkg/dasync"
	"github.com/mavolin/disstate/pkg/state"
	"go.uber.org/multierr"

	"github.com/mavolin/stormy/pkg/config"
)

const maxMsgLen = 2000

func (b *Bot) HandleReactionAdd(s *state.State, e *state.MessageReactionAddEvent) error {
	if e.GuildID == 0 || e.Member.User.Bot {
		return nil
	}

	for _, c := range b.Config.ChannelConfigs {
		if c.ChannelID == e.ChannelID {
			return b.handleReaction(s, e, c)
		}
	}

	return nil
}

func (b *Bot) handleReaction(s *state.State, e *state.MessageReactionAddEvent, c config.ChannelConfig) error {
	for _, r := range c.RepostReactions {
		if e.Emoji.APIString() == r.Reaction {
			return b.handleCrosspost(s, e, r)
		}
	}

	return nil
}

func (b *Bot) handleCrosspost(s *state.State, e *state.MessageReactionAddEvent, r config.RepostReaction) error {
	for _, id := range r.Rights.UserIDs {
		if id == e.UserID {
			return b.crosspost(s, e, r)
		}
	}

	for _, searchID := range r.Rights.RoleIDs {
		for _, id := range e.Member.RoleIDs {
			if searchID == id {
				return b.crosspost(s, e, r)
			}
		}
	}

	// the user wasn't explicitly granted permission, but maybe he is an admin

	gf := dasync.Guild(s, e.GuildID)
	cf := dasync.Channel(s, e.ChannelID)

	g, err := gf()
	if err != nil {
		return err
	}

	c, err := cf()
	if err != nil {
		return err
	}

	perms := discord.CalcOverwrites(*g, *c, *e.Member)
	if perms.Has(discord.PermissionAdministrator) {
		return b.crosspost(s, e, r)
	}

	return nil
}

type templateFields struct {
	Message            string
	MessageQuoted      string
	Author             string
	AuthorMention      string
	Crossposter        string
	CrossposterMention string
	SourceChannel      string
	Time               string
	Date               string
}

func (b *Bot) crosspost(s *state.State, e *state.MessageReactionAddEvent, r config.RepostReaction) (err error) {
	msgf := dasync.Message(s, e.ChannelID, e.MessageID)
	memf := dasync.Member(s, e.GuildID, e.UserID)
	rf := dasync.DeleteReactions(s, e.ChannelID, e.MessageID, e.Emoji.APIString())

	msg, err := msgf()
	if err != nil {
		return err
	}

	mem, err := memf()
	if err != nil {
		return err
	}

	authorName := mem.Nick
	if authorName == "" {
		authorName = msg.Author.Username
	}

	crossposterName := e.Member.Nick
	if crossposterName == "" {
		crossposterName = e.Member.User.Username
	}

	f := templateFields{
		Message:            msg.Content,
		MessageQuoted:      strings.ReplaceAll(msg.Content, "\n", "\n> "),
		Author:             authorName,
		AuthorMention:      mem.Mention(),
		Crossposter:        crossposterName,
		CrossposterMention: e.Member.Mention(),
		SourceChannel:      fmt.Sprintf("<#%d>", msg.ChannelID),
		Time:               msg.Timestamp.Format(b.Config.TimeFormat),
		Date:               msg.Timestamp.Format(b.Config.DateFormat),
	}

	t, err := template.New("crosspost").Parse(r.Message)
	if err != nil {
		return err
	}

	builder := new(strings.Builder)
	builder.Grow(maxMsgLen) // expecting msg max, although the message could be longer

	err = t.Execute(builder, f)
	if err != nil {
		return err
	}

	post := builder.String()

	if len(post) <= maxMsgLen {
		_, err := s.SendText(r.Target, post)
		if err != nil {
			return err
		}

		return rf()
	}

	msgs := splitMessage(post)

	errCallbacks := make([]func() (*discord.Message, error), len(msgs))

	for i, m := range msgs {
		errCallbacks[i] = dasync.SendText(s, r.Target, m)

		time.Sleep(1 * time.Millisecond) // make sure we keep the correct order
	}

	for _, c := range errCallbacks {
		_, err2 := c()
		err = multierr.Append(err, err2)
	}

	err = multierr.Append(err, rf())

	return
}

func splitMessage(msg string) (msgs []string) {
	words := strings.Fields(strings.TrimSpace(msg))
	if len(words) == 0 {
		return hardSplit(msg)
	}

	var current string
	spaceLeft := maxMsgLen
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			msgs = append(msgs, current)

			current = word
			spaceLeft = maxMsgLen - len(word)
		} else {
			current += " " + word
			spaceLeft -= 1 + len(word)
		}
	}

	return
}

func hardSplit(msg string) (msgs []string) {
	msgs = make([]string, len(msg)/maxMsgLen)

	for len(msg) > maxMsgLen {
		msgs = append(msgs, msg[:maxMsgLen])
		msg = msg[maxMsgLen:]
	}

	return
}
