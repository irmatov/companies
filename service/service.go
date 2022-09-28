// Contains the implementation of Companies service
package service

import (
	"context"

	"github.com/irmatov/companies/filter"
	"github.com/irmatov/companies/types"
)

// Companies provides an interface to perform various operations on companies.
type Companies struct {
	storage types.Storage
}

// New creates a new instance of Companies service using a provided storage implementation.
func New(storage types.Storage) *Companies {
	return &Companies{storage}
}

// Get returns a list of companies that match the provided filter.
func (c *Companies) Get(ctx context.Context, f filter.Filter) ([]types.Company, error) {
	var companies []types.Company
	var err error
	err = c.storage.Tx(ctx, func(tx types.Tx) error {
		companies, err = tx.Get(f)
		return err
	})
	return companies, err
}

// Create creates a new company and returns its ID.
func (c *Companies) Create(ctx context.Context, company types.Company) (int, error) {
	var id int
	err := c.storage.Tx(ctx, func(tx types.Tx) error {
		existing, err := tx.Get(filter.Filter{Expr: "name = $1", Arguments: []interface{}{company.Name}})
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			company.Id = existing[0].Id
			if company == existing[0] {
				id = company.Id
				return nil
			}
			return types.ErrAlreadyExists
		}
		id, err = tx.Create(company)
		return err
	})
	return id, err
}

// Update updates an existing company.
func (c *Companies) Update(ctx context.Context, company types.Company) error {
	err := c.storage.Tx(ctx, func(tx types.Tx) error {
		existing, err := tx.Get(filter.Filter{Expr: "id = $1", Arguments: []interface{}{company.Id}})
		if err != nil {
			return err
		}
		if len(existing) == 0 {
			return types.ErrNotFound
		}
		if company == existing[0] {
			return nil
		}
		return tx.Update(company)
	})
	return err
}

// Delete deletes an existing company.
func (c *Companies) Delete(ctx context.Context, id int) error {
	err := c.storage.Tx(ctx, func(tx types.Tx) error {
		existing, err := tx.Get(filter.Filter{Expr: "id = $1", Arguments: []interface{}{id}})
		if err != nil {
			return err
		}
		if len(existing) == 0 {
			return types.ErrNotFound
		}
		return tx.Delete(id)
	})
	return err
}
