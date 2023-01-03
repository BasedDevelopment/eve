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

		r.Route("/admin", func(r chi.Router) {
			// Hypervisor management
			r.Route("/hypervisors", func(r chi.Router) {
				r.Get("/", admin.GetHVs)
				//r.Post("/", admin.CreateHV)
				r.Route("/{hypervisor}", func(r chi.Router) {
					r.Get("/", admin.GetHV)
					//r.Patch("/", admin.UpdateHV)
					//r.Delete("/", admin.DeleteHV)
					r.Route("/virtual_machines", func(r chi.Router) {
						r.Get("/", admin.GetVMs)
						//r.Post("/", admin.CreateVM)
						r.Route("/{virtual_machine}", func(r chi.Router) {
							r.Get("/", admin.GetVM)
							//r.Patch("/", admin.UpdateVM)
							//r.Delete("/", admin.DeleteVM)
						})
					})
				})
			})
			r.Route("/users", func(r chi.Router) {
				r.Post("/", admin.CreateUser)
				//r.Get("/", admin.GetUsers)
				//r.Route("/{user}", func(r chi.Router) {
				//r.Get("/", admin.GetUser)
				//r.Patch("/", admin.UpdateUser)
				//r.Delete("/", admin.DeleteUser)
			})
		})
	})

	// User endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Use(middleware.UserContext)

		r.Route("/users", func(r chi.Router) {
			r.Get("/me", users.GetSelf)
			//r.Patch("/me", users.UpdateSelf)
			r.Route("/virtual_machines", func(r chi.Router) {
				r.Get("/", users.GetVMs)
				r.Route("/{virtual_machines}", func(r chi.Router) {
					//r.Get("/", users.GetVirtualMachine)
					//r.Patch("/", users.UpdateVirtualMachine)
				})
			})
		})

		r.Post("/logout", routes.Logout)
	})

	return r
}
