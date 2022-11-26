package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func GetHVs() (HVs []HV, err error) {
	/*
		inputRows := [][]interface{}{
			{"85833bb8-2f0a-4b1e-981f-f9cb3597904c", "dorm0.sit.eric.si", "10.10.9.4", 16509, "sit"},
			{"70ada0c5-b641-4633-88bc-2d58c8a387a5", "broke.sit.eric.si", "10.10.222.1", 16509, "sit"},
		}

		copyCount, err := pool.CopyFrom(context.Background(), pgx.Identifier{"hv"}, []string{"id", "hostname", "ip", "port", "site"}, pgx.CopyFromRows(inputRows))
		if err != nil {
			fmt.Println("Error copying from rows: ", err)
		}
		if int(copyCount) != len(inputRows) {
			fmt.Println("Len mismatch", copyCount, len(inputRows))
		} else {
			fmt.Println("Copied", copyCount, "rows")
		}
	*/

	// Reading HVs
	rows, queryErr := pool.Query(context.Background(), "SELECT * FROM hv")
	if queryErr != nil {
		return nil, fmt.Errorf("Error reading hv: %w", queryErr)
	}
	HVs, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[HV])
	if collectErr != nil {
		return nil, fmt.Errorf("Error collecting hv: %w", collectErr)
	}
	return
}
