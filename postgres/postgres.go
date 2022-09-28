// Implements storage interface using PostgreSQL database.
package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/irmatov/companies/filter"
	"github.com/irmatov/companies/types"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgres struct {
	db *sql.DB
}

type wrappedTx struct {
	tx *sql.Tx
}

// New creates a new storage instance using provided database.
func New(db *sql.DB) types.Storage {
	return &postgres{db}
}

// Tx creates a new transaction and executes the action function, automatically committing or rolling back the transaction.
func (p *postgres) Tx(ctx context.Context, action func(types.Tx) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = action(&wrappedTx{tx})
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Get returns a list of companies that match the given filter.
func (tx *wrappedTx) Get(f filter.Filter) ([]types.Company, error) {
	q := `SELECT id, name, code, country, website, phone FROM companies`
	if f.Expr != "" {
		q += ` WHERE ` + f.Expr
	}
	q += ` ORDER BY name`
	log.Printf("query: %s", q)
	rows, err := tx.tx.Query(q, f.Arguments...)
	if err != nil {
		return nil, err
	}
	companies := make([]types.Company, 0)
	for rows.Next() {
		var c types.Company
		err := rows.Scan(&c.Id, &c.Name, &c.Code, &c.Country, &c.Website, &c.Phone)
		if err != nil {
			rows.Close()
			return nil, err
		}
		companies = append(companies, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return companies, nil
}

// Create creates a new company and returns its ID.
func (tx *wrappedTx) Create(c types.Company) (int, error) {
	const q = `INSERT INTO companies (name, code, country, website, phone) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := tx.tx.QueryRow(q, c.Name, c.Code, c.Country, c.Website, c.Phone)
	if err := row.Scan(&c.Id); err != nil {
		return 0, err
	}
	return c.Id, nil
}

// Update updates the company with the given ID.
func (tx *wrappedTx) Update(c types.Company) error {
	_, err := tx.tx.Exec(`UPDATE companies SET name = $1, code = $2, country = $3, website = $4, phone = $5 WHERE id = $6`, c.Name, c.Code, c.Country, c.Website, c.Phone, c.Id)
	return err
}

// Delete deletes the company with the given ID.
func (tx *wrappedTx) Delete(id int) error {
	_, err := tx.tx.Exec(`DELETE FROM companies WHERE id = $1`, id)
	return err
}
