package qry

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func Transaction(db Querier, f func(tx pgx.Tx) error) (err error) {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed db.Begin: %s", err)
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			err = fmt.Errorf("failed tx.Rollback: %s", err)
		}
	}()

	if err = f(tx); err != nil {
		return fmt.Errorf("failed f: %s", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed tx.Commit: %s", err)
	}

	return nil
}
