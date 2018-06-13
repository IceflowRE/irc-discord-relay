package ircDiscordRelay

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func StartDiscord() error {
	session, err := discordgo.New("Bot " + *Config.Discord.Token)
	if err != nil {
		return err
	}
	session.AddHandler(func(session *discordgo.Session, msg *discordgo.Ready) { session.UpdateStatus(0, *Config.Irc.Channel+" relay") })
	session.AddHandler(onDiscord)
	err = session.Open()
	if err != nil {
		return err
	}
	Relay.dSession = session
	return nil
}

func SendDiscord(msg string) {
	_, err := Relay.dSession.ChannelMessageSend(*Config.Discord.ChannelId, msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func onDiscord(_ *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.Bot || !Relay.isReady() || msg.ChannelID != *Config.Discord.ChannelId { // ignore message from bots (including myself) and if not ready
		return
	}
	SendIrc("<" + msg.Author.Username + "> " + msg.Content)
}
