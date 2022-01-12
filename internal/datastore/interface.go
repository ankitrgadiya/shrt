package datastore

import (
	"context"

	"argc.in/shrt/internal/model"
)

type RouteStore interface {
	Close() error
	Save(ctx context.Context, r *model.Route) error
	Query(ctx context.Context, r *model.Route) error
	QueryAll(ctx context.Context) ([]model.Route, error)
	Delete(ctx context.Context, r *model.Route) error
}
