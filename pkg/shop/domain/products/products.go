package products

import (
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/common/price"
	"errors"
)

type ID string

var (
	ErrEmptyID   = errors.New("empty product ID")
	ErrEmptyName = errors.New("empty product name")
)

type Product struct {
	id ID

	name        string
	description string

	price price.Price
}

func NewProduct(id ID, name string, description string, price price.Price) (*Product, error) {
	if len(id) == 0 {
		return nil, ErrEmptyID
	}
	if len(name) == 0 {
		return nil, ErrEmptyName
	}

	return &Product{id, name, description, price}, nil
}

func (p Product) ID() ID {
	return p.id
}

func (p Product) Name() string {
	return p.name
}

func (p Product) Description() string {
	return p.description
}

func (p Product) Price() price.Price {
	return p.price
}
