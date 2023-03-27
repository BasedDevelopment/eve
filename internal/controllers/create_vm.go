package controllers

import (
	"context"

	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (hv *HV) CreateVM(ctx context.Context, vm *util.VMCreateRequest, hvid uuid.UUID) (uuid.UUID, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return uuid, err
	}

	idStr := uuid.String()

	// Pick connection from pool
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return uuid, err
	}
	defer conn.Release()

	// Prepare the db entry
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	_, err = tx.Exec(ctx, idStr, "PREPARE TRANSACTION 'INSERT INTO vm (id, hv_id, hostname, profile_id, cpu, memory) VALUES ($1, $2, $3, $4, $5, $6)';",
		vm.Id,
		hvid,
		vm.Hostname,
		vm.User,
		vm.CPU,
		vm.Memory,
	)
	if err != nil {
		return uuid, err
	}

	// Call auto
	err = hv.Auto.CreateVM(vm)
	if err != nil {
		tx.Rollback(ctx)
		return uuid, err
	}

	tx.Commit(ctx)
	return uuid, nil
}
