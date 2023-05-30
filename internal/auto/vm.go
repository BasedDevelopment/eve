package auto

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BasedDevelopment/auto/pkg/models"
	"github.com/BasedDevelopment/eve/internal/util"
	eStatus "github.com/BasedDevelopment/eve/pkg/status"
	"github.com/google/uuid"
)

func (a *Auto) GetLibvirtVMs() (vms []models.VM, err error) {
	url := a.Url + "/libvirt/domains"
	respBytes, status, err := a.httpReq("GET", url, nil)

	if err != nil {
		return
	}

	if status != http.StatusOK {
		respBody := string(respBytes)
		return nil, fmt.Errorf("status code %d: %s", status, respBody)
	}

	err = json.Unmarshal(respBytes, &vms)

	return
}

func (a *Auto) GetLibvirtVM(vmid string) (vm models.VM, err error) {
	url := a.Url + "/libvirt/domain/" + vmid
	respBytes, status, err := a.httpReq("GET", url, nil)

	if err != nil {
		return
	}

	if status != http.StatusOK {
		respBody := string(respBytes)
		return vm, fmt.Errorf("status code %d: %s", status, respBody)
	}

	err = json.Unmarshal(respBytes, &vm)

	return
}

func (a *Auto) GetVMState(vmid string) (state models.VMState, err error) {
	url := a.Url + "/libvirt/domains/" + vmid + "/state"
	respBytes, status, err := a.httpReq("GET", url, nil)

	if err != nil {
		return
	}

	if status != http.StatusOK {
		respBody := string(respBytes)
		return state, fmt.Errorf("status code %d: %s", status, respBody)
	}

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
	reqUrl := a.Url + "/libvirt/domains/" + vmid + "/state"
	reqBody := map[string]string{
		"state": stateStr(state),
	}

	respBytes, status, err := a.httpReq("POST", reqUrl, reqBody)

	if err != nil {
		return
	}

	if status != http.StatusOK {
		respBody := string(respBytes)
		return respState, fmt.Errorf("status code %d: %s", status, respBody)
	}

	err = json.Unmarshal(respBytes, &respState)

	if respState.State == eStatus.StatusUnknown {
		respStr := string(respBytes)
		return respState, fmt.Errorf(respStr)
	}

	return
}

func (a *Auto) CreateVM(req *util.VMCreateRequest, vmid uuid.UUID) (err error) {
	reqUrl := a.Url + "/libvirt/domain/" + vmid.String()
	respBytes, status, err := a.httpReq("PUT", reqUrl, req)

	if err != nil {
		return
	}

	if status != http.StatusCreated {
		respBody := string(respBytes)
		return fmt.Errorf("status code %d: %s", status, respBody)
	}

	return
}

func (a *Auto) DeleteVM(vmid string) error {
	reqUrl := a.Url + "/libvirt/domain/" + vmid
	respBytes, status, err := a.httpReq("DELETE", reqUrl, nil)

	if err != nil {
		return err
	}

	if status != http.StatusOK {
		respBody := string(respBytes)
		return fmt.Errorf("status code %d: %s", status, respBody)
	}

	return nil
}
