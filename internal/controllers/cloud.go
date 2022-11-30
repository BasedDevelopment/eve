package controllers

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type HVList struct {
	mutex sync.Mutex
	HVs   map[string]*HV `json:"hvs"`
}

var Cloud *HVList

func InitCloud() *HVList {
	Cloud = new(HVList)
	if err := getHVs(Cloud); err != nil {
		log.Fatal().Err(err).Msg("Failed to get HVs")
	} else {
		count := len(Cloud.HVs)
		log.Info().Int("hv", count).Msg("Found hypervisors")
	}
	return Cloud
}
