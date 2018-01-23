package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/pkg/errors"
)

// Config exposes a JSON file as a read-only key-value store.
type Config struct {
	mu  sync.RWMutex
	cfg map[string]json.RawMessage

	filename string // immutable
}

// New creates a Config and loads the file.
func New(filename string) (*Config, error) {
	c := &Config{filename: filename}
	err := c.load()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) load() error {
	b, err := ioutil.ReadFile(c.filename)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	cfg := make(map[string]json.RawMessage)
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	c.mu.Lock()
	c.cfg = cfg
	c.mu.Unlock()

	return nil
}

// Reload the config file from disk, replacing the old values.
//func (c *Config) Reload() error {
//	return c.load()
//}

// Exists reports whether a key exists.
func (c *Config) Exists(key string) bool {
	c.mu.RLock()
	_, ok := c.cfg[key]
	c.mu.RUnlock()
	return ok
}

// Get unmarshals a key's value into an interface using the same rules as json.Unmarshal.
// It returns an error if the key is not set.
func (c *Config) Get(key string, value interface{}) error {
	c.mu.RLock()
	b, ok := c.cfg[key]
	c.mu.RUnlock()
	if !ok {
		return errors.Errorf("key %q not found in config", key)
	}

	err := json.Unmarshal(b, value)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal key %q", key)
	}

	return nil
}
