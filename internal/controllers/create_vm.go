package controllers

import (
	"context"

	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/google/uuid"
)

func (hv *HV) CreateVM(ctx context.Context, vm *util.VMCreateRequest, hvid uuid.UUID) (uuid.UUID, error) {

	vmid, err := uuid.NewRandom()
	if err != nil {
		return vmid, err
	}

	err = hv.Auto.CreateVM(vm, vmid)
	if err != nil {
		return vmid, err
	}

	//don't use vm.Id here, it's not set

	_, err = db.Pool.Exec(
		ctx,
		"INSERT INTO vm (id, hv_id, hostname, profile_id, cpu, memory) VALUES ($1, $2, $3, $4, $5, $6)",
		vmid,
		hvid,
		vm.Hostname,
		vm.User,
		vm.CPU,
		vm.Memory,
	)

	if err != nil {
		return vmid, err
	}

	//refresh hv's vm list
	_, err = hv.getVMsFromDB()

	return vmid, err
}
