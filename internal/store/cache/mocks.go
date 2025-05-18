package cache

import (
	"context"

	"github.com/yunsuk-jeung/social/internal/store"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (s *MockUserStore) Get(ctx context.Context, userID int64) (*store.User, error) { return nil, nil }

func (s *MockUserStore) Set(ctx context.Context, user *store.User) error { return nil }
