package shop

import (
	shop_app "github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/shop/application"
)

func LoadShopFixtures(productsService shop_app.ProductsService) error {
	err := productsService.AddProduct(shop_app.AddProductCommand{
		ID:            "1",
		Name:          "Product 1",
		Description:   "Some extra description",
		PriceCents:    422,
		PriceCurrency: "USD",
	})
	if err != nil {
		return err
	}

	return productsService.AddProduct(shop_app.AddProductCommand{
		ID:            "2",
		Name:          "Product 2",
		Description:   "Another extra description",
		PriceCents:    333,
		PriceCurrency: "EUR",
	})
}
