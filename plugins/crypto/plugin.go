package crypto

import (
	"fmt"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("crypto", func(b *bot.Bot) error {
	b.AddCommand(command)
	return nil
})

var command = bot.SimpleCommand("crypto", execute, bot.SimpleCommandInfo{
	Comment:     "check exchange rates",
	Usage:       []string{"crypto <pair>"},
	Description: "Check cryptocurrency exchange rates.\nExample pair: `btcusd`",
})

const (
	increase = 47369    // #00b909
	decrease = 12977670 // #c60606
)

func ftos(f float64) string {
	return strconv.FormatFloat(f, 'G', -1, 64)
}

func execute(s *dg.Session, m *dg.Message) error {
	if m.Content == "" {
		return nil
	}

	pr, err := GetPair(m.Content)
	if err != nil {
		return err
	}
	if len(pr.Result.Markets) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Invalid pair.")
		return err
	}
	market := pr.Result.Markets[0]

	er, err := GetExchange(market.Exchange)
	if err != nil {
		return err
	}
	exchange := er.Result

	msr, err := GetMarketSummary(market.Exchange, market.Pair)
	if err != nil {
		return err
	}
	summary := msr.Result

	msg := &dg.MessageEmbed{
		URL:   "https://cryptowat.ch/" + market.Exchange + "/" + market.Pair,
		Title: strings.Title(exchange.Name) + ": " + strings.ToUpper(market.Pair),
		Fields: []*dg.MessageEmbedField{
			{Name: "Latest", Value: ftos(summary.Price.Last), Inline: true},
			{Name: "High", Value: ftos(summary.Price.High), Inline: true},
			{Name: "Low", Value: ftos(summary.Price.Low), Inline: true},
			{Name: "Change (24H)", Value: fmt.Sprintf("%+.3f%% (%+G)", 100*summary.Price.Change.Percentage, summary.Price.Change.Absolute), Inline: true},
			{Name: "Volume", Value: ftos(summary.Volume), Inline: true},
		},
		Footer: &dg.MessageEmbedFooter{
			Text: "Data provided by https://cryptowat.ch/",
		},
	}
	if msr.Result.Price.Change.Absolute > 0 {
		msg.Color = increase
	} else {
		msg.Color = decrease
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	return err
}
