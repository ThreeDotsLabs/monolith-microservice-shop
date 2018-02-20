package orders_test

import (
	"testing"

	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/orders/domain/orders"
	"github.com/stretchr/testify/assert"
)

func TestNewAddress(t *testing.T) {
	testCases := []struct {
		TestName string

		Name     string
		Street   string
		City     string
		PostCode string
		Country  string

		ExpectedErr bool
	}{
		{
			TestName:    "valid",
			Name:        "test",
			Street:      "test",
			City:        "test",
			PostCode:    "test",
			Country:     "test",
			ExpectedErr: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.TestName, func(t *testing.T) {
			address, err := orders.NewAddress(c.Name, c.Street, c.City, c.PostCode, c.Country)

			if c.ExpectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, c.Name, address.Name())
				assert.EqualValues(t, c.Street, address.Street())
				assert.EqualValues(t, c.City, address.City())
				assert.EqualValues(t, c.PostCode, address.PostCode())
				assert.EqualValues(t, c.Country, address.Country())
			}
		})
	}
}
