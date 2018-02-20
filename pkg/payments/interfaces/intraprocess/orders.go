package intraprocess

import (
	"log"
	"sync"

	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/common/price"
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/payments/application"
)

type OrderToProcess struct {
	ID    string
	Price price.Price
}

type PaymentsInterface struct {
	orders            <-chan OrderToProcess
	service           application.PaymentsService
	orderProcessingWg *sync.WaitGroup
	runEnded          chan struct{}
}

func NewPaymentsInterface(orders <-chan OrderToProcess, service application.PaymentsService) PaymentsInterface {
	return PaymentsInterface{
		orders,
		service,
		&sync.WaitGroup{},
		make(chan struct{}, 1),
	}
}

func (o PaymentsInterface) Run() {
	defer func() {
		o.runEnded <- struct{}{}
	}()

	for order := range o.orders {
		go func(orderToPay OrderToProcess) {
			o.orderProcessingWg.Add(1)
			defer o.orderProcessingWg.Done()

			if err := o.service.InitializeOrderPayment(orderToPay.ID, orderToPay.Price); err != nil {
				log.Print("Cannot initialize payment:", err)
			}
		}(order)
	}
}

func (o PaymentsInterface) Close() {
	o.orderProcessingWg.Wait()
	<-o.runEnded
}
