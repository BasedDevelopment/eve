/*
 * eve - management toolkit for libvirt servers
 * Copyright (C) 2022-2023  BNS Services LLC

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

var (
	k      = koanf.New(".")
	parser = toml.Parser()

	Config struct {
		Name string `koanf:"name"`

		API struct {
			Host        string `koanf:"host"`
			Port        int    `koanf:"port"`
			BehindProxy bool   `koanf:"behind_proxy"`
		} `koanf:"api"`

		Database struct {
			URL string `koanf:"url"`
		} `koanf:"database"`
	}
)

func Load(configPath string) (err error) {
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
