package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	idRelay "github.com/IceflowRE/irc-discord-relay/pkg"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "./config.json", "config file")
	flag.Parse()

	err := idRelay.LoadConfig(configFile)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = idRelay.StartIRC()
	if err != nil {
		log.Println("irc"+err.Error())
		return
	}
	err = idRelay.StartDiscord()
	if err != nil {
		log.Println("discord"+err.Error())
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
		log.Println("discord"+discordErr.Error())
	}
	if ircErr != nil {
		log.Println("irc"+ircErr.Error())
	}
}
