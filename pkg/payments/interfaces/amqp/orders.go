package amqp

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/common/price"
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/payments/application"
	"github.com/streadway/amqp"
)

type OrderToProcessView struct {
	ID    string `json:"id"`
	Price PriceView
}

type PriceView struct {
	Cents    uint   `json:"cents"`
	Currency string `json:"currency"`
}

type PaymentsInterface struct {
	conn    *amqp.Connection
	queue   amqp.Queue
	channel *amqp.Channel

	service application.PaymentsService
}

func NewPaymentsInterface(url string, queueName string, service application.PaymentsService) (PaymentsInterface, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return PaymentsInterface{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return PaymentsInterface{}, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return PaymentsInterface{}, err
	}

	return PaymentsInterface{conn, q, ch, service}, nil
}

func (o PaymentsInterface) Run(ctx context.Context) error {
	msgs, err := o.channel.Consume(
		o.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	done := ctx.Done()
	defer func() {
		if err := o.conn.Close(); err != nil {
			log.Print("cannot close conn: ", err)
		}
		if err := o.channel.Close(); err != nil {
			log.Print("cannot close channel: ", err)
		}
	}()

	for {
		select {
		case msg := <-msgs:
			err := o.processMsg(msg)
			if err != nil {
				log.Printf("cannot process msg: %s, err: %s", msg.Body, err)
			}
		case <-done:
			return nil
		}
	}
}

func (o PaymentsInterface) processMsg(msg amqp.Delivery) error {
	orderView := OrderToProcessView{}
	err := json.Unmarshal(msg.Body, &orderView)
	if err != nil {
		log.Printf("cannot decode msg: %s, error: %s", string(msg.Body), err)
	}

	orderPrice, err := price.NewPrice(orderView.Price.Cents, orderView.Price.Currency)
	if err != nil {
		log.Printf("cannot decode price for msg %s: %s", string(msg.Body), err)

	}

	return o.service.InitializeOrderPayment(orderView.ID, orderPrice)
}
