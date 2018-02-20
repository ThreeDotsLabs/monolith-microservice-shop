package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/orders/domain/orders"
	orders_http_interface "github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/orders/interfaces/public/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCases = []struct {
	Name                 string
	OrdersServiceAddress string
	ShopServiceAddress   string
}{
	{
		Name:                 "monolith",
		OrdersServiceAddress: "http://monolith:8080", // running from container, so :8080, not :8090
		ShopServiceAddress:   "http://monolith:8080",
	},
	{
		Name:                 "microservices",
		OrdersServiceAddress: "http://orders:8080", // running from container, so :8080, not :8070
		ShopServiceAddress:   "http://shop:8080",
	},
}

func TestOrderPath(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			orderID := placeOrder(t, tc.OrdersServiceAddress, "1")

			timeout := time.Now().Add(time.Second)
			for {
				if isOrderPaid(t, tc.OrdersServiceAddress, orderID) {
					break
				}

				if time.Now().After(timeout) {
					t.Fatal("timeouted: order is not paid")
					break
				}

				time.Sleep(time.Millisecond * 100)
			}
		})
	}

}

func isOrderPaid(t *testing.T, ordersServiceAddress string, orderID string) bool {
	resp := makeRequest(t, "GET", fmt.Sprintf("%s/orders/%s/paid", ordersServiceAddress, orderID), nil)
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	paidResponse := orders_http_interface.OrderPaidView{}
	err = json.Unmarshal(respBody, &paidResponse)
	assert.NoError(t, err)

	return paidResponse.IsPaid
}

func placeOrder(t *testing.T, ordersServiceAdddress string, productID string) string {
	resp := makeRequest(t, "POST", ordersServiceAdddress+"/orders", orders_http_interface.PostOrderRequest{
		ProductID: orders.ProductID(productID),
		Address: orders_http_interface.PostOrderAddress{
			Name:     "test name",
			Street:   "test street",
			City:     "test city",
			PostCode: "test post code",
			Country:  "test country",
		},
	})

	responseData := orders_http_interface.PostOrdersResponse{}

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(b, &responseData))

	require.EqualValues(t, http.StatusOK, resp.StatusCode)

	return responseData.OrderID
}

func TestProducts(t *testing.T) {
	expectedProducts := `
[
   {
      "id":"1",
      "name":"Product 1",
      "description":"Some extra description",
      "price":{
         "cents":422,
         "currency":"USD"
      }
   },
   {
      "id":"2",
      "name":"Product 2",
      "description":"Another extra description",
      "price":{
         "cents":333,
         "currency":"EUR"
      }
   }
]`

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			url := fmt.Sprintf("%s/products", tc.ShopServiceAddress)
			resp := makeRequest(t, "GET", url, nil)
			defer resp.Body.Close()

			greeting, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.JSONEq(t, expectedProducts, string(greeting))
		})
	}
}

func makeRequest(t *testing.T, method string, path string, data interface{}) (*http.Response) {
	var body []byte

	if data != nil {
		var err error
		body, err = json.Marshal(data)
		assert.NoError(t, err)
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	// resp.Body is io.Reader and is treated as stream, so you can read from it once.
	// We must set body again to allow read it again.
	bodyCopy, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyCopy))

	fmt.Printf("Request %s %s done, response status: %d: body: %s\n", method, path, resp.StatusCode, bodyCopy)

	return resp
}
