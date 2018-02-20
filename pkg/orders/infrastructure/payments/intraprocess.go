package payments

import (
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/common/price"
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/orders/domain/orders"
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/payments/interfaces/intraprocess"
)

type IntraprocessService struct {
	orders chan <- intraprocess.OrderToProcess
}

func NewIntraprocessService(ordersChannel chan <- intraprocess.OrderToProcess) IntraprocessService {
	return IntraprocessService{ordersChannel}
}

func (i IntraprocessService) InitializeOrderPayment(id orders.ID, price price.Price) error {
	i.orders <- intraprocess.OrderToProcess{string(id), price}
	return nil
}
