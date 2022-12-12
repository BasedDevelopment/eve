package admin

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
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
		util.WriteError(err, w, http.StatusInternalServerError)
		return
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
		util.WriteError(err, w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(outJson)
}
