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

package util

import (
	"encoding/json"
	"net/http"

	"github.com/BasedDevelopment/eve/pkg/status"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type Validatable[T any] interface {
	Validate() error
	*T
}

type Request interface {
	UserCreateRequest |
		LoginRequest |
		SetStateRequest |
		VMCreateRequest
}

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Disabled bool   `json:"disabled"`
	IsAdmin  bool   `json:"is_admin"`
	Remarks  string `json:"remarks"`
}

func (s UserCreateRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required, validation.Length(2, 20)),
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Password, validation.Required, validation.Length(8, 0), is.PrintableASCII), // todo: PrintableASCII includes spaces, password shouldn't
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s LoginRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Password, validation.Required, validation.Length(8, 0), is.PrintableASCII), // todo: PrintableASCII includes spaces, password shouldn't
	)
}

type SetStateRequest struct {
	State string `json:"state"`
}

func (s SetStateRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.State, validation.Required, validation.In("start", "reboot", "poweroff", "stop", "reset")),
	)
}

type VMCreateRequest struct {
	Id        uuid.UUID     `json:"id"`
	User      uuid.UUID     `json:"user"`
	Hostname  string        `json:"hostname"`
	CPU       int           `json:"cpu"`
	Memory    int           `json:"memory"`
	State     status.Status `json:"state"`
	Image     string        `json:"image"`
	Cloud     bool          `json:"cloud"`
	OS        string        `json:"os"`
	OSVariant string        `json:"os_variant"`
	userData  string        `json:"userData"`
	metaData  string        `json:"metaData"`
	Disk      []struct {
		Id   int    `json:"id"`
		Size int    `json:"size"`
		Disk string `json:"disk"`
	} `json:"disk"`
	Iface []struct {
		Bridge string `json:"bridge"`
		MAC    string `json:"mac"`
	} `json:"iface"`
}

func (s VMCreateRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Hostname, validation.Required, validation.Length(2, 20)),
		validation.Field(&s.CPU, validation.Required, validation.Min(1)),
		validation.Field(&s.Memory, validation.Required, validation.Min(1)),
		//TODO
	)
}

func ParseRequest[R Request, T Validatable[R]](r *http.Request, rq T) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&rq); err != nil {
		return err
	}

	return rq.Validate()
}
