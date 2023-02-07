package auto

import (
	"encoding/json"
	"io/ioutil"

	"github.com/BasedDevelopment/auto/pkg/models"
)

func (a *Auto) GetHVSpecs() (hv models.HV, err error) {
	c := a.getHttpsClient()
	url := a.Url + "/libvirt"

	// Make request
	resp, err := c.Get(url)
	if err != nil {
		return
	}

	// Read response
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Unmarshal response
	err = json.Unmarshal(respBytes, &hv)

	return
}
