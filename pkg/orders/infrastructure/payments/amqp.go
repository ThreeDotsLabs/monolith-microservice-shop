package payments

import (
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/common/price"
	"github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/orders/domain/orders"
	payments_amqp_interface "github.com/ThreeDotsLabs/monolith-microservice-shop/pkg/payments/interfaces/amqp"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type AMQPService struct {
	queue   amqp.Queue
	channel *amqp.Channel
}

func NewAMQPService(url, queueName string) (AMQPService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return AMQPService{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return AMQPService{}, err
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
		return AMQPService{}, err
	}

	return AMQPService{q, ch}, nil
}

func (i AMQPService) InitializeOrderPayment(id orders.ID, price price.Price) error {
	order := payments_amqp_interface.OrderToProcessView{
		ID: string(id),
		Price: payments_amqp_interface.PriceView{
			Cents:    price.Cents(),
			Currency: price.Currency(),
		},
	}

	b, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "cannot marshal order for amqp")
	}

	err = i.channel.Publish(
		"",
		i.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	if err != nil {
		return errors.Wrap(err, "cannot send order to amqp")
	}

	log.Printf("sent order %s to amqp", id)

	return nil
}

func (i AMQPService) Close() error {
	return i.channel.Close()
}
