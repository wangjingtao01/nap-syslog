package dispatch

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Responser struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	uri          string
	exchangeName string
	exchangeType string
}

func NewResponser(d DispatcherInstance) (*Responser, error) {

	r := &Responser{
		conn:    nil,
		channel: nil,

		exchangeName: d.Exchange,
		exchangeType: d.ExchangeType,
		uri:          d.URI,
	}

	var err error
	log.Printf("Dialing %q\n", r.uri)
	r.conn, err = amqp.Dial(r.uri)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %s", err)
	}

	log.Printf("got Channel, declaring %q Exchange (%q)\n", r.exchangeType, r.exchangeName)
	if err := r.channel.ExchangeDeclare(
		r.exchangeName, // name
		r.exchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return nil, fmt.Errorf("exchange declare: %s", err)
	}

	// declare trigger queues
	for _, v := range d.Triggers {
		queue, err := r.channel.QueueDeclare(
			v.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("queue declare: %s", err)
		}

		log.Printf("Queue %s bound to Exchange %s\n", queue.Name, r.exchangeName)
		if err = r.channel.QueueBind(
			v.Queue,        // name of the queue
			v.RoutingKey,   // bindingKey
			r.exchangeName, // sourceExchange
			false,          // noWait
			nil,            // arguments
		); err != nil {
			return nil, fmt.Errorf("queue binding: %s", err)
		}
	}

	return r, nil
}

func (s *Responser) Send(response interface{}, trigger Trigger) {
	b, err := json.Marshal(response)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("sending message %s", string(b))

	if err = s.channel.Publish(
		s.exchangeName,     // publish to an exchange
		trigger.RoutingKey, // routing to 0 or more queues
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			Headers: map[string]interface{}{
				"__TypeId__": "net.skycloud.nap.messaging.model.LogEvent",
			},
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			Body:            b,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		log.Printf("exchange publish: %s", err)
	}
}
