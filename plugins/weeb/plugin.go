package weeb

import (
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
	sendError = b.SendError
	err := configure(b.Config)
	if err != nil {
		return err
	}
	b.Session.AddHandler(handle)
	return nil
})

func configure(c *config.Config) error {
	d := 5 * time.Minute
	if c.Exists("weeb") {
		var x struct {
			Cooldown string
		}
		err := c.Get("weeb", &x)
		if err != nil {
			return err
		}
		d, err = time.ParseDuration(x.Cooldown)
		if err != nil {
			return err
		}
	}
	cooldown = newCooldownTimer(d)
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

type cooldownTimer struct {
	mu    sync.Mutex
	dur   time.Duration
	users map[string]time.Time
}

func newCooldownTimer(d time.Duration) *cooldownTimer {
	return &cooldownTimer{
		dur:   d,
		users: make(map[string]time.Time),
	}
}

func (c *cooldownTimer) ended(userID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return time.Since(c.users[userID]) > c.dur
}

func (c *cooldownTimer) update(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.users[userID] = time.Now()
}

var (
	cooldown  *cooldownTimer
	sendError func(error)
)

func handle(s *dg.Session, mc *dg.MessageCreate) {
	m := mc.Message
	if containsJapanese(m.Content) {
		defer cooldown.update(m.Author.ID)
		if !cooldown.ended(m.Author.ID) {
			return
		}
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
