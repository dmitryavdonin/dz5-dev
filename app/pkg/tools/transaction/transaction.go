package transaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Finish(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("error rollback transaction: %s", rollbackErr.Error())
		}
		return err
	} else {
		if commitErr := tx.Commit(ctx); commitErr != nil {
			return fmt.Errorf("error commit transaction: %s", commitErr.Error())
		}
		return nil
	}
}
