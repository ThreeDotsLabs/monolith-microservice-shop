package shop

import (
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/orders/domain/orders"
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/shop/interfaces/private/intraprocess"
)

type IntraprocessService struct {
	intraprocessInterface intraprocess.ProductInterface
}

func NewIntraprocessService(intraprocessInterface intraprocess.ProductInterface) IntraprocessService {
	return IntraprocessService{intraprocessInterface}
}

func (i IntraprocessService) ProductByID(id orders.ProductID) (orders.Product, error) {
	shopProduct, err := i.intraprocessInterface.ProductByID(string(id))
	if err != nil {
		return orders.Product{}, err
	}

	return OrderProductFromIntraprocess(shopProduct)
}

func OrderProductFromIntraprocess(shopProduct intraprocess.Product) (orders.Product, error) {
	return orders.NewProduct(orders.ProductID(shopProduct.ID), shopProduct.Name, shopProduct.Price)
}
