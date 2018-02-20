package products_test

import (
	"testing"

	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/common/price"
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/shop/domain/products"
	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	testPrice, err := price.NewPrice(42, "USD")
	assert.NoError(t, err)

	testCases := []struct {
		TestName string

		ID          products.ID
		Name        string
		Description string
		Price       price.Price

		ExpectedErr error
	}{
		{
			TestName:    "valid",
			ID:          "1",
			Name:        "foo",
			Description: "bar",
			Price:       testPrice,
		},
		{
			TestName:    "empty_id",
			ID:          "",
			Name:        "foo",
			Description: "bar",
			Price:       testPrice,

			ExpectedErr: products.ErrEmptyID,
		},
		{
			TestName:    "empty_name",
			ID:          "1",
			Name:        "",
			Description: "bar",
			Price:       testPrice,

			ExpectedErr: products.ErrEmptyName,
		},
	}

	for _, c := range testCases {
		t.Run(c.TestName, func(t *testing.T) {
			_, err := products.NewProduct(c.ID, c.Name, c.Description, c.Price)
			assert.EqualValues(t, c.ExpectedErr, err)
		})
	}
}
