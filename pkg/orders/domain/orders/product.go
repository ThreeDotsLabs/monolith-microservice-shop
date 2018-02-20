package orders

import (
	"errors"

	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/common/price"
)

type ProductID string

var ErrEmptyProductID = errors.New("empty product ID")

type Product struct {
	id    ProductID
	name  string
	price price.Price
}

func NewProduct(id ProductID, name string, price price.Price) (Product, error) {
	if len(id) == 0 {
		return Product{}, ErrEmptyProductID
	}

	return Product{id, name, price}, nil
}

func (p Product) ID() ProductID {
	return p.id
}

func (p Product) Name() string {
	return p.name
}

func (p Product) Price() price.Price {
	return p.price
}
