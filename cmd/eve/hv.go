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

package main

import (
	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/rs/zerolog/log"
)

func connHV(hv *controllers.HV) {
	// Connect to hypervisors, fetch VMs, and check for VM consistency
	if err := hv.Init(); err != nil {
		log.Warn().
			Err(err).
			Str("hostname", hv.Hostname).
			Msg("Failed to connect to hypervisor and fetch virtual machines")
	} else {
		log.Info().
			Str("hostname", hv.Hostname).
			Str("hvs", hv.Hostname).
			Int("vms", len(hv.VMs)).
			Msg("Connected to hypervisor and fetched virtual machines")
	}
}
