package xkcd

import (
	"fmt"
	"regexp"
	"strconv"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("xkcd", func(b *bot.Bot) error {
	logf = b.Logf
	b.AddCommand(command)
	return nil
})

var logf func(format string, v ...interface{})

var command = bot.SimpleCommand("xkcd", execute, bot.SimpleCommandInfo{
	Comment:     "get xkcd comics",
	Usage:       []string{"xkcd", "xkcd <number>", "xkcd random"},
	Description: "Get xkcd comics.",
})

var numRegexp = regexp.MustCompile("^[1-9][0-9]*$")

func execute(s *dg.Session, m *dg.Message) {
	var (
		info *Info
		err  error
	)

	switch {
	case m.Content == "random":
		info, err = GetRandom()
	case m.Content == "", numRegexp.MatchString(m.Content):
		info, err = Get(m.Content)
	default:
		err = fmt.Errorf("invalid argument %q", m.Content)
	}
	if err != nil {
		logf("[xkcd] %v", err)
		return
	}

	msg := &dg.MessageEmbed{
		URL:   "https://xkcd.com/" + strconv.Itoa(info.Num) + "/",
		Title: "xkcd: " + info.Title,
		Image: &dg.MessageEmbedImage{
			URL: info.Img,
		},
		Footer: &dg.MessageEmbedFooter{
			Text: fmt.Sprintf("#%d, posted %s-%s-%s", info.Num, info.Year, info.Month, info.Day),
		},
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		logf("[xkcd] %v", err)
	}
}
