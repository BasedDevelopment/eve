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
	log.Info().Msg("|               v" + version + "              |")
	log.Info().Msg("+-----------------------------------+")

	// Load configuration
	log.Info().Msg("Loading configuration")

	if err := config.Load(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Init database
	log.Info().Msg("Connecting to database")

	if err := db.Init(config.Config.Database.URL); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	} else {
		log.Info().Msg("Connected to database")
	}
}

func main() {
	// Creating the Cloud
	cloud := controllers.InitCloud()

	// Get HVs
	log.Info().Msg("Getting HVs")

	for i := range cloud.HVs {
		hv := cloud.HVs[i]

		log.Info().
			Str("hostname", hv.Hostname).
			Msg("Connecting to HV")

		if err := libvirt.InitHVs(cloud.HVs[i]); err != nil {
			log.Warn().
				Err(err).
				Str("hostname", hv.Hostname).
				Msg("Failed to connect to HV")
		} else {
			log.Info().
				Str("hostname", hv.Hostname).
				Msg("Connected to HV")
		}
	}

	// Report amount of online HVs
	var c int
	for i := range cloud.HVs {
		hv := cloud.HVs[i]
		if hv.Status == "Online" {
			c++
		}
	}

	// TODO: Report amount of VMs found

	log.Info().
		Int("online hypervisors", c).
		Int("total hypervisors", len(cloud.HVs))

	// Start server
	listenAddress := config.Config.API.Host + ":" + strconv.Itoa(config.Config.API.Port)

	log.Info().
		Str("host", config.Config.API.Host).
		Int("port", config.Config.API.Port).
		Msg("Started HTTP server")

	if err := http.ListenAndServe(listenAddress, server.Start()); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to start HTTP server")
	}
}
