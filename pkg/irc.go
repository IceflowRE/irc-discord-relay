package ircDiscordRelay

import (
	"crypto/tls"

	"github.com/thoj/go-ircevent"
)

func onIrcPrivmsg(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick {
		return
	}
	SendDiscord("**<" + e.Nick + ">** " + e.Message())
}

func onIrcCtcpAction(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick {
		return
	}
	SendDiscord("**<" + e.Nick + ">** *" + e.Message() + "*")
}

func onIrcJoin(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick {
		return
	}
	SendDiscord("*<" + e.Nick + ">* has joined the channel")
}

func onIrcPart(e *irc.Event) {
	if !Relay.isReady() {
		return
	}
	SendDiscord("*<" + e.Nick + ">* has left the channel (" + e.Message() + ")")
}

func onIrcQuit(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick {
		return
	}
	SendDiscord("*<" + e.Nick + ">* has quit (" + e.Message() + ")")
}

func onIrcKick(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick || len(e.Arguments) < 2 {
		return
	}
	SendDiscord("*<" + e.Nick + ">* kicked *" + e.Arguments[1] + "*")
}

func onIrcMode(e *irc.Event) {
	if !Relay.isReady() || e.Nick == *Config.Irc.Nick || e.Nick == "" || len(e.Arguments) < 3 {
		return
	}
	SendDiscord("*<" + e.Nick + ">* set `" + e.Arguments[1] + "` on *" + e.Arguments[2] + "*")
}

func SendIrc(msg string) {
	Relay.iConn.Privmsg(*Config.Irc.Channel, msg)
}

func StartIRC() error {
	iConn := irc.IRC(*Config.Irc.Nick, *Config.Irc.Nick)
	iConn.VerboseCallbackHandler = false // DEV, use true for DEV
	iConn.UseTLS = true
	iConn.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	iConn.AddCallback("001", func(e *irc.Event) {
		iConn.Join(*Config.Irc.Channel)
		//iConn.SendRaw("/msg NickServ identify Blizzard " + password)
	},
	)
	iConn.AddCallback("PRIVMSG", onIrcPrivmsg)
	iConn.AddCallback("CTCP_ACTION", onIrcCtcpAction)
	iConn.AddCallback("JOIN", onIrcJoin)
	iConn.AddCallback("PART", onIrcPart)
	iConn.AddCallback("QUIT", onIrcQuit)
	iConn.AddCallback("KICK", onIrcKick)
	iConn.AddCallback("MODE", onIrcMode)

	err := iConn.Connect(*Config.Irc.Server)
	if err != nil {
		return err
	}
	Relay.iConn = iConn
	go func(iConn *irc.Connection) {
		iConn.Loop()
	}(iConn)
	return nil
}
