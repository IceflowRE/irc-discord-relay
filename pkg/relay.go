package ircDiscordRelay

import (
	"github.com/bwmarrin/discordgo"
	"github.com/thoj/go-ircevent"
)

var Relay *idRelay = &idRelay{}

type idRelay struct {
	dSession *discordgo.Session
	dGuildId string
	iConn    *irc.Connection
}

func (relay *idRelay) isReady() bool {
	return relay.dSession != nil && relay.iConn != nil
}

func (relay *idRelay) Close() (error, error) { // Discord exit, IRC exit; irc exit is nil everytime
	Relay.iConn.Quit()
	return Relay.dSession.Close(), nil
}
