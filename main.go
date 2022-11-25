package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ericzty/eve/lib/db"
	"github.com/ericzty/eve/lib/libvirt"
	"github.com/ericzty/eve/lib/routes"

	"github.com/go-chi/chi/v5"
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
	config = Config{}
)

func init() {
	stopwatch := time.Now()

	// Init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("+----------------------------------+")
	log.Info().Msg("|     EVE Virtural Environment     |")
	log.Info().Msg("|              v" + version + "              |")
	log.Info().Msg("+----------------------------------+")
	log.Info().Msg("Loading configuration")

	// Load configuration
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	} else {
		log.Info().Msg("Loaded configuration from " + configPath)
	}

	// Unmarshal configuration
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration")
	} else {
		log.Info().Msg("Parsed configuration")
	}

	// Check for required configuration
	if config.Name == "" {
		log.Fatal().Msg("Configuration: name is required")
	}
	if config.API.Host == "" {
		log.Fatal().Msg("Configuration: api.host is required")
	}
	if config.API.Port == 0 {
		log.Fatal().Msg("Configuration: api.port is required")
	}
	if config.Database.URL == "" {
		log.Fatal().Msg("Configuration: database.url is required")
	}

	// Init database
	log.Info().Msg("Connecting to database")
	if dberr := db.Init(k.String("database.url")); dberr != nil {
		log.Fatal().Err(dberr).Msg("Failed to connect to database")
	} else {
		log.Info().Msg("Connected to database")
	}

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

	log.Info().Msg("EVE database and libvirt initialized in " + time.Since(stopwatch).String())
}

func main() {
	// Init router
	r := chi.NewRouter()

	// Register routes
	r.Get("/health", routes.Health)

	// Start server
	listenAddress := k.String("api.host") + ":" + strconv.Itoa(k.Int("api.port"))
	log.Info().Msg("Starting the web server at " + listenAddress)
	if err := http.ListenAndServe(listenAddress, r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
