package auto

import (
	"encoding/json"

	"github.com/BasedDevelopment/auto/pkg/models"
)

func (a *Auto) GetHVSpecs() (hv models.HV, err error) {
	url := a.Url + "/libvirt"

	respBytes, status, err := a.httpReq("GET", url, nil)

	if err != nil {
		return hv, err
	}

	if status != 200 {
		return hv, err
	}

	// Unmarshal response
	err = json.Unmarshal(respBytes, &hv)

	return
}
