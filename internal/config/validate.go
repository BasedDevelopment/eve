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
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Check for required fields in the config file
func validate() error {
	if err := validation.Validate(Config.Hostname, validation.Required, is.DNSName); err != nil {
		return fmt.Errorf("Configuration(hostname): %w", err)
	}

	if Config.TLSPath == "" {
		return fmt.Errorf("Configuration(tls_path): is required")
	}

	if err := validation.Validate(Config.API.Host, validation.Required, is.Host); err != nil {
		return fmt.Errorf("Configuration(api.host): %w", err)
	}

	if err := validation.Validate(Config.API.Port, validation.Required, validation.Min(1), validation.Max(65535)); err != nil {
		return fmt.Errorf("Configuration(api.port): %w", err)
	}

	if err := validation.Validate(Config.Database.URL, validation.Required); err != nil {
		return fmt.Errorf("Configuration(database.url): %w", err)
	}

	return nil
}
