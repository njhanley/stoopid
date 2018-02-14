package weeb

import (
	"log"
	"unicode"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("weeb", func(b *bot.Bot) error {
	b.Session.AddHandler(handle)
	return nil
})

func handle(s *dg.Session, mc *dg.MessageCreate) {
	m := mc.Message
	for _, r := range m.Content {
		if unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana) {
			_, err := s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" is a filthy WEEB!")
			if err != nil {
				log.Print(err)
			}
			return
		}
	}
}
