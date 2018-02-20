package orders

import "github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/orders/domain/orders"

type MemoryRepository struct {
	orders []orders.Order
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{[]orders.Order{}}
}

func (m *MemoryRepository) Save(orderToSave *orders.Order) error {
	for i, p := range m.orders {
		if p.ID() == orderToSave.ID() {
			m.orders[i] = *orderToSave
			return nil
		}
	}

	m.orders = append(m.orders, *orderToSave)
	return nil
}

func (m MemoryRepository) ByID(id orders.ID) (*orders.Order, error) {
	for _, p := range m.orders {
		if p.ID() == id {
			return &p, nil
		}
	}

	return nil, orders.ErrNotFound
}
