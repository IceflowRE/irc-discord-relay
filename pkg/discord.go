package idrelay

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

var discordNickReg = regexp.MustCompile("@[a-zA-Z0-9]*")

// StartDiscord connects the Discord bot to discord and starts the handlers
func StartDiscord() error {
	session, err := discordgo.New("Bot " + *config.Discord.Token)
	if err != nil {
		return err
	}

	session.AddHandler(func(session *discordgo.Session, msg *discordgo.Ready) {
		err = session.UpdateGameStatus(0, *config.Irc.Channel+" relay")
		if err != nil {
			log.Println(err.Error())
		}
	})
	valid := false
	for _, value := range *config.Discord.Sharing {
		switch value {
		case "message":
			valid = true
			session.AddHandler(onDiscordMsg)
		default:
			log.Println("Invalid discord.sharing value '" + value + "' will be ignored.")
		}
	}
	if !valid {
		return errors.New("no valid values in discord.sharing")
	}

	err = session.Open()
	if err != nil {
		return err
	}
	Relay.dSession = session

	chn, err := Relay.dSession.Channel(*config.Discord.ChannelID)
	if err != nil {
		return err
	}
	Relay.dGuildID = chn.GuildID

	return nil
}

// send message on discord
func sendDiscord(msg string) {
	_, err := Relay.dSession.ChannelMessageSend(*config.Discord.ChannelID, msg)
	if err != nil {
		log.Println(err.Error())
	}
}

var emojiRe = regexp.MustCompile("(<)a?(:.*:)[0-9]*(>)")

// removes the unique id from the emoji part
func stripEmoji(msg string) string {
	return emojiRe.ReplaceAllString(msg, "$1$2$3")
}

// on discord message
func onDiscordMsg(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// ignore message from bots (including myself) and if not ready
	if msg.Author.Bot || !Relay.isReady() || msg.ChannelID != *config.Discord.ChannelID {
		return
	}
	msgText, err := msg.ContentWithMoreMentionsReplaced(session)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var sender string
	memb, err := session.State.Member(Relay.dGuildID, msg.Author.ID)
	if err != nil {
		log.Println("Could not get the nickname, fallback to username!")
		sender = msg.Author.Username
	} else if memb.Nick == "" {
		sender = msg.Author.Username
	} else {
		sender = memb.Nick
	}

	// remove the emoji id of the emoji string, affects mostly only server specific emojis
	msgText = stripEmoji(msgText)
	// send all lines of the discord message
	for _, msgPart := range strings.Split(msgText, "\n") {
		sendIrc("<" + sender + "> " + msgPart)
	}
	// if message contains an attachment
	if msg.Attachments != nil && len(msg.Attachments) > 0 {
		for _, att := range msg.Attachments {
			sendIrc("<" + msg.Author.Username + "> " + att.URL)
		}
	}
}
