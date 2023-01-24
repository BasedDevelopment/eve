package auto

import (
	"encoding/json"
	"io/ioutil"

	"github.com/BasedDevelopment/auto/pkg/models"
)

func (a *Auto) GetLibvirtVMs() (vms []models.VM, err error) {
	c := a.getClient()
	url := a.Url + "/libvirt/domains"

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
	err = json.Unmarshal(respBytes, &vms)

	return
}

func (a *Auto) GetLibvirtVM(vmid string) (vm models.VM, err error) {
	c := a.getClient()
	url := a.Url + "/libvirt/domains/" + vmid

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
	err = json.Unmarshal(respBytes, &vm)

	return
}
