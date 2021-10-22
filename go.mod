module github.com/IceflowRE/irc-discord-relay

go 1.16

require (
	github.com/bwmarrin/discordgo v0.23.2
	github.com/gorilla/websocket v1.4.1 // indirect; force to ovreride v1.4.0 due to security issues
	github.com/pkg/errors v0.9.1
	github.com/thoj/go-ircevent v0.0.0-20210723090443-73e444401d64
)
