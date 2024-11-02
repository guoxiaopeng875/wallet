package pg

import (
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors"
	"github.com/jackc/pgx/v5"
)

func wrapError(err error) *errors.Error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.RecordNotFound
	}
	// TODO handle more specific errors
	return errors.InternalDB
}
