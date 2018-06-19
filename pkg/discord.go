package ircDiscordRelay

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

var discordNickReg = regexp.MustCompile("@[a-zA-Z0-9]*")

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

	chn, err := Relay.dSession.Channel(*Config.Discord.ChannelId)
	if err != nil {
		return err
	}
	Relay.dGuildId = chn.GuildID

	return nil
}

// send message on discord
func SendDiscord(msg string) {
	_, err := Relay.dSession.ChannelMessageSend(*Config.Discord.ChannelId, msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

var emojiRe = regexp.MustCompile("(<:.*:)[0-9]*(>)")
// removes the unique id from the emoji part
func stripEmoji(msg string) string {
	return emojiRe.ReplaceAllString(msg, "$1$2")
}

// on discord message
func onDiscordMsg(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// ignore message from bots (including myself) and if not ready
	if msg.Author.Bot || !Relay.isReady() || msg.ChannelID != *Config.Discord.ChannelId {
		return
	}
	msgText, err := msg.ContentWithMoreMentionsReplaced(session);
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	SendIrc(msg.Author.Username, stripEmoji(msgText))
	// if message contains an attachment
	if msg.Attachments != nil && len(msg.Attachments) > 0 {
		for _, att := range msg.Attachments {
			SendIrc(msg.Author.Username, att.URL)
		}
	}
}
