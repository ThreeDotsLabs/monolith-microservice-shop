package products

import (
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/shop/domain/products"
)

type MemoryRepository struct {
	products []products.Product
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{[]products.Product{}}
}

func (m *MemoryRepository) Save(productToSave *products.Product) error {
	for i, p := range m.products {
		if p.ID() == productToSave.ID() {
			m.products[i] = *productToSave
			return nil
		}
	}

	m.products = append(m.products, *productToSave)
	return nil
}

func (m MemoryRepository) ByID(id products.ID) (*products.Product, error) {
	for _, p := range m.products {
		if p.ID() == id {
			return &p, nil
		}
	}

	return nil, products.ErrNotFound
}

func (m MemoryRepository) AllProducts() ([]products.Product, error) {
	return m.products, nil
}
