package bot

// Plugin is the interface for extending a bot.
// Plugins are identified by their names; loading multiple
// plugins with the same name is allowed but all but the last
// will be overwritten.
type Plugin interface {
	Name() string
	Load(*Bot) error
}

type LoadFunc func(*Bot) error

func SimplePlugin(name string, load LoadFunc) Plugin {
	return &simplePlugin{name, load}
}

type simplePlugin struct {
	name string
	load LoadFunc
}

func (p *simplePlugin) Name() string {
	return p.name
}

func (p *simplePlugin) Load(b *Bot) error {
	return p.load(b)
}
