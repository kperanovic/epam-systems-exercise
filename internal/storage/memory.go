package storage

import (
	"github.com/google/uuid"
	"github.com/kperanovic/epam-systems/api/v1/types"
)

type memoryStorage struct {
	store map[uuid.UUID]*types.Company
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		store: make(map[uuid.UUID]*types.Company, 0),
	}
}

func (mem *memoryStorage) Connect() error {
	return nil
}

func (mem *memoryStorage) SaveCompany(company *types.Company) error {
	mem.store[company.ID] = company

	return nil
}

func (mem *memoryStorage) GetCompany(id uuid.UUID) (*types.Company, error) {
	return mem.store[id], nil
}

func (mem *memoryStorage) UpdateCompany(id uuid.UUID, company *types.Company) error {
	mem.store[id] = company

	return nil
}

func (mem *memoryStorage) DeleteCompany(id uuid.UUID) error {
	delete(mem.store, id)

	return nil
}
