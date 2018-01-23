package name

import (
	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("name", func(b *bot.Bot) error {
	b.AddCommand(bot.ToOwnerCommand(command))
	return nil
})

var command = bot.SimpleCommand("name", execute, bot.SimpleCommandInfo{
	Comment:     "change bot nickname",
	Usage:       []string{"name [<nickname>]"},
	Description: "Change or reset the bot's nickname in the guild.",
})

func execute(s *dg.Session, m *dg.Message) error {
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}
	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}
	return s.GuildMemberNickname(ch.GuildID, "@me", m.Content)
}
