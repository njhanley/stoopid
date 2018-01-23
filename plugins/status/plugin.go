package status

import (
	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("status", func(b *bot.Bot) error {
	b.AddCommand(bot.ToOwnerCommand(command))
	return nil
})

var command = bot.SimpleCommand("status", execute, bot.SimpleCommandInfo{
	Comment:     "change bot status",
	Usage:       []string{"status [<game>]"},
	Description: "Change the bot's status.",
})

func execute(s *dg.Session, m *dg.Message) error {
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}
	return s.UpdateStatus(0, m.Content)
}
