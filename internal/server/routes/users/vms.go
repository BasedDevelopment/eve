package users

import (
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
	"github.com/google/uuid"
)

func GetVirtualMachines(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile := controllers.Profile{ID: ctx.Value("owner").(uuid.UUID)}
	profile, err := profile.Get(ctx)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Fetch users virtual machines
	cloud := controllers.Cloud
	var hvs []map[string]interface{}

	for _, hv := range cloud.HVs {
		hv.VMs = make(map[uuid.UUID]*controllers.VM)
		hvs = append(hvs, map[string]interface{}{
			"hv":  hv.Hostname,
			"vms": hv.VMs,
		})
	}

	// Send response
	if err := util.WriteResponse(hvs, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}

}
