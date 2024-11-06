package pg

import (
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors"
	"github.com/jackc/pgx/v5"
)

func wrapError(err error) error {
	if err == nil {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.RecordNotFound.WithCause(err)
	}
	// TODO handle more specific errors
	return errors.InternalDB.WithCause(err)
}
