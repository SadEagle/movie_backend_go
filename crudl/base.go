package crudl

import (
	"database/sql"
	"errors"
	"fmt"
)

var EmptyDeletionError = errors.New("0 values was deleted")

func checkNonEmptyDeletion(res sql.Result) error {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("calculate affected rows by delete: %w", err)
	}
	if rowsAffected == 0 {
		return EmptyDeletionError
	}
	return nil
}
