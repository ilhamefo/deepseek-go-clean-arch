package domain

import (
	"context"
)

type UserMeilisearchRepository interface {
	SetupIndex() (err error)
	SeedIndex() error
	Search(ctx context.Context, keyword string) (users []*UserVCC, err error)
	Update(ctx context.Context, user *UserVCC) error
	CheckHealth() error
}
