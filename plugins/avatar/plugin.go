package avatar

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("avatar", func(b *bot.Bot) error {
	logf = b.Logf
	b.AddCommand(bot.ToOwnerCommand(command))
	return nil
})

var logf func(format string, v ...interface{})

var command = bot.SimpleCommand("avatar", execute, bot.SimpleCommandInfo{
	Comment:     "change avatar",
	Usage:       []string{"avatar"},
	Description: "Change the bot's avatar to the attached image or reset it to default if no image is attached with the command.",
})

func execute(s *dg.Session, m *dg.Message) {
	n := len(m.Attachments)
	if n > 1 {
		logf("[avatar] more than one attachment")
		return
	}

	avatar := "data:;base64,"
	if n != 0 {
		r, err := http.Get(m.Attachments[0].URL)
		if err != nil {
			logf("[avatar] %v", err)
			return
		}
		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			logf("[avatar] %v", err)
			return
		}

		switch mime := http.DetectContentType(b); mime {
		case "image/gif", "image/jpeg", "image/png":
			avatar = "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(b)
		default:
			logf("[avatar] invalid MIME type %q", mime)
			return
		}
	}

	_, err := s.UserUpdate("", "", "", avatar, "")
	if err != nil {
		logf("[avatar] %v", err)
		return
	}

	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		logf("[avatar] %v", err)
	}
}
