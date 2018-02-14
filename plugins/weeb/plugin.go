package weeb

import (
	"log"
	"sync"
	"time"
	"unicode"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("weeb", func(b *bot.Bot) error {
	lastCallout = make(map[string]time.Time)
	b.Session.AddHandler(handle)
	return nil
})

func containsJapanese(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana) {
			return true
		}
	}
	return false
}

var (
	lastCallout map[string]time.Time
	mutex       sync.Mutex
)

func handle(s *dg.Session, mc *dg.MessageCreate) {
	m := mc.Message
	mutex.Lock()
	defer mutex.Unlock()

	if time.Since(lastCallout[m.Author.ID]) < time.Minute {
		return
	}

	if containsJapanese(m.Content) {
		lastCallout[m.Author.ID] = time.Now()
		_, err := s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" is a filthy WEEB!")
		if err != nil {
			log.Print(err)
		}
	}
}
