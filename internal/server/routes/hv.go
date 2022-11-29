package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/rs/zerolog/log"
)

func GetHVs(w http.ResponseWriter, r *http.Request) {
	cloud := controllers.Cloud
	var out []*controllers.HV
	for _, hv := range cloud.HVs {
		hv.VMs = make(map[string]*controllers.VM)
		out = append(out, hv)
	}
	outJson, err := json.Marshal(out)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal cloud json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(outJson)
}

func GetHV(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	hvid := parts[3]
	hv := controllers.Cloud.HVs[hvid]
	outJson, err := json.Marshal(hv)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal hv json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(outJson)
}
