package ircDiscordRelay

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/thoj/go-ircevent"
)

func StartIRC() error {
	iConn := irc.IRC(*Config.Irc.Nick, *Config.Irc.Nick)
	iConn.VerboseCallbackHandler = false // DEV, use true for DEV
	iConn.UseTLS = true
	iConn.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// add callbacks
	valid := false
	for _, value := range *Config.Discord.Sharing {
		switch value {
		case "message":
			valid = true
			iConn.AddCallback("PRIVMSG", onIrcPrivmsg)
		case "me":
			valid = true
			iConn.AddCallback("CTCP_ACTION", onIrcCtcpAction)
		case "join":
			valid = true
			iConn.AddCallback("JOIN", onIrcJoin)
		case "leaving":
			valid = true
			iConn.AddCallback("PART", onIrcPart)
		case "nick":
			valid = true
			iConn.AddCallback("NICK", onIrcNick)
		case "quit":
			valid = true
			iConn.AddCallback("QUIT", onIrcQuit)
		case "kick":
			valid = true
			iConn.AddCallback("KICK", onIrcKick)
		case "mode":
			valid = true
			iConn.AddCallback("MODE", onIrcMode)
		default:
			fmt.Println("Invalid irc.sharing value '" + value + "' will be ignored.")
		}
	}
	if !valid {
		return errors.New("No valid values in irc.sharing.")
	}

	// connect to server callback
	iConn.AddCallback(
		"001",
		func(e *irc.Event) {
			iConn.Join(*Config.Irc.Channel)
			for target, msg := range *Config.Irc.OnConnection {
				iConn.Privmsg(target, msg)
			}
		},
	)
	/*
	iConn.AddCallback( // join successful
		"353",
		func(e *irc.Event) {
			if len(e.Arguments) < 3 {
				return
			}
			if e.Arguments[2] == *Config.Irc.Channel {

			}
		},
	)*/

	// connect to IRC server
	err := iConn.Connect(*Config.Irc.Server)
	if err != nil {
		return err
	}
	Relay.iConn = iConn

	// start main callback loop
	go func(iConn *irc.Connection) {
		iConn.Loop()
	}(iConn)
	return nil
}

func SendIrc(msg string) {
	Relay.iConn.Privmsg(*Config.Irc.Channel, msg)
}

// if it does not contain any mention it returns the msg itself
func messageWithMention(msg string) string {
	if strings.ContainsRune(msg, '@') { // contains the message possible mentions
		// get the members, have the Discord API limit in mind!
		members, err := Relay.dSession.GuildMembers(Relay.dGuildId, "", 1000)
		if err != nil { // something went wrong
			fmt.Println(err.Error())
			return ""
		}
		arr := strings.Split(msg, " ")
		for ind, val := range arr { // go through the message and check the parts
			if discordNickReg.MatchString(val) { // it is a mention
				for _, member := range members { // get the regarding discord id if its exists
					if val == ("@"+member.Nick) || val == ("@"+member.User.Username) {
						arr[ind] = "<@" + member.User.ID + ">"
						break;
					}
				}
			}
		}
		return strings.Join(arr, " ")
	}
	return msg
}

func onIrcPrivmsg(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick {
		return
	}
	SendDiscord("**<" + e.Nick + ">** " + messageWithMention(e.Message()))
}

func onIrcCtcpAction(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick {
		return
	}
	SendDiscord("**<" + e.Nick + ">** *" + messageWithMention(e.Message()) + "*")
}

func onIrcJoin(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick {
		return
	}
	SendDiscord("*" + e.Nick + "* has joined the channel")
}

func onIrcPart(e *irc.Event) {
	if !Relay.isReady() {
		return
	}
	SendDiscord("*" + e.Nick + "* has left the channel (" + e.Message() + ")")
}

func onIrcNick(e *irc.Event) {
	if !Relay.isReady() {
		return
	}
	SendDiscord("*" + e.Nick + "* is now known as *" + e.Message() + "*")
}

func onIrcQuit(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick {
		return
	}
	SendDiscord("*" + e.Nick + "* has quit (" + e.Message() + ")")
}

func onIrcKick(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick || len(e.Arguments) < 2 {
		return
	}
	SendDiscord("*" + e.Nick + "* kicked *" + e.Arguments[1] + "*")
}

func onIrcMode(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick || e.Nick == "" || len(e.Arguments) < 3 {
		return
	}
	SendDiscord("*" + e.Nick + "* set `" + e.Arguments[1] + "` on *" + e.Arguments[2] + "*")
}
