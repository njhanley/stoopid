package bot

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/njhanley/stoopid/config"
	"github.com/pkg/errors"
)

type Bot struct {
	Config  *config.Config
	Session *dg.Session

	commandsMu sync.RWMutex
	commands   map[string]Command

	pluginsMu sync.RWMutex
	plugins   map[string]Plugin

	defers []func()

	logger *log.Logger

	// immutable
	token   string
	owner   string
	sigil   string
	logpath string
}

func NewBot(cfg *config.Config) (*Bot, error) {
	bot := &Bot{
		Config:   cfg,
		commands: make(map[string]Command),
		plugins:  make(map[string]Plugin),
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

	err = b.Config.Get("sigil", &b.sigil)
	if err != nil {
		return err
	}

	if b.Config.Exists("logpath") {
		err = b.Config.Get("logpath", &b.logpath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) initLogger() error {
	out := io.Writer(os.Stderr)
	if b.logpath != "" {
		file, err := os.Create(filepath.Join(b.logpath, time.Now().Format(time.RFC3339)+".log"))
		if err != nil {
			return err
		}
		b.Defer(func() { file.Close() })
		out = io.MultiWriter(out, file)
	}
	b.logger = log.New(out, "", log.LstdFlags)

	msgType := []string{
		dg.LogError:         "ERROR",
		dg.LogWarning:       "WARNING",
		dg.LogInformational: "INFO",
		dg.LogDebug:         "DEBUG",
	}
	dg.Logger = func(level, _ int, format string, v ...interface{}) {
		b.Logf("[DG %s] %s\n", msgType[level], fmt.Sprintf(format, v...))
	}
	b.Session.LogLevel = dg.LogWarning

	return nil
}

func (b *Bot) Log(v ...interface{}) {
	b.logger.Print(v...)
}

func (b *Bot) Logf(format string, v ...interface{}) {
	b.logger.Printf(format, v...)
}

func (b *Bot) Logln(v ...interface{}) {
	b.logger.Println(v...)
}

func (b *Bot) Defer(fn func()) {
	b.defers = append(b.defers, fn)
}

func (b *Bot) connect() error {
	err := b.Session.Open()
	if err != nil {
		return err
	}
	b.Defer(func() { b.Session.Close() })
	return nil
}

func (b *Bot) Run() error {
	err := b.initLogger()
	if err != nil {
		return err
	}
	return b.connect()
}

func (b *Bot) Stop() {
	for i := len(b.defers) - 1; i >= 0; i-- {
		b.defers[i]()
	}
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

	n := strings.Index(msg.Content, " ")
	if n < 0 {
		n = len(msg.Content)
	}
	name := msg.Content[:n]
	msg.Content = strings.TrimSpace(msg.Content[n:])

	cmd := b.GetCommand(name)

	if cmd != nil && (!IsOwnerCommand(cmd) || msg.Author.ID == b.owner) {
		err := cmd.Execute(s, msg)
		if err != nil {
			b.Log(err)
		}
	}
}
