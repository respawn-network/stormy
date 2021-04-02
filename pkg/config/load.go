package config

import (
	"fmt"
	"os"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/spf13/viper"
)

// these types take in the raw config data, before all fields get properly
// marshaled into a real Config.
type (
	loadableConfig struct {
		Token string

		Status       gateway.Status
		ActivityType string
		ActivityName string

		DateFormat string
		TimeFormat string
		Location   string

		ChannelConfigs []ChannelConfig
	}
)

// Load attempts to load the config from the working directory or the home
// directory.
func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("stormy")
	v.AddConfigPath(".")
	if home := os.Getenv("HOME"); home != "" {
		v.AddConfigPath(home)
	}

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	setupDefaults(v)

	lc, err := load(v)
	if err != nil {
		return nil, err
	}

	return loadableToConfig(lc)
}

// LoadFromPath loads the config file from the specified path.
func LoadFromPath(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	setupDefaults(v)

	lc, err := load(v)
	if err != nil {
		return nil, err
	}

	return loadableToConfig(lc)
}

func setupDefaults(v *viper.Viper) {
	v.SetDefault("status", gateway.OnlineStatus)
	v.SetDefault("activityType", "playing")
	v.SetDefault("dateFormat", "January 2, 2006")
	v.SetDefault("timeFormat", "3:04 PM")
	v.SetDefault("location", time.Local.String())
}

func load(v *viper.Viper) (c *loadableConfig, err error) {
	err = v.ReadInConfig()
	if err != nil {
		return
	}

	err = v.Unmarshal(&c)

	return c, err
}

func loadableToConfig(lc *loadableConfig) (c *Config, err error) {
	c = &Config{
		Token:          lc.Token,
		ActivityName:   lc.ActivityName,
		DateFormat:     lc.DateFormat,
		TimeFormat:     lc.TimeFormat,
		ChannelConfigs: lc.ChannelConfigs,
	}

	err = isValidStatus(lc.Status)
	if err != nil {
		return
	}

	c.ActivityType, err = parseActivityType(lc.ActivityType)
	if err != nil {
		return
	}

	for i, conf := range c.ChannelConfigs {
		for j, r := range conf.RepostReactions {
			if r.Message == "" {
				c.ChannelConfigs[i].RepostReactions[j].Message = "{{.Message}}"
			}
		}
	}

	c.Location, err = time.LoadLocation(lc.Location)
	return
}

func isValidStatus(s gateway.Status) error {
	if s == gateway.OnlineStatus || s == gateway.DoNotDisturbStatus || s == gateway.IdleStatus ||
		s == gateway.InvisibleStatus {
		return nil
	}

	return fmt.Errorf("%s is not a valid status", s)
}

func parseActivityType(activityType string) (discord.ActivityType, error) {
	switch activityType {
	case "playing":
		return discord.GameActivity, nil
	case "streaming":
		return discord.StreamingActivity, nil
	case "listening":
		return discord.ListeningActivity, nil
	case "watching":
		return discord.WatchingActivity, nil
	default:
		return 0, fmt.Errorf("%s is not a valid activity type", activityType)
	}
}
