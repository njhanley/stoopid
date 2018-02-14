package weeb

import (
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
	sendError = b.SendError
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

func getDisplayName(st *dg.State, channelID, userID string) (name string, err error) {
	ch, err := st.Channel(channelID)
	if err != nil {
		return "", err
	}
	mem, err := st.Member(ch.GuildID, userID)
	if err != nil {
		return "", err
	}
	if name = mem.Nick; name == "" {
		name = mem.User.Username
	}
	return name, nil
}

var (
	lastCallout map[string]time.Time
	sendError   func(error)
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
		name, err := getDisplayName(s.State, m.ChannelID, m.Author.ID)
		if err != nil {
			sendError(err)
		}
		_, err = s.ChannelMessageSend(m.ChannelID, name+" is a filthy WEEB!")
		if err != nil {
			sendError(err)
		}
	}
}
