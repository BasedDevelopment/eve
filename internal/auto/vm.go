package auto

import (
	"encoding/json"
	"io/ioutil"

	"github.com/BasedDevelopment/auto/pkg/models"
)

func (a *Auto) GetLibvirtVMs() (vms []models.VM, err error) {
	c := a.getClient()
	url := a.Url + "/libvirt/domain"

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

func (a *Auto) GetVMState(vmid string) (state models.VMState, err error) {
	c := a.getClient()
	url := a.Url + "/libvirt/domains/" + vmid + "/state"

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
	err = json.Unmarshal(respBytes, &state)

	return
}
