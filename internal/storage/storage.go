package storage

import (
	"github.com/google/uuid"
	"github.com/kperanovic/epam-systems/api/v1/types"
)

type Storage interface {
	Connect() error
	SaveCompany(*types.Company) error
	GetCompany(id uuid.UUID) (*types.Company, error)
	UpdateCompany(uuid.UUID, *types.Company) error
	DeleteCompany(id uuid.UUID) error
}
