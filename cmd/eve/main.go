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

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/BasedDevelopment/eve/internal/config"
	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/BasedDevelopment/eve/internal/server"

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

	db.Init(config.Config.Database.URL)
}

func main() {
	// Creating the Cloud
	cloud := controllers.InitCloud()

	// Get HVs
	log.Info().Msg("Fetching HVs")

	for i := range cloud.HVs {
		hv := cloud.HVs[i]
		go connHV(hv)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    config.Config.API.Host + ":" + strconv.Itoa(config.Config.API.Port),
		Handler: server.Service(),
	}

	srvCtx, srvStopCtx := context.WithCancel(context.Background())

	// Watch for OS signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(srvCtx, 30*time.Second)
		defer shutdownCtxCancel() // release srvCtx if we take too long to shut down

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Warn().Msg("Graceful shutdown timed out... forcing regular exit.")
			}
		}()

		// Actually shutdown server, *gracefully*
		err := srv.Shutdown(shutdownCtx)

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to shutdown HTTP listener")
		}

		// Close database pool while we're at it
		db.Pool.Close()

		srvStopCtx()
	}()

	// Start the server
	err := srv.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Err(err).
			Msg("Failed to start HTTP listener")
	}

	// @todo Fix
	log.Info().
		Str("host", config.Config.API.Host).
		Int("port", config.Config.API.Port).
		Msg("Started HTTP server")

	// Wait for server context to be stopped
	<-srvCtx.Done()
	log.Info().Msg("Gracefully shutting down")
}

func connHV(hv *controllers.HV) {
	log.Info().
		Str("hostname", hv.Hostname).
		Msg("Connecting to HV and fetching VMs")

	// Connect to hypervisors, fetch VMs, and check for VM consistency
	if err := hv.Init(); err != nil {
		log.Warn().
			Err(err).
			Str("hostname", hv.Hostname).
			Msg("Failed to connect to HV and fetch VMs")
	} else {
		log.Info().
			Str("hostname", hv.Hostname).
			Str("hv", hv.Hostname).
			Int("vms", len(hv.VMs)).
			Msg("Connected to HV and fetched VMs")
	}
}
