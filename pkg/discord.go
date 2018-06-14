package ircDiscordRelay

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func StartDiscord() error {
	session, err := discordgo.New("Bot " + *Config.Discord.Token)
	if err != nil {
		return err
	}

	session.AddHandler(func(session *discordgo.Session, msg *discordgo.Ready) { session.UpdateStatus(0, *Config.Irc.Channel+" relay") })
	valid := false
	for _, value := range *Config.Discord.Sharing {
		switch value {
		case "message":
			valid = true
			session.AddHandler(onDiscordMsg)
		default:
			fmt.Println("Invalid discord.sharing value '" + value + "' will be ignored.")
		}
	}
	if !valid {
		return errors.New("No valid values in discord.sharing.")
	}

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

func onDiscordMsg(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.Bot || !Relay.isReady() || msg.ChannelID != *Config.Discord.ChannelId { // ignore message from bots (including myself) and if not ready
		return
	}
	msgText, err := msg.ContentWithMoreMentionsReplaced(session);
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	SendIrc("<" + msg.Author.Username + "> " + msgText)
	// if message contains an attachment
	if msg.Attachments != nil && len(msg.Attachments) > 0 {
		for _, att := range msg.Attachments {
			SendIrc("<" + msg.Author.Username + "> " + att.URL)
		}
	}
}
