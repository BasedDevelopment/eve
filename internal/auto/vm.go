package auto

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

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

const (
	Start uint8 = iota
	Reboot
	Poweroff
	Stop
	Reset
)

func stateStr(state uint8) string {
	switch state {
	case Start:
		return "start"
	case Reboot:
		return "reboot"
	case Poweroff:
		return "poweroff"
	case Stop:
		return "stop"
	case Reset:
		return "reset"
	}
	return ""
}

func (a *Auto) SetVMState(vmid string, state uint8) (respState models.VMState, err error) {
	c := a.getClient()
	reqUrl := a.Url + "/libvirt/domains/" + vmid + "/state"
	url, err := url.Parse(reqUrl)
	if err != nil {
		return
	}

	reqBody := map[string]string{
		"state": stateStr(state),
	}

	reqBodyBytes, err := json.Marshal(reqBody)

	req := http.Request{
		Method: "PATCH",
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		URL:  url,
		Body: ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes)),
	}

	resp, err := c.Do(&req)

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	err = json.Unmarshal(respBytes, &respState)

	return
}
