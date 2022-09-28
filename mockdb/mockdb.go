// Package contains a mock implementation of the storage interface. It is useful for testing.
package mockdb

import (
	"context"
	"errors"
	"sync"

	"github.com/irmatov/companies/filter"
	"github.com/irmatov/companies/types"
)

type mockStorage struct {
	mutex sync.Mutex
	data  []types.Company
	seq   int
}

type mockTx struct {
	data []types.Company
	seq  int
}

func New() types.Storage {
	return &mockStorage{}
}

func (m *mockStorage) Tx(ctx context.Context, action func(types.Tx) error) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tx := &mockTx{data: make([]types.Company, len(m.data)), seq: m.seq}
	copy(tx.data, m.data)
	err := action(tx)
	if err != nil {
		return err
	}
	m.data = tx.data
	m.seq = tx.seq
	return nil
}

func (tx *mockTx) Get(f filter.Filter) ([]types.Company, error) {
	switch f.Expr {
	case "":
		return tx.data, nil
	case "id = $1":
		id := f.Arguments[0].(int)
		for _, c := range tx.data {
			if c.Id == id {
				return []types.Company{c}, nil
			}
		}
	case "name = $1":
		name := f.Arguments[0].(string)
		for _, c := range tx.data {
			if c.Name == name {
				return []types.Company{c}, nil
			}
		}
	default:
		return nil, errors.New("not implemented")
	}
	return nil, nil
}

func (tx *mockTx) Create(newCompany types.Company) (int, error) {
	tx.seq++
	for _, c := range tx.data {
		if c.Name == newCompany.Name {
			return 0, errors.New("already exists")
		}
	}
	newCompany.Id = tx.seq
	tx.data = append(tx.data, newCompany)
	return tx.seq, nil
}

func (tx *mockTx) Update(c types.Company) error {
	for i, existing := range tx.data {
		if existing.Id == c.Id {
			tx.data[i] = c
			return nil
		}
	}
	return errors.New("not found")
}

func (tx *mockTx) Delete(id int) error {
	for i, c := range tx.data {
		if c.Id == id {
			tx.data = append(tx.data[:i], tx.data[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
