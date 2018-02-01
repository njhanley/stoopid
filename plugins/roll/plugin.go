package roll

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/bot"
	"github.com/njhanley/stoopid/config"
	"github.com/pkg/errors"
)

func Plugin() bot.Plugin {
	return plugin
}

var plugin = bot.SimplePlugin("roll", func(b *bot.Bot) error {
	rand.Seed(time.Now().UnixNano())
	err := configure(b.Config)
	if err != nil {
		return err
	}
	b.AddCommand(command)
	b.AddCommand(bot.ToHiddenCommand(emojiCommand))
	return nil
})

var command = bot.SimpleCommand("roll", execute, commandInfo)
var emojiCommand = bot.SimpleCommand("\U0001F3B2", execute, commandInfo)

var commandInfo = bot.SimpleCommandInfo{
	Comment:     "roll dice",
	Usage:       []string{"roll [<number of dice>]d<number of sides>[+|-<modifier>] [<text>]"},
	Description: fmt.Sprintf("Roll %d to %d dice each with %d to %d sides with an optional modifier between %d and %d. If <number of dice> is missing, it will default to the minimum. Additional text may be included after the command.", cfg.Dice.Min, cfg.Dice.Max, cfg.Sides.Min, cfg.Sides.Max, cfg.Modifier.Min, cfg.Modifier.Max),
}

type minmax struct {
	Min, Max int
}

var cfg = struct {
	Dice     minmax
	Sides    minmax
	Modifier minmax
}{
	Dice:     minmax{1, 100},
	Sides:    minmax{2, 1000},
	Modifier: minmax{-1000000, 1000000},
}

func configure(c *config.Config) error {
	if c.Exists("roll") {
		err := c.Get("roll", &cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

var rollRegexp = regexp.MustCompile("^([1-9][0-9]*)?d([1-9][0-9]*)([+-][1-9][0-9]*)?")

func execute(s *dg.Session, m *dg.Message) error {
	// match roll pattern
	loc := rollRegexp.FindStringSubmatchIndex(m.Content)
	if loc == nil {
		return nil
	}

	// convert
	var err error
	var tmp string
	var dice, sides, modifier int
	if loc[2] >= 0 {
		tmp = m.Content[loc[2]:loc[3]]
		dice, err = strconv.Atoi(tmp)
		if err != nil {
			return errors.Wrapf(err, "strconv.Atoi failed on dice (%q)", tmp)
		}
	} else {
		dice = cfg.Dice.Min
	}
	tmp = m.Content[loc[4]:loc[5]]
	sides, err = strconv.Atoi(tmp)
	if err != nil {
		return errors.Wrapf(err, "strconv.Atoi failed on sides (%q)", tmp)
	}
	if loc[6] >= 0 {
		tmp = m.Content[loc[6]:loc[7]]
		modifier, err = strconv.Atoi(tmp)
		if err != nil {
			return errors.Wrapf(err, "strconv.Atoi failed on modifier (%q)", tmp)
		}
	}

	// check limits
	if !(cfg.Dice.Min <= dice && dice <= cfg.Dice.Max) {
		return errors.Errorf("roll failed: dice out of bounds (%d, min = %d, max = %d)", dice, cfg.Dice.Min, cfg.Dice.Max)
	}
	if !(cfg.Sides.Min <= sides && sides <= cfg.Sides.Max) {
		return errors.Errorf("roll failed: sides out of bounds (%d, min = %d, max = %d)", sides, cfg.Sides.Min, cfg.Sides.Max)
	}
	if !(cfg.Modifier.Min <= modifier && modifier <= cfg.Modifier.Max) {
		return errors.Errorf("roll failed: modifier out of bounds (%d, min = %d, max = %d)", modifier, cfg.Modifier.Min, cfg.Modifier.Max)
	}

	total := modifier
	rolls := make([]string, dice)
	for i := range rolls {
		n := rand.Intn(sides) + 1
		total += n
		rolls[i] = strconv.Itoa(n)
	}

	text := strings.Join(rolls, " + ")
	if modifier != 0 {
		sign := " + "
		if modifier < 0 {
			modifier = -modifier
			sign = " - "
		}
		text += sign + strconv.Itoa(modifier)
	}
	if dice > 1 || modifier != 0 {
		text += " = " + strconv.Itoa(total)
	}

	_, err = s.ChannelMessageSend(m.ChannelID, text)
	return err
}
