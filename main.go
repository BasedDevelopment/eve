package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ericzty/eve/internal/db"
	"github.com/ericzty/eve/internal/libvirt"
	"github.com/ericzty/eve/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	configPath      = "/etc/eve/config.toml"
	shutdownTimeout = 5 * time.Second
	version         = "0.0.1"
)

var (
	k      = koanf.New(".")
	parser = toml.Parser()
)

func init() {
	// Init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("+-----------------------------------+")
	log.Info().Msg("|      EVE Virtual Environment      |")
	log.Info().Msg("|               v" + version + "               |")
	log.Info().Msg("+-----------------------------------+")
	log.Info().Msg("Loading configuration")

	// Load configuration
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	} else {
		log.Info().Msg("Loaded configuration from " + configPath)
	}

	// Check config
	log.Info().Msg("Checking config")
	if err := checkConfig(); err != nil {
		log.Fatal().Err(err).Msg("Invalid configuration")
	} else {
		log.Info().Msg("Configuration is valid")
	}

	// Init database
	log.Info().Msg("Connecting to database")
	if dberr := db.Init(k.String("database.url")); dberr != nil {
		log.Fatal().Err(dberr).Msg("Failed to connect to database")
	} else {
		log.Info().Msg("Connected to database")
	}

}

func main() {
	// Get HVs
	log.Info().Msg("Getting HVs")
	hvs, hvErr := db.GetHVs()
	if hvErr != nil {
		log.Fatal().Err(hvErr).Msg("Failed to get HVs")
	} else {
		hvCount := len(hvs)
		log.Info().Msg("Got " + strconv.Itoa(hvCount) + " HVs")
	}

	// Init libvirt
	log.Info().Msg("Connecting to HVs via libvirt")
	for i := range hvs {
		hv := hvs[i]
		log.Info().Msg("Connecting to " + hv.Hostname)
		if err := libvirt.Init(&hvs[i]); err != nil {
			log.Warn().Err(err).Msg("Failed to connect to HV " + hv.Hostname)
		} else {
			hv := hvs[i]
			log.Info().Msg("Connected to " + hv.Hostname + ",libvirt version: " + hv.Version)
		}
	}
	// Report amount of online HVs
	var c int
	for i := range hvs {
		hv := hvs[i]
		if hv.Status == "online" {
			c++
		}
	}
	log.Info().Msg("Online HVs: " + strconv.Itoa(c) + "/" + strconv.Itoa(len(hvs)))
	// Init router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	// Rate limiter
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// Public routes
	r.Get("/health", routes.Health)
	r.Post("/login", routes.Login)
	/*

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(routes.AdminAuth)
			r.Get("/admin/hvs", routes.GetHVs)
			r.Get("/admin/hvs/{id}", routes.GetHV)
			r.Get("/admin/hvs/{id}/vms", routes.GetVMs)
			r.Get("/admin/hvs/{id}/vms/{vmid}", routes.GetVM)
			r.Post("/admin/hvs/{id}/vms", routes.CreateVM)
			r.Put("/admin/hvs/{id}/vms/{vmid}", routes.UpdateVM)
			r.Delete("/admin/hvs/{id}/vms/{vmid}", routes.DeleteVM)
			r.Post("/admin/users", routes.CreateUser)
			r.Get("/admin/users", routes.GetUsers)
			r.Get("/admin/users/{id}", routes.GetUser)
			r.Put("/admin/users/{id}", routes.UpdateUser)
			r.Delete("/admin/users/{id}", routes.DeleteUser)
		})

		// User routes
		r.Group(func(r chi.Router) {
			r.Use(routes.UserAuth)
			r.Get("/users/me", routes.GetUser)
			r.Put("/users/me", routes.UpdateUser)
			r.Get("/users/me/vms", routes.GetVMs)
			r.Get("/users/me/vms/{id}", routes.GetVM)
		})
	*/
	// Start server
	listenAddress := k.String("api.host") + ":" + strconv.Itoa(k.Int("api.port"))
	log.Info().Msg("Starting the web server at " + listenAddress)
	if err := http.ListenAndServe(listenAddress, r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
