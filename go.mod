module github.com/IceflowRE/irc-discord-relay

go 1.18

require (
	github.com/bwmarrin/discordgo v0.25.0
	github.com/gorilla/websocket v1.5.0 // indirect; force to ovreride v1.4.0 due to security issues
	github.com/pkg/errors v0.9.1
	github.com/thoj/go-ircevent v0.0.0-20210723090443-73e444401d64
	golang.org/x/crypto v0.0.0-20220427172511-eb4f295cb31f // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sys v0.0.0-20220429233432-b5fbb4746d32 // indirect
)
