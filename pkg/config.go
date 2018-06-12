package ircDiscordRelay

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

var Config *Settings

type Discord struct {
	ChannelId *string `json:"´channelId,omitempty"`
	Token     *string `json:"´token,omitempty"`
}

type Irc struct {
	Channel *string `json:"channel,omitempty"`
	Nick    *string `json:"nick,omitempty"`
	Server  *string `json:"server,omitempty"`
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
	config := Settings{}
	err = json.Unmarshal(content, &config)
	if err != nil {
		return err
	}
	Config = &config
	return checkConfig(&config)
}

func checkConfig(config *Settings) error {
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
	case config.Discord == nil:
		return errors.New("discord is not set.")
	case config.Discord.ChannelId == nil:
		return errors.New("discord.channelId is not set.")
	case config.Discord.Token == nil:
		return errors.New("discord.token is not set.")
	default:
		return nil
	}
}
