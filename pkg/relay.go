package idrelay

import (
	"github.com/bwmarrin/discordgo"
	ircE "github.com/thoj/go-ircevent"
)

// Relay is the variable which unite the irc and discord access and will contain the variables for the running bots
var Relay = &idRelay{}

type idRelay struct {
	dSession *discordgo.Session
	dGuildID string
	iConn    *ircE.Connection
}

func (relay *idRelay) isReady() bool {
	return relay.dSession != nil && relay.iConn != nil
}

func (relay *idRelay) Close() (error, error) { // Discord exit, IRC exit; irc exit is nil everytime
	Relay.iConn.Quit()
	return Relay.dSession.Close(), nil
}
