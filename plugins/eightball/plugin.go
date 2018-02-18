package eightball

import (
	"math/rand"
	"regexp"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
	"github.com/njhanley/stoopid/config"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("8ball", func(b *bot.Bot) error {
	rand.Seed(time.Now().UnixNano())
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
	Text   []string
	Weight float64
}

type responses struct {
	resp []response
	rng  *wrng
}

func newResponses(resps []response) *responses {
	var sum float64
	for _, r := range resps {
		sum += r.Weight
	}

	weights := make([]float64, len(resps))
	for i := range weights {
		resps[i].Weight /= sum
		weights[i] = resps[i].Weight
	}

	rng := newWRNG(time.Now().UnixNano(), weights)

	return &responses{resps, rng}
}

func (rs *responses) choose() response {
	return rs.resp[rs.rng.get()]
}

var (
	answers, insults *responses
	defaultAnswers   = []response{
		{[]string{"It is certain."}, 1},
		{[]string{"It is decidedly so."}, 1},
		{[]string{"Without a doubt."}, 1},
		{[]string{"Yes definitely."}, 1},
		{[]string{"You may rely on it."}, 1},
		{[]string{"As I see it, yes."}, 1},
		{[]string{"Most likely."}, 1},
		{[]string{"Outlook good."}, 1},
		{[]string{"Yes."}, 1},
		{[]string{"Signs point to yes."}, 1},
		{[]string{"Reply hazy try again."}, 1},
		{[]string{"Ask again later."}, 1},
		{[]string{"Better not tell you now."}, 1},
		{[]string{"Cannot predict now."}, 1},
		{[]string{"Concentrate and ask again."}, 1},
		{[]string{"Don't count on it."}, 1},
		{[]string{"My reply is no."}, 1},
		{[]string{"My sources say no."}, 1},
		{[]string{"Outlook not so good."}, 1},
		{[]string{"Very doubtful."}, 1},
	}
	defaultInsults = []response{
		{[]string{"How should I know?"}, 1},
		{[]string{"What kind of question is that?"}, 1},
		{[]string{"I don't think you understand the meaning of \"yes or no\"."}, 1},
	}
)

func configure(c *config.Config) error {
	var x struct {
		Answers []response
		Insults []response
	}
	if c.Exists("8ball") {
		err := c.Get("8ball", &x)
		if err != nil {
			return err
		}
	}

	if len(x.Answers) == 0 {
		x.Answers = defaultAnswers
	}
	answers = newResponses(x.Answers)

	if len(x.Insults) == 0 {
		x.Insults = defaultInsults
	}
	insults = newResponses(x.Insults)

	return nil
}

var wrongQuestion = regexp.MustCompile("^(?i:how|what|when|where|which|who|why)")

func execute(s *dg.Session, m *dg.Message) error {
	var resp response
	if wrongQuestion.MatchString(m.Content) {
		resp = insults.choose()
	} else {
		resp = answers.choose()
	}
	for _, t := range resp.Text {
		_, err := s.ChannelMessageSend(m.ChannelID, t)
		if err != nil {
			return err
		}
	}
	return nil
}
