package eightball

import (
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
	"github.com/njhanley/stoopid/config"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("8ball", func(b *bot.Bot) error {
	err := configure(b.Config)
	if err != nil {
		return err
	}
	b.AddCommand(command)
	b.AddCommand(bot.ToHiddenCommand(emojiCommand))
	return nil
})

var command = bot.SimpleCommand("8ball", execute, commandInfo)
var emojiCommand = bot.SimpleCommand("\U0001F3B1", execute, commandInfo)

var commandInfo = bot.SimpleCommandInfo{
	Comment:     "ask a yes-no question",
	Usage:       []string{"8ball [<question>]"},
	Description: "Ask the 8ball a yes or no question (question optional).",
}

type response struct {
	Text   string
	Weight float64
}

var (
	responses []response
	rng       *wrng
)

var defaultResponses = []response{
	{"It is certain.", 1},
	{"It is decidedly so.", 1},
	{"Without a doubt.", 1},
	{"Yes definitely.", 1},
	{"You may rely on it.", 1},
	{"As I see it, yes.", 1},
	{"Most likely.", 1},
	{"Outlook good.", 1},
	{"Yes.", 1},
	{"Signs point to yes.", 1},
	{"Reply hazy try again.", 1},
	{"Ask again later.", 1},
	{"Better not tell you now.", 1},
	{"Cannot predict now.", 1},
	{"Concentrate and ask again.", 1},
	{"Don't count on it.", 1},
	{"My reply is no.", 1},
	{"My sources say no.", 1},
	{"Outlook not so good.", 1},
	{"Very doubtful.", 1},
}

func configure(c *config.Config) error {
	if c.Exists("8ball") {
		err := c.Get("8ball", &responses)
		if err != nil {
			return err
		}
	} else {
		responses = defaultResponses
	}

	var sum float64
	for _, r := range responses {
		sum += r.Weight
	}

	weights := make([]float64, len(responses))
	for i := range weights {
		weights[i] = responses[i].Weight / sum
	}

	rng = newWRNG(time.Now().UnixNano(), weights)

	return nil
}

func execute(s *dg.Session, m *dg.Message) error {
	_, err := s.ChannelMessageSend(m.ChannelID, responses[rng.get()].Text)
	return err
}
