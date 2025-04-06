package inmemory

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"sync"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[int64]*domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int64]*domain.User),
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("user is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %d not found", id)
	}

	return user, nil
}
