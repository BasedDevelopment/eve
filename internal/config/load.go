package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

const configPath = "/etc/eve/config.toml"

var (
	k      = koanf.New(".")
	parser = toml.Parser()

	Config struct {
		Name string `koanf:"name"`

		API struct {
			Host string `koanf:"host"`
			Port int    `koanf:"port"`
		} `koanf:"api"`

		Database struct {
			URL string `koanf:"url"`
		} `koanf:"database"`
	}
)

func Load() (err error) {
	// Load from toml
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		return err
	}

	// Marshal into struct
	if err := k.Unmarshal("", &Config); err != nil {
		return err
	}

	// Validate config
	if err := validate(); err != nil {
		return err
	}
	return
}
