package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ericzty/eve/internal/config"
	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/db"
	"github.com/ericzty/eve/internal/libvirt"
	"github.com/ericzty/eve/internal/server"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	shutdownTimeout = 5 * time.Second
	version         = "0.0.1"
)

func init() {
	// Init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	log.Info().Msg("+-----------------------------------+")
	log.Info().Msg("|      EVE Virtual Environment      |")
	log.Info().Msg("|               v" + version + "               |")
	log.Info().Msg("+-----------------------------------+")

	// Load configuration
	log.Info().Msg("Loading configuration")

	config.Load()

	// Init database
	log.Info().Msg("Connecting to database")

	if err := db.Init(config.Config.Database.URL); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	} else {
		log.Info().Msg("Connected to database")
	}
}

func main() {
	// Get HVs
	log.Info().Msg("Getting HVs")

	hvs, hvErr := controllers.GetHVs()

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
	// Public routes

	// Start server
	listenAddress := config.Config.API.Host + ":" + strconv.Itoa(config.Config.API.Port)

	log.Info().Msg("Starting the web server at " + listenAddress)

	if err := http.ListenAndServe(listenAddress, server.Start()); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}