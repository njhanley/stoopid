package weeb

import (
	"strings"
	"sync"
	"time"
	"unicode"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
	"github.com/njhanley/stoopid/config"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("weeb", func(b *bot.Bot) error {
	log = b.Log
	sigil = b.Sigil()
	err := configure(b.Config)
	if err != nil {
		return err
	}
	b.Session.AddHandler(handle)
	return nil
})

var (
	cooldown = 5 * time.Minute
	log      func(v ...interface{})
	sigil    string

	mutex sync.Mutex
	weebs = make(map[string]time.Time)
)

func configure(c *config.Config) error {
	if c.Exists("weeb") {
		var x struct {
			Cooldown string
		}
		err := c.Get("weeb", &x)
		if err != nil {
			return err
		}
		cooldown, err = time.ParseDuration(x.Cooldown)
		if err != nil {
			return err
		}
	}
	return nil
}

func containsJapanese(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana) {
			return true
		}
	}
	return false
}

func getDisplayName(st *dg.State, channelID, userID string) (string, error) {
	ch, err := st.Channel(channelID)
	if err != nil {
		return "", err
	}

	mem, err := st.Member(ch.GuildID, userID)
	if err != nil {
		return "", err
	}

	if mem.Nick != "" {
		return mem.Nick, nil
	}
	return mem.User.Username, nil
}

func handle(s *dg.Session, mc *dg.MessageCreate) {
	m := mc.Message
	if m.Author.ID == s.State.User.ID ||
		strings.HasPrefix(m.Content, sigil) ||
		!containsJapanese(m.Content) {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	last := weebs[m.Author.ID]
	weebs[m.Author.ID] = time.Now()
	if time.Since(last) < cooldown {
		return
	}

	name, err := getDisplayName(s.State, m.ChannelID, m.Author.ID)
	if err != nil {
		log(err)
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, name+" is a filthy WEEB!")
	if err != nil {
		log(err)
	}
}
