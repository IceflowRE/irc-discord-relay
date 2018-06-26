package idrelay

import (
	"crypto/tls"
	"errors"
	"fmt"
	"regexp"
	"strings"

	ircE "github.com/thoj/go-ircevent"
)

// StartIRC starts the IRC connection and adds the handlers
func StartIRC() error {
	iConn := ircE.IRC(*config.Irc.Nick, *config.Irc.Nick)
	iConn.VerboseCallbackHandler = false // DEV, use true for DEV
	iConn.UseTLS = true
	iConn.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// add callbacks
	valid := false
	for _, value := range *config.Irc.Sharing {
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
		return errors.New("no valid values in irc.sharing")
	}

	// connect to server callback
	iConn.AddCallback(
		"001",
		func(e *ircE.Event) {
			iConn.Join(*config.Irc.Channel)
			for _, msg := range *config.Irc.OnConnection {
				iConn.SendRaw(msg)
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
				if e.Arguments[2] == *config.irc.Channel {

				}
			},
		)*/

	// connect to IRC server
	err := iConn.Connect(*config.Irc.Server)
	if err != nil {
		return err
	}
	Relay.iConn = iConn

	// start main callback loop
	go func(iConn *ircE.Connection) {
		iConn.Loop()
	}(iConn)
	return nil
}

// send a message to the IRC channel
func sendIrc(msg string) {
	Relay.iConn.Privmsg(*config.Irc.Channel, msg)
}

// if it does not contain any mention it returns the message itself
func messageWithMention(msg string) string {
	if strings.ContainsRune(msg, '@') { // contains the message possible mentions
		// get the members, have the discord API limit in mind!
		members, err := Relay.dSession.GuildMembers(Relay.dGuildID, "", 1000)
		if err != nil { // something went wrong
			fmt.Println(err.Error())
			return msg
		}
		arr := strings.Split(msg, " ")
		for ind, val := range arr { // go through the message and check the parts
			if discordNickReg.MatchString(val) { // it is a mention
				for _, member := range members { // get the regarding discord id if its exists
					if val == ("@"+member.Nick) || val == ("@"+member.User.Username) {
						arr[ind] = "<@" + member.User.ID + ">"
						break
					}
				}
			}
		}
		return strings.Join(arr, " ")
	}
	return msg
}

// https://stackoverflow.com/a/10567935/6754440
// \x02: bold
// \x1F: underline
// \x16: italics
// \x1D: italics
// \x0F: normal
// \x03: colors
var ircFormat = regexp.MustCompile(`[\x02\x1F\x16\x1D\x0F]|\x03(\d\d?(,\d\d?)?)?`)

// PRIVMSG callback
// https://tools.ietf.org/html/rfc1459#section-4.4.1
func onIrcPrivmsg(e *ircE.Event) {
	if !Relay.isReady() || e.Nick == *config.Irc.Nick {
		return
	}

	msg := ircFormat.ReplaceAllString(e.Message(), "")
	sendDiscord("**<" + e.Nick + ">** " + messageWithMention(msg))
}

// CTCP_ACTION callback
func onIrcCtcpAction(e *ircE.Event) {
	if !Relay.isReady() || e.Nick == *config.Irc.Nick {
		return
	}

	msg := ircFormat.ReplaceAllString(e.Message(), "")
	sendDiscord("\\* " + e.Nick + " *" + messageWithMention(msg) + "*")
}

// JOIN callback
// https://tools.ietf.org/html/rfc2813#section-4.2.1
func onIrcJoin(e *ircE.Event) {
	if !Relay.isReady() || e.Nick == *config.Irc.Nick {
		return
	}
	sendDiscord("*" + e.Nick + "* has joined the channel")
}

// PART callback
// https://tools.ietf.org/html/rfc1459#section-4.2.2
func onIrcPart(e *ircE.Event) {
	if !Relay.isReady() {
		return
	}
	sendDiscord("*" + e.Nick + "* has left the channel (" + e.Message() + ")")
}

// NICK callback
// https://tools.ietf.org/html/rfc2813#section-4.1.3
func onIrcNick(e *ircE.Event) {
	if !Relay.isReady() {
		return
	}
	sendDiscord("*" + e.Nick + "* is now known as *" + e.Message() + "*")
}

// QUIT callback
// https://tools.ietf.org/html/rfc2813#section-4.1.5
func onIrcQuit(e *ircE.Event) {
	if !Relay.isReady() || e.Nick == *config.Irc.Nick {
		return
	}
	sendDiscord("*" + e.Nick + "* has quit (" + e.Message() + ")")
}

// KICK callback
// https://tools.ietf.org/html/rfc1459#section-4.2.8
func onIrcKick(e *ircE.Event) {
	if !Relay.isReady() || e.Nick == *config.Irc.Nick || len(e.Arguments) < 2 {
		return
	}
	sendDiscord("*" + e.Nick + "* kicked *" + e.Arguments[1] + "*")
}

// MODE callback
// https://tools.ietf.org/html/rfc1459#section-4.2.3
func onIrcMode(e *ircE.Event) {
	if !Relay.isReady() || e.Nick == *config.Irc.Nick || e.Nick == "" || len(e.Arguments) < 3 {
		return
	}
	sendDiscord("*" + e.Nick + "* set `" + e.Arguments[1] + "` on *" + e.Arguments[2] + "*")
}
