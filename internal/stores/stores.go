package stores

import (
	"context"
	"errors"
)

type Store interface {
	Upload(ctx context.Context, sourceRepo string, objects ...string) error
	Retrieve(ctx context.Context, sourceRepo string, objects ...string) error
	Options() Options
}

var (
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)
