package server

import (
	"time"

	middleware "github.com/ericzty/eve/internal/server/middleware"
	"github.com/ericzty/eve/internal/server/routes"
	"github.com/ericzty/eve/internal/server/routes/admin"
	"github.com/ericzty/eve/internal/server/routes/users"
	"github.com/go-chi/chi/v5"
	cm "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func Start() *chi.Mux {
	r := chi.NewMux()

	// Middlewares
	r.Use(cm.RealIP)
	r.Use(cm.RequestID)
	r.Use(cm.Logger)
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
		// r.Get("/admin/hypervisors/{id}/virtual_machines", routes.GetVMs)
		// r.Get("/admin/hypervisors/{id}/virtual_machines/{vmid}", routes.GetVM)
		// r.Post("/admin/hypervisors/{id}/virtual_machines", routes.CreateVM)
		// r.Patch("/admin/hypervisors/{id}/virtual_machines/{vmid}", routes.UpdateVM)
		// r.Delete("/admin/hypervisors/{id}/virtual_machines/{vmid}", routes.DeleteVM)

		// User management
		r.Post("/admin/users", admin.CreateUser)
		// r.Get("/admin/users", routes.GetUsers)
		// r.Get("/admin/users/{id}", routes.GetUser)
		// r.Patch("/admin/users/{id}", routes.UpdateUser)
		// r.Delete("/admin/users/{id}", routes.DeleteUser)
	})

	// User endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Use(middleware.UserContext)

		r.Get("/users/me", users.GetSelf)
		// r.Patch("/users/me", routes.UpdateUser)
		// r.Get("/users/me/virtual_machines", routes.GetVMs)
		// r.Get("/users/me/virtual_machines/{id}", routes.GetVM)
		// r.Patch("/users/me/virtual_machines/{id}", routes.UpdateVM)

		r.Post("/logout", routes.Logout)
	})

	return r
}
