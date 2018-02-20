package products_test

import (
	"testing"

	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/common/price"
	products_domain "github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/shop/domain/products"
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/shop/infrastructure/products"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRepository(t *testing.T) {
	repo := products.NewMemoryRepository()

	assertAllProducts(t, repo, []products_domain.Product{})

	product1 := addProduct(t, repo, "1")
	// test idempotency
	_ = addProduct(t, repo, "1")

	assertAllProducts(t, repo, []products_domain.Product{*product1})
	repoProduct1, err := repo.ByID("1")
	assert.NoError(t, err)
	assert.EqualValues(t, *product1, *repoProduct1)

	product2 := addProduct(t, repo, "2")

	assertAllProducts(t, repo, []products_domain.Product{*product1, *product2})
	repoProduct2, err := repo.ByID("2")
	assert.NoError(t, err)
	assert.EqualValues(t, *product2, *repoProduct2)
}

func assertAllProducts(t *testing.T, repo *products.MemoryRepository, expectedProducts []products_domain.Product) {
	products, err := repo.AllProducts()

	assert.EqualValues(t, expectedProducts, products)
	assert.NoError(t, err)
}

func addProduct(t *testing.T, repo *products.MemoryRepository, id string) *products_domain.Product {
	price, err := price.NewPrice(42, "USD")
	assert.NoError(t, err)

	p, err := products_domain.NewProduct(products_domain.ID(id), "foo " + id, "bar " + id, price)
	assert.NoError(t, err)

	err = repo.Save(p)
	assert.NoError(t, err)

	return p
}
