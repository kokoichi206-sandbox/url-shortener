package usecase

import "context"

type database interface {
	Health(ctx context.Context) error
}
