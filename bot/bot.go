package bot

import (
	"strings"
	"sync"
	"unicode"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/config"
	"github.com/pkg/errors"
)

// Bot is a Discord bot.
type Bot struct {
	Config  *config.Config
	Session *dg.Session

	commandsMu sync.RWMutex
	commands   map[string]Command

	pluginsMu sync.RWMutex
	plugins   map[string]Plugin

	errMu sync.RWMutex
	err   chan<- error

	// immutable
	token string
	owner string
	sigil string
}

const DefaultSigil = "!"

// NewBot creates a bot.
func NewBot(c *config.Config) (*Bot, error) {
	bot := &Bot{
		Config:   c,
		commands: make(map[string]Command),
		plugins:  make(map[string]Plugin),
		sigil:    DefaultSigil,
	}

	err := bot.loadCfg()
	if err != nil {
		return nil, err
	}

	bot.Session, err = dg.New("Bot " + bot.token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	bot.Session.AddHandler(bot.messageCreate)

	bot.AddCommand(helpCommand{bot})

	return bot, nil
}

func (b *Bot) loadCfg() error {
	err := b.Config.Get("token", &b.token)
	if err != nil {
		return err
	}

	err = b.Config.Get("owner", &b.owner)
	if err != nil {
		return err
	}

	if b.Config.Exists("sigil") {
		err = b.Config.Get("sigil", &b.sigil)
		if err != nil {
			return err
		}
	}

	return nil
}

// Connect to Discord.
func (b *Bot) Connect() error {
	err := b.Session.Open()
	if err != nil {
		return errors.Wrap(err, "failed to connect")
	}
	return nil
}

// Disconnect from Discord.
func (b *Bot) Disconnect() error {
	err := b.Session.Close()
	if err != nil {
		return errors.Wrap(err, "failed to disconnect gracefully")
	}
	return nil
}

// NotifyOnError causes all errors returned from command execution or plugins to be sent to c.
// Sends will not block; it is the caller's responsibility to ensure c has a sufficient buffer.
func (b *Bot) NotifyOnError(c chan<- error) {
	b.errMu.Lock()
	b.err = c
	b.errMu.Unlock()
}

// SendError allows plugins to send their own errors for centralized logging.
func (b *Bot) SendError(err error) {
	b.errMu.RLock()
	select {
	case b.err <- err:
	default:
	}
	b.errMu.RUnlock()
}

// AddPlugin loads a plugin into the bot.
// A plugin with an empty name will have its Load method called
// but will not be retrievable with GetPlugin.
func (b *Bot) AddPlugin(p Plugin) error {
	b.pluginsMu.Lock()
	defer b.pluginsMu.Unlock()

	err := p.Load(b)
	if err != nil {
		return errors.Wrapf(err, "load plugin %q failed", p.Name())
	}

	if name := p.Name(); name != "" {
		b.plugins[name] = p
	}
	return nil
}

// GetPlugin retrieves a plugin by name.
// Nil is returned if no plugin with the given name is loaded.
func (b *Bot) GetPlugin(name string) Plugin {
	b.pluginsMu.RLock()
	defer b.pluginsMu.RUnlock()
	return b.plugins[name]
}

// AddCommand adds a command to the bot.
// Commands can be replaced by adding a command with the same name.
func (b *Bot) AddCommand(cmd Command) {
	b.commandsMu.Lock()
	b.commands[cmd.Name()] = cmd
	b.commandsMu.Unlock()
}

// GetCommand retrieves a command by name.
// Nil is returned if no command with the given name is loaded.
func (b *Bot) GetCommand(name string) Command {
	b.commandsMu.RLock()
	defer b.commandsMu.RUnlock()
	return b.commands[name]
}

// Sigil returns the string marking the beginning of a command.
func (b *Bot) Sigil() string {
	return b.sigil
}

func (b *Bot) messageCreate(s *dg.Session, m *dg.MessageCreate) {
	msg := m.Message

	if msg.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(msg.Content, b.sigil) {
		return
	}
	msg.Content = msg.Content[len(b.sigil):]

	n := strings.IndexFunc(msg.Content, unicode.IsSpace)
	if n < 0 {
		n = len(msg.Content)
	}
	name := msg.Content[:n]
	msg.Content = strings.TrimSpace(msg.Content[n:])

	cmd := b.GetCommand(name)

	if cmd != nil && (!IsOwnerCommand(cmd) || msg.Author.ID == b.owner) {
		err := cmd.Execute(s, msg)
		if err != nil {
			b.SendError(err)
		}
	}
}
