package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	idRelay "github.com/iceflowre/ircDiscordRelay/pkg"
)

func main() {
	idRelay.LoadConfig()

	err := idRelay.StartIRC()
	if err != nil {
		panic(err)
	}
	err = idRelay.StartDiscord()
	if err != nil {
		panic(err)
	}

	// waiting for close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	select {
	case <-sc:
	}
	discordErr, ircErr := idRelay.Relay.Close()
	if discordErr  != nil {
		fmt.Errorf(discordErr.Error())
	}
	if ircErr  != nil {
		fmt.Errorf(ircErr.Error())
	}
}
