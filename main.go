package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/njhanley/stoopid/bot"
	"github.com/njhanley/stoopid/config"
	"github.com/njhanley/stoopid/plugins/avatar"
	"github.com/njhanley/stoopid/plugins/crypto"
	"github.com/njhanley/stoopid/plugins/eightball"
	"github.com/njhanley/stoopid/plugins/name"
	"github.com/njhanley/stoopid/plugins/roll"
	"github.com/njhanley/stoopid/plugins/say"
	"github.com/njhanley/stoopid/plugins/status"
	"github.com/njhanley/stoopid/plugins/weeb"
	"github.com/njhanley/stoopid/plugins/xkcd"
	"golang.org/x/sys/unix"
)

var plugins = []bot.Plugin{
	avatar.Plugin(),
	crypto.Plugin(),
	eightball.Plugin(),
	name.Plugin(),
	roll.Plugin(),
	say.Plugin(),
	status.Plugin(),
	weeb.Plugin(),
	xkcd.Plugin(),
}

var configfile = flag.String("c", "config.json", "specify config location")

func main() {
	flag.Parse()

	cfg, err := config.New(*configfile)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := bot.NewBot(cfg)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan error, 1)
	go func() {
		for err := range c {
			log.Print(err)
		}
	}()
	bot.NotifyOnError(c)

	for _, p := range plugins {
		err = bot.AddPlugin(p)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = bot.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Disconnect()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, unix.SIGINT, unix.SIGTERM)
	<-sc
}
