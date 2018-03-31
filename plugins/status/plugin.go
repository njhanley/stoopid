package status

import (
	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("status", func(b *bot.Bot) error {
	logf = b.Logf
	b.AddCommand(bot.ToOwnerCommand(command))
	return nil
})

var logf func(format string, v ...interface{})

var command = bot.SimpleCommand("status", execute, bot.SimpleCommandInfo{
	Comment:     "change bot status",
	Usage:       []string{"status [<game>]"},
	Description: "Change the bot's status.",
})

func execute(s *dg.Session, m *dg.Message) {
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		logf("[status] %v", err)
		return
	}
	err = s.UpdateStatus(0, m.Content)
	if err != nil {
		logf("[status] %v", err)
	}
}
