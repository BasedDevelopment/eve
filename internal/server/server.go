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
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(cm.Heartbeat("/health"))
	r.Use(cm.Recoverer)

	// Login
	r.Post("/login", routes.Login)
	r.Post("/health", routes.Health)

	// Admin endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Use(middleware.UserContext)
		r.Use(middleware.MustBeAdmin)

		// Hypervisor management
		r.Get("/admin/hvs", admin.GetHVs)
		r.Get("/admin/hv/{id}", admin.GetHV)
		// r.Post("/admin/hvs/{id}", routes.CreateHV)
		// r.Patch("/admin/hvs/{id}", routes.UpdateHV)
		// r.Delete("/admin/hvs/{id}", routes.RemoveHV)

		// VM management
		// r.Get("/admin/hvs/{id}/vms", routes.GetVMs)
		// r.Get("/admin/hvs/{id}/vms/{vmid}", routes.GetVM)
		// r.Post("/admin/hvs/{id}/vms", routes.CreateVM)
		// r.Patch("/admin/hvs/{id}/vms/{vmid}", routes.UpdateVM)
		// r.Delete("/admin/hvs/{id}/vms/{vmid}", routes.DeleteVM)

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
		// r.Get("/users/me/vms", routes.GetVMs)
		// r.Get("/users/me/vms/{id}", routes.GetVM)
		// r.Patch("/users/me/vms/{id}", routes.UpdateVM)

		r.Post("/logout", routes.Logout)
	})

	return r
}
