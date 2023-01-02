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

import "fmt"

// Check for required fields in the config file
func validate() error {
	if Config.Name == "" {
		return fmt.Errorf("Configuration: name is required")
	}

	if Config.API.Host == "" {
		return fmt.Errorf("Configuration: api.host is required")
	}

	if Config.API.Port == 0 {
		return fmt.Errorf("Configuration: api.port is required")
	}

	if Config.Database.URL == "" {
		return fmt.Errorf("Configuration: database.url is required")
	}
	return nil
}
