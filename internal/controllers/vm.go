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
	"net"
	"time"

	"github.com/google/uuid"
)

type VM struct {
	ID       uuid.UUID
	Hostname string
	User     *Profile `db:"-"`
	CPU      int
	Memory   int
	Nics     map[string]VMNic     `db:"-"`
	Storages map[string]VMStorage `db:"-"`
	Created  time.Time
	Updated  time.Time
	Remarks  string
	State    string `db:"-"`
}

type VMNic struct {
	ID      uuid.UUID
	name    string
	MAC     string
	IP      []net.IP `db:"-"`
	Created time.Time
	Updated time.Time
	Remarks string
	State   string `db:"-"`
}

type VMStorage struct {
	ID      uuid.UUID
	Size    int
	Created time.Time
	Updated time.Time
	Remarks string
}
