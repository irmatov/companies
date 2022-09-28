package service

import (
	"context"
	"testing"

	"github.com/irmatov/companies/filter"
	"github.com/irmatov/companies/mockdb"
	"github.com/irmatov/companies/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompanies(t *testing.T) {
	svc := New(mockdb.New())
	ctx := context.Background()

	// ensure there are no companies
	cs, err := svc.Get(ctx, filter.Filter{})
	require.NoError(t, err)
	assert.Empty(t, cs)

	// add a company
	c1 := types.Company{
		Name:    "Hydrogen",
		Code:    "HYDRO",
		Country: "US",
		Website: "http://hydrogen.io",
		Phone:   "+1234567890",
	}
	c1.Id, err = svc.Create(ctx, c1)
	require.NoError(t, err)

	// ensure it is there
	got, err := svc.Get(ctx, filter.Filter{Expr: "id = $1", Arguments: []interface{}{c1.Id}})
	require.NoError(t, err)
	require.Equal(t, []types.Company{c1}, got)

	// adding the same company twice gives no error
	c1.Id, err = svc.Create(ctx, c1)
	require.NoError(t, err)

	// add another company
	c2 := types.Company{
		Name:    "Beryllium",
		Code:    "BERYL",
		Country: "FR",
		Website: "http://beryllium.io",
		Phone:   "+987654321",
	}
	c2.Id, err = svc.Create(ctx, c2)
	require.NoError(t, err)

	// ensure it is there
	got, err = svc.Get(ctx, filter.Filter{Expr: "id = $1", Arguments: []interface{}{c2.Id}})
	require.NoError(t, err)
	require.Equal(t, []types.Company{c2}, got)

	// list all companies
	got, err = svc.Get(ctx, filter.Filter{})
	require.NoError(t, err)
	require.Equal(t, []types.Company{c1, c2}, got)

	// update the first company
	c1.Website = "http://example.com/"
	require.NoError(t, svc.Update(ctx, c1))

	// list all companies, ensure the first one is updated
	got, err = svc.Get(ctx, filter.Filter{})
	require.NoError(t, err)
	require.Equal(t, []types.Company{c1, c2}, got)

	// drop the second company
	require.NoError(t, svc.Delete(ctx, c2.Id))

	// list all companies, ensure the second is missing
	got, err = svc.Get(ctx, filter.Filter{})
	require.NoError(t, err)
	require.Equal(t, []types.Company{c1}, got)

	// attempt to delete the same company again results in an error
	require.Equal(t, types.ErrNotFound, svc.Delete(ctx, c2.Id))

	// list all companies
	got, err = svc.Get(ctx, filter.Filter{})
	require.NoError(t, err)
	require.Equal(t, []types.Company{c1}, got)
}
