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
)

type Storage struct {
	Mutex      sync.Mutex `json:"-" db:"-"`
	ID         uuid.UUID  `json:"id"`
	HV         uuid.UUID  `json:"hv" db:"hv_id"`
	Enabeld    bool       `json:"enabled" db:"enabled"`
	Type       string     `json:"type" db:"type"`
	Path       string     `json:"path" db:"path"`
	Iso        bool       `json:"iso" db:"iso"`
	Disk       bool       `json:"disk" db:"disk"`
	CloudImage bool       `json:"cloud_image" db:"cloud_image"`
	Remarks    string     `json:"remarks" db:"remarks"`
}
