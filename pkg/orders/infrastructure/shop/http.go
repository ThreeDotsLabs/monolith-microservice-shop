package shop

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/common/price"
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/orders/domain/orders"
	http_interface "github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/shop/interfaces/private/http"
	"github.com/pkg/errors"
)

type HTTPClient struct {
	address string
}

func NewHTTPClient(address string) HTTPClient {
	return HTTPClient{address}
}

func (h HTTPClient) ProductByID(id orders.ProductID) (orders.Product, error) {
	resp, err := http.Get(fmt.Sprintf("%s/products/%s", h.address, id))
	if err != nil {
		return orders.Product{}, errors.Wrap(err, "request to shop failed")
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return orders.Product{}, errors.Wrap(err, "cannot read response")
	}

	productView := http_interface.ProductView{}
	if err := json.Unmarshal(b, &productView); err != nil {
		return orders.Product{}, errors.Wrapf(err, "cannot decode response: %s", b)
	}

	return OrderProductFromHTTP(productView)
}


func OrderProductFromHTTP(shopProduct http_interface.ProductView) (orders.Product, error) {
	productPrice, err := OrderProductPriceFromHTTP(shopProduct.Price)
	if err != nil {
		return orders.Product{}, errors.Wrap(err, "cannot decode price")
	}

	return orders.NewProduct(orders.ProductID(shopProduct.ID), shopProduct.Name, productPrice)
}

func OrderProductPriceFromHTTP(priceView http_interface.PriceView) (price.Price, error) {
	return price.NewPrice(priceView.Cents, priceView.Currency)
}