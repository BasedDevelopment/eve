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

package server

import (
	"time"

	"github.com/BasedDevelopment/eve/internal/server/middleware"
	"github.com/BasedDevelopment/eve/internal/server/routes"
	"github.com/BasedDevelopment/eve/internal/server/routes/admin"
	"github.com/BasedDevelopment/eve/internal/server/routes/users"
	"github.com/go-chi/chi/v5"
	cm "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func Service() *chi.Mux {
	r := chi.NewMux()

	// Middlewares
	r.Use(cm.RealIP)
	r.Use(cm.RequestID)
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(cm.Heartbeat("/health"))
	r.Use(cm.Recoverer)

	// Login
	r.Post("/login", routes.Login)

	// Admin endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Use(middleware.UserContext)
		r.Use(middleware.MustBeAdmin)

		// Hypervisor management
		r.Get("/admin/hypervisors", admin.GetHVs)
		r.Get("/admin/hypervisors/{id}", admin.GetHV)
		// r.Post("/admin/hypervisors/{id}", routes.CreateHV)
		// r.Patch("/admin/hypervisors/{id}", routes.UpdateHV)
		// r.Delete("/admin/hypervisors/{id}", routes.RemoveHV)

		// VM management
		r.Get("/admin/hypervisors/{id}/virtual_machines", admin.GetVMs)
		// r.Get("/admin/hypervisors/{id}/virtual_machines/{vmid}", routes.GetVM)
		// r.Post("/admin/hypervisors/{id}/virtual_machines", routes.CreateVM)
		// r.Patch("/admin/hypervisors/{id}/virtual_machines/{vmid}", routes.UpdateVM)
		// r.Delete("/admin/hypervisors/{id}/virtual_machines/{vmid}", routes.DeleteVM)

		// User management
		r.Post("/admin/users", admin.CreateUser)
		// r.Get("/admin/users", routes.GetUsers)
		// r.Patch("/admin/users/{id}", routes.UpdateUser)
		// r.Delete("/admin/users/{id}", routes.DeleteUser)
	})

	// User endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Use(middleware.UserContext)

		r.Get("/users/me", users.GetSelf)
		// r.Patch("/users/me", routes.UpdateUser)
		r.Get("/users/me/virtual_machines", users.GetVirtualMachines)
		// r.Get("/users/me/virtual_machines/{id}", routes.GetVM)
		// r.Patch("/users/me/virtual_machines/{id}", routes.UpdateVM)

		r.Post("/logout", routes.Logout)
	})

	return r
}
