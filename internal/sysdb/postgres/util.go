package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Returnes values that are in the "a" but not in the "b"
func difference(a, b []int) []int {
	set := make(map[int]struct{}, len(b))
	for _, id := range b {
		set[id] = struct{}{}
	}

	var diff []int
	for _, id := range a {
		if _, found := set[id]; !found {
			diff = append(diff, id)
		}
	}
	return diff
}

func addRelationsTx(ctx context.Context, tx pgx.Tx, query string, entityID int, idsToAdd []int) error {

	if len(idsToAdd) != 0 {

		batch := &pgx.Batch{}

		for _, id := range idsToAdd {
			batch.Queue(query, entityID, id).Exec(func(ct pgconn.CommandTag) error {
				if ct.RowsAffected() == 0 {
					return fmt.Errorf("batch insert error: %w", ct.String())
				}
				return nil
			})
		}

		if err := tx.SendBatch(ctx, batch).Close(); err != nil {
			return fmt.Errorf("failed to add relation: %w", err)
		}
	}
	return nil
}
