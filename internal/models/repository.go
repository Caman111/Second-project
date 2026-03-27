package models

import "sync"

type Repository struct {
	data   map[uint]BiznessCreate
	nextID uint
	mu     sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		data:   make(map[uint]BiznessCreate),
		nextID: 1,
	}
}

func (r *Repository) Create(b BiznessCreate) BiznessCreate {
	r.mu.Lock()
	defer r.mu.Unlock()

	b.ID = r.nextID

	r.data[b.ID] = b

	r.nextID++

	return b
}
