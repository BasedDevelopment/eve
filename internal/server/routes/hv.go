package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/rs/zerolog/log"
)

func GetHVs(w http.ResponseWriter, r *http.Request) {
	cloud := controllers.Cloud
	for _, hv := range cloud.HVs {
		hv.VMs = nil
	}
	cloudJson, err := json.Marshal(cloud)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal cloud json")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(cloudJson)
}
