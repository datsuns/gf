package main

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/pelletier/go-toml/v2"
)

const (
	DefaultScrollLines = 7
)

type ConfigEntry struct {
	LeftPath    Path
	RightPath   Path
	ScrollLines int
	JumpList    map[string]string
	Editor      string
}

type Config struct {
	Path Path
	Body *ConfigEntry
}

func LoadConfig(path Path) (*Config, error) {
	var entry ConfigEntry
	var err error
	raw, err := os.ReadFile(string(path))
	if err != nil {
		cfg, e := GenDefaultConfig(path)
		if e != nil {
			return nil, errors.Wrap(e, fmt.Sprintf("GenDefaultConfig(%v)", path))
		}
		return &Config{Path: path, Body: cfg}, nil
	}
	err = toml.Unmarshal(raw, &entry)
	if err != nil {
		return nil, errors.Wrap(err, "toml.Unmarshal()")
	}
	if entry.ScrollLines == 0 {
		entry.ScrollLines = DefaultScrollLines
	}
	if entry.JumpList == nil {
		entry.JumpList = map[string]string{}
	}
	return &Config{Path: path, Body: &entry}, nil
}

func GenDefaultConfig(path Path) (*ConfigEntry, error) {
	c, _ := os.Getwd()
	d := &ConfigEntry{
		LeftPath:    Path(c),
		RightPath:   Path(c),
		ScrollLines: DefaultScrollLines,
		JumpList:    map[string]string{},
		Editor:      "",
	}
	err := save(path, d)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("save(%v)", path))
	}
	return d, nil
}

func (c *Config) Save() error {
	return save(c.Path, c.Body)
}

func (c *Config) LeftPath() Path {
	return c.Body.LeftPath
}

func (c *Config) RightPath() Path {
	return c.Body.RightPath
}

func save(path Path, e *ConfigEntry) error {
	raw, err := toml.Marshal(*e)
	if err != nil {
		return errors.WithStack(err)
	}
	if e := os.WriteFile(string(path), raw, 0666); e != nil {
		return errors.WithStack(e)
	}
	return nil
}
