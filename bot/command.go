package bot

import (
	"sort"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

// Command is the interface for commands.
// Comment, Usage, and Description are used for help information.
//
// Commands are invoked by messages beginning with the command sigil
// and the name of a command. The sigil and name are removed from
// the message, as well as leading and trailing whitespace. The
// message is then passed into Execute.
type Command interface {
	Name() string        // command name
	Comment() string     // a short description
	Usage() []string     // command syntax
	Description() string // a detailed description
	Execute(*dg.Session, *dg.Message) error
}

// HiddenCommand is an interface for
// commands that should not appear in the help list.
type HiddenCommand interface {
	Command
	Hidden()
}

// IsHiddenCommand reports if a command implements HiddenCommand.
func IsHiddenCommand(cmd Command) bool {
	_, ok := cmd.(HiddenCommand)
	return ok
}

type hiddenCommand struct{ Command }

func (c hiddenCommand) Hidden() {}

// ToHiddenCommand decorates a command with
// a Hidden method, implementing HiddenCommand.
func ToHiddenCommand(cmd Command) HiddenCommand {
	return hiddenCommand{cmd}
}

// OwnerCommand is an interface for
// commands that only the owner can use.
type OwnerCommand interface {
	Command
	Owner()
}

// IsOwnerCommand reports if a command implements OwnerCommand.
func IsOwnerCommand(cmd Command) bool {
	_, ok := cmd.(OwnerCommand)
	return ok
}

type ownerCommand struct{ Command }

func (c ownerCommand) Owner() {}

// ToOwnerCommand decorates a command with
// an Owner method, implementing OwnerCommand.
func ToOwnerCommand(cmd Command) OwnerCommand {
	return ownerCommand{cmd}
}

// SimpleCommandInfo specifies the help information for SimpleCommand.
type SimpleCommandInfo struct {
	Comment     string
	Usage       []string
	Description string
}

// ExecuteFunc is a function that implements Execute for SimpleCommand.
type ExecuteFunc func(*dg.Session, *dg.Message) error

// SimpleCommand is a convenience function
// for creating commands from functions.
func SimpleCommand(name string, fn ExecuteFunc, i SimpleCommandInfo) Command {
	return &simpleCommand{name: name, exec: fn, info: i}
}

type simpleCommand struct {
	name string
	exec ExecuteFunc
	info SimpleCommandInfo
}

func (c *simpleCommand) Name() string {
	return c.name
}

func (c *simpleCommand) Comment() string {
	return c.info.Comment
}

func (c *simpleCommand) Usage() []string {
	// create copy to enforce immutability
	return append([]string(nil), c.info.Usage...)
}

func (c *simpleCommand) Description() string {
	return c.info.Description
}

func (c *simpleCommand) Execute(s *dg.Session, m *dg.Message) error {
	return c.exec(s, m)
}

type helpCommand struct {
	*Bot
}

func (c helpCommand) Name() string {
	return "help"
}

func (c helpCommand) Comment() string {
	return "get info about commands"
}

func (c helpCommand) Usage() []string {
	return []string{"help [<command>]"}
}

func (c helpCommand) Description() string {
	return "Get information about commands. If a command is not specified, list all commands."
}

func (c helpCommand) Execute(s *dg.Session, m *dg.Message) error {
	if m.Content == "" {
		return c.helplist(s, m)
	}
	return c.help(s, m)
}

const missingText = "<undefined>"

func (c helpCommand) help(s *dg.Session, m *dg.Message) error {
	cmd := c.GetCommand(m.Content)
	if cmd == nil {
		_, err := s.ChannelMessageSend(m.ChannelID, "Command not found.")
		return err
	}

	ownercmd := IsOwnerCommand(cmd)
	if ownercmd && m.Author.ID != c.owner {
		_, err := s.ChannelMessageSend(m.ChannelID, "You do not have permission to use that command.")
		return err
	}

	usages := cmd.Usage()
	for i, s := range usages {
		usages[i] = "`" + c.sigil + s + "`"
	}
	usage := strings.Join(usages, "\n")
	if usage == "" {
		usage = missingText
	}

	description := cmd.Description()
	if description == "" {
		description = missingText
	}
	if ownercmd {
		description += "\n\nOwner only."
	}

	fields := []*dg.MessageEmbedField{
		{Name: "Usage:", Value: usage},
		{Name: "Description:", Value: description},
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &dg.MessageEmbed{
		Fields: fields,
		Footer: &dg.MessageEmbedFooter{
			Text: "For a list of all commands, use " + c.sigil + "help",
		},
	})
	return err
}

func (c helpCommand) helplist(s *dg.Session, m *dg.Message) error {
	c.commandsMu.RLock()
	defer c.commandsMu.RUnlock()

	commands := make([]*dg.MessageEmbedField, 0, len(c.commands))
	for name, cmd := range c.commands {
		ownercmd := IsOwnerCommand(cmd)
		if IsHiddenCommand(cmd) || (ownercmd && m.Author.ID != c.owner) {
			continue
		}

		command := &dg.MessageEmbedField{Name: name, Value: cmd.Comment()}
		if command.Name == "" {
			command.Name = missingText
		}
		if command.Value == "" {
			command.Value = missingText
		}
		if ownercmd {
			command.Value += "\nowner only"
		}
		commands = append(commands, command)
	}

	sort.Slice(commands, func(i, j int) bool { return commands[i].Name < commands[j].Name })

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &dg.MessageEmbed{
		Title:  "Commands:",
		Fields: commands,
		Footer: &dg.MessageEmbedFooter{Text: "For more information, use " + c.sigil + "help <command>"},
	})
	return err
}
