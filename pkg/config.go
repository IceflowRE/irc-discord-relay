// To know more about the configuration see the Readme

package idrelay

import (
	"encoding/json"
	"errors"
	"os"
)

var config *settings

type discord struct {
	ChannelID *string   `json:"channelId,omitempty"`
	Token     *string   `json:"token,omitempty"`
	Sharing   *[]string `json:"sharing,omitempty"`
}

type irc struct {
	Channel      *string   `json:"channel,omitempty"`
	Nick         *string   `json:"nick,omitempty"`
	Server       *string   `json:"server,omitempty"`
	OnConnection *[]string `json:"onConnection,omitempty"`
	Sharing      *[]string `json:"sharing,omitempty"`
}

type settings struct {
	Irc     *irc     `json:"irc,omitempty"`
	Discord *discord `json:"discord,omitempty"`
}

// LoadConfig loads the configuration file from the given path and saves it to config
func LoadConfig(file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	if config.Irc.OnConnection == nil {
		config.Irc.OnConnection = &[]string{}
	}
	if config.Irc.Sharing == nil {
		config.Irc.Sharing = &[]string{"message", "me", "join", "leaving", "quit", "kick", "mode"}
	}
	if config.Discord.Sharing == nil {
		config.Discord.Sharing = &[]string{"message"}
	}

	return config.checkConfig()
}

// sharing values will be checked at the starting routines
func (config *settings) checkConfig() error {
	switch {
	case config == nil:
		return errors.New("settings is nil")
	case config.Irc == nil:
		return errors.New("irc is not set")
	case config.Irc.Channel == nil:
		return errors.New("irc.channel is not set")
	case config.Irc.Nick == nil:
		return errors.New("irc.nick is not set")
	case config.Irc.Server == nil:
		return errors.New("irc.server is not set")
	case len(*config.Irc.Sharing) == 0:
		return errors.New("irc.sharing is empty")
	case config.Discord == nil:
		return errors.New("discord is not set")
	case config.Discord.ChannelID == nil:
		return errors.New("discord.channelId is not set")
	case config.Discord.Token == nil:
		return errors.New("discord.token is not set")
	case len(*config.Discord.Sharing) == 0:
		return errors.New("discord.sharing is empty")
	}

	return nil
}
