package server

import (
	"time"

	"github.com/ericzty/eve/internal/server/middlewares"
	"github.com/ericzty/eve/internal/server/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func Start() *chi.Mux {
	r := chi.NewMux()

	// Middlewares
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middlewares.Logger)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(middleware.Heartbeat("/health"))

	// Login
	r.Post("/login", routes.Login)

	// Admin endpoints
	r.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)
		r.Use(middlewares.MustBeAdmin)

		r.Get("/admin/health", routes.Health)
		r.Get("/admin/hvs", routes.GetHVs)
		r.Get("/admin/hv/{id}", routes.GetHV)
		// r.Post("/admin/hvs/{id}", routes.CreateHV)
		// r.Put("/admin/hvs/{id}", routes.UpdateHV)
		// r.Delete("/admin/hvs/{id}", routes.RemoveHV)
		// r.Get("/admin/hvs/{id}/vms", routes.GetVMs)
		// r.Get("/admin/hvs/{id}/vms/{vmid}", routes.GetVM)
		// r.Post("/admin/hvs/{id}/vms", routes.CreateVM)
		// r.Put("/admin/hvs/{id}/vms/{vmid}", routes.UpdateVM)
		// r.Delete("/admin/hvs/{id}/vms/{vmid}", routes.DeleteVM)
		r.Post("/admin/users", routes.CreateUser)
		// r.Get("/admin/users", routes.GetUsers)
		// r.Get("/admin/users/{id}", routes.GetUser)
		// r.Put("/admin/users/{id}", routes.UpdateUser)
		// r.Delete("/admin/users/{id}", routes.DeleteUser)
	})

	// User endpoints
	r.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Post("/logout", routes.Logout)
		r.Get("/user/health", routes.Health)
		// r.Get("/users/me", routes.GetUser)
		// r.Put("/users/me", routes.UpdateUser)
		// r.Get("/users/me/vms", routes.GetVMs)
		// r.Get("/users/me/vms/{id}", routes.GetVM)
		// r.Put("/users/me/vms/{id}", routes.UpdateVM)
	})

	return r
}
