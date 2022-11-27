package controllers

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type Cloud struct {
	Mutex sync.Mutex
	HVs   map[string]*HV
}

var C *Cloud

func InitCloud() *Cloud {
	C = new(Cloud)
	if err := getHVs(C); err != nil {
		log.Fatal().Err(err).Msg("Failed to get HVs")
	} else {
		count := len(C.HVs)
		log.Info().Msgf("Found %d HVs", count)
	}
	return C
}
