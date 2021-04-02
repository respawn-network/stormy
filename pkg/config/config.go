// Package config provides a Config containing all information needed to run
// bot and means to load it.
package config

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type (
	// Config is the configuration of the bot.
	Config struct {
		// Token contains the bot token including the `Bot ` prefix.
		Token string

		// Status is the status of the bot.
		//
		// Default: online
		Status gateway.Status
		// ActivityType is the ActivityType of the bot.
		//
		// Default: Playing
		ActivityType discord.ActivityType
		// ActivityName is the name of the activity.
		// If this is empty, no activity will be displayed.
		ActivityName string

		// DateFormat is the format used to encode dates.
		//
		// Default: January 2, 2006
		DateFormat string
		// TimeFormat is the format used to encode times.
		//
		// Default: 3:04 PM
		TimeFormat string
		// Location is the time zone used to generate time stamps.
		//
		// Defaults to the system's local time zone.
		// See time.Local for more
		Location *time.Location

		// ChannelConfigs contains the ChannelConfigs that are being watched.
		ChannelConfigs []ChannelConfig
	}

	// ChannelConfig is the configuration for a specific channel
	ChannelConfig struct {
		// ChannelID is the id of the channel that is being watched.
		ChannelID discord.ChannelID
		// AutoReactions contains the emojis that will be added to every
		// message.
		AutoReactions []discord.APIEmoji
		// ScanReactions contains the emojis that will be added, if found in
		// the message.
		ScanReactions []discord.APIEmoji
		// RepostReactions contains the reactions that can be used to repost
		// the message.
		RepostReactions []RepostReaction
	}

	// RepostReaction is the configuration of a reaction, that can be used to
	// repost the message.
	RepostReaction struct {
		// Target is the id of the channel the message will be reposted to.
		Target discord.ChannelID
		// Reaction is the reaction that is being watched.
		Reaction discord.APIEmoji
		// AutoReact defines if the Reaction should be added by the bot.
		//
		// Default: false
		AutoReact bool
		// AutoDelete defines if the Reaction should be deleted upon crosspost.
		//
		// Default: false
		AutoDelete bool
		// DeleteOnRepost specifies whether the original message will be
		// deleted, if crossposted.
		//
		// Default: false
		DeleteOnRepost bool
		// Message is the message sent as formatted by template.Template.
		//
		// Available variables are
		// 		Message - the original message
		//		MessageQuoted - the original message, but quoted
		//		Author - the name of the author without descriptor
		//		AuthorMention - a mention of the author
		//		Crossposter - the name of the user who authorized the crosspost
		//		CrossposterMention - a mention of the user who authorized the crosspost
		// 		SourceChannel - a mention of the original channel
		//		Time - the time the original message was sent
		//		Date - the date the original message was sent
		//
		// The delimiters are {{ and }} respectively.
		//
		// If the resulting message is longer than 2000 characters (Discord's
		// maximum message length, the bot will attempt to split the message
		// at a line break.
		// If that is not possible either the split will be done at the last
		// word before 2000 characters.
		Message string
		// Rights defines the users and roles that are allowed to authorize a
		// crosspost.
		Rights RepostReactionRights
	}

	// RepostReactionRights defines the rights needed to authorize a crosspost.
	//
	// Note that administrators will always be allowed to crosspost, even
	// if not mentioned in the RepostReactionRights.
	RepostReactionRights struct {
		// UserIDs is a list of users that are allowed to crosspost, no matter
		// their roles.
		UserIDs []discord.UserID
		// RoleIDs is a list of roles of users that are allowed to crosspost.
		RoleIDs []discord.RoleID
	}
)
