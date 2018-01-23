package say

import (
	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("say", func(b *bot.Bot) error {
	b.AddCommand(bot.ToOwnerCommand(command))
	return nil
})

var command = bot.SimpleCommand("say", execute, bot.SimpleCommandInfo{
	Comment:     "say a message",
	Usage:       []string{"say <message>"},
	Description: "Make the bot say the message.",
})

func execute(s *dg.Session, m *dg.Message) error {
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}
	if m.Content != "" {
		_, err = s.ChannelMessageSend(m.ChannelID, m.Content)
	}
	return err
}
