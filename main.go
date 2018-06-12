package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	idRelay "github.com/iceflowRE/irc-discord-relay/pkg"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "./config.json", "config file")
	flag.Parse()

	err := idRelay.LoadConfig(configFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = idRelay.StartIRC()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = idRelay.StartDiscord()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// waiting for close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	select {
	case <-sc:
	}
	discordErr, ircErr := idRelay.Relay.Close()
	if discordErr != nil {
		fmt.Println(discordErr.Error())
	}
	if ircErr != nil {
		fmt.Println(ircErr.Error())
	}
}
