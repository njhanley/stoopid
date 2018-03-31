package say

import (
	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("say", func(b *bot.Bot) error {
	logf = b.Logf
	b.AddCommand(bot.ToOwnerCommand(command))
	return nil
})

var logf func(format string, v ...interface{})

var command = bot.SimpleCommand("say", execute, bot.SimpleCommandInfo{
	Comment:     "say a message",
	Usage:       []string{"say <message>"},
	Description: "Make the bot say the message.",
})

func execute(s *dg.Session, m *dg.Message) {
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		logf("[say] %v", err)
		return
	}

	if m.Content == "" {
		logf("[say] no argument")
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, m.Content)
	if err != nil {
		logf("[say] %v", err)
	}
}
