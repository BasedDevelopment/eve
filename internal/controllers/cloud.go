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

package controllers

import (
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type HVList struct {
	mutex sync.Mutex
	HVs   map[uuid.UUID]*HV `json:"hvs"`
}

// The cloud struct that will be used by the rest of the app
var Cloud *HVList

// Initialize the cloud struct that will hold all of the hypervisors
func InitCloud() *HVList {
	Cloud = new(HVList)
	if err := getHVs(Cloud); err != nil {
		log.Fatal().Err(err).Msg("Failed to get HVs")
	} else {
		count := len(Cloud.HVs)
		log.Info().Int("hvs", count).Msg("Found hypervisors")
	}
	return Cloud
}
