package types

import (
	"context"

	"github.com/irmatov/companies/filter"
)

type Company struct {
	Id      int
	Name    string
	Code    string
	Country string
	Website string
	Phone   string
}

type Storage interface {
	Tx(ctx context.Context, action func(Tx) error) error
}

type Tx interface {
	Get(f filter.Filter) ([]Company, error)
	Create(c Company) (int, error)
	Update(c Company) error
	Delete(id int) error
}

type CompanyService interface {
	Get(ctx context.Context, f filter.Filter) ([]Company, error)
	Create(ctx context.Context, c Company) (int, error)
	Update(ctx context.Context, c Company) error
	Delete(ctx context.Context, id int) error
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrAlreadyExists = Error("already exists")
	ErrNotFound      = Error("not found")
)
