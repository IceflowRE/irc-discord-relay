package ircDiscordRelay

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

var Config *Settings

type Discord struct {
	ChannelId *string   `json:"channelId,omitempty"`
	Token     *string   `json:"token,omitempty"`
	Sharing   *[]string `json:"sharing,omitempty"`
}

type Irc struct {
	Channel      *string   `json:"channel,omitempty"`
	Nick         *string   `json:"nick,omitempty"`
	Server       *string   `json:"server,omitempty"`
	OnConnection *[]string `json:"onConnection,omitempty"`
	Sharing      *[]string `json:"sharing,omitempty"`
}

type Settings struct {
	Irc     *Irc     `json:"irc,omitempty"`
	Discord *Discord `json:"discord,omitempty"`
}

func LoadConfig(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &Config)
	if err != nil {
		return err
	}

	if Config.Irc.OnConnection == nil {
		Config.Irc.OnConnection = &[]string{}
	}
	if Config.Irc.Sharing == nil {
		Config.Irc.Sharing = &[]string{"message", "me", "join", "leaving", "quit", "kick", "mode"}
	}
	if Config.Discord.Sharing == nil {
		Config.Discord.Sharing = &[]string{"message"}
	}

	return Config.checkConfig()
}

// sharing values will be checked at the starting routines
func (config *Settings) checkConfig() error {
	switch {
	case config == nil:
		return errors.New("Settings is nil.")
	case config.Irc == nil:
		return errors.New("irc is not set.")
	case config.Irc.Channel == nil:
		return errors.New("irc.channel is not set.")
	case config.Irc.Nick == nil:
		return errors.New("irc.nick is not set.")
	case config.Irc.Server == nil:
		return errors.New("irc.server is not set.")
	case len(*config.Irc.Sharing) == 0:
		return errors.New("irc.sharing is empty.")
	case config.Discord == nil:
		return errors.New("discord is not set.")
	case config.Discord.ChannelId == nil:
		return errors.New("discord.channelId is not set.")
	case config.Discord.Token == nil:
		return errors.New("discord.token is not set.")
	case len(*config.Discord.Sharing) == 0:
		return errors.New("discord.sharing is empty.")
	}

	return nil
}
