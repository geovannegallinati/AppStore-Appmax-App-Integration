//go:build appmax_live

package appmax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateCustomer(t *testing.T) {
	id := ensureCustomerID(t)
	assert.Greater(t, id, 0)
	t.Logf("Customer ready: id=%d", id)
}
