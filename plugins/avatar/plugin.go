package avatar

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("avatar", func(b *bot.Bot) error {
	b.AddCommand(bot.ToOwnerCommand(command))
	return nil
})

var command = bot.SimpleCommand("avatar", execute, bot.SimpleCommandInfo{
	Comment:     "change avatar",
	Usage:       []string{"avatar"},
	Description: "Change the bot's avatar to the attached image or reset it to default if no image is attached with the command.",
})

func execute(s *dg.Session, m *dg.Message) error {
	n := len(m.Attachments)
	if n > 1 {
		return errors.New("more than one attachment")
	}

	avatar := "data:;base64,"
	if n != 0 {
		r, err := http.Get(m.Attachments[0].URL)
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			return err
		}

		switch mime := http.DetectContentType(b); mime {
		case "image/gif", "image/jpeg", "image/png":
			avatar = "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(b)
		default:
			return errors.New("invalid MIME type: " + mime)
		}
	}

	_, err := s.UserUpdate("", "", "", avatar, "")
	if err != nil {
		return err
	}

	return s.ChannelMessageDelete(m.ChannelID, m.ID)
}
