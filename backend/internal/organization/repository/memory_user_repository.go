package repository

//mocked data

import (
	"errors"
	"strings"
	"sync"


	"api-on/internal/auth/model"
	"github.com/google/uuid"
)

type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*model.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[uuid.UUID]*model.User),
	}
}

func (r *MemoryUserRepository) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.users {
		if strings.EqualFold(existing.Email, user.Email) {
			return errors.New("email already exists")
		}
	}

	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) FindByEmail(email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if strings.EqualFold(user.Email, email) {
			copyUser := *user
			return &copyUser, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *MemoryUserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	copyUser := *user
	return &copyUser, nil
}

func (r *MemoryUserRepository) List() ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, *user)
	}

	return users, nil
}

func (r *MemoryUserRepository) Update(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.users[user.ID]
	if !ok {
		return errors.New("user not found")
	}

	for id, existing := range r.users {
		if id != user.ID && strings.EqualFold(existing.Email, user.Email) {
			return errors.New("email already exists")
		}
	}

	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}