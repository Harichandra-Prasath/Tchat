package broker

import (
	"fmt"

	"github.com/Harichandra-Prasath/Tchat/configs"
	ampq "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Type string
	DSN  string
	Data []byte
}

var rmqChn *ampq.Channel

func IntialiseBroker() error {
	conn, err := ampq.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", configs.GnCfg.RMQUser, configs.GnCfg.RMQPassword, configs.GnCfg.RMQHost, configs.GnCfg.RMQPort))
	if err != nil {
		return fmt.Errorf("creating broker conn: %s", err.Error())
	}
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("creating conn channel: %s", err.Error())
	}

	err = ch.ExchangeDeclare(configs.GnCfg.RMQUserExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declaring the exchange: %s", err.Error())
	}

	rmqChn = ch

	return nil
}

func PublishEvents(event *Event) error {
	q, err := rmqChn.QueueDeclare("users."+event.DSN, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declaring queue: %s", err.Error())
	}

	err = rmqChn.Publish(configs.GnCfg.RMQUserExchange, q.Name, false, false, ampq.Publishing{ContentType: "application/json", Body: event.Data})
	if err != nil {
		return fmt.Errorf("publishing message: %s", err.Error())
	}

	return nil
}

func ConsumeEvents(ch chan *Event, DSN string) error {
	q, err := rmqChn.QueueDeclare("users."+DSN, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declaring queue: %s", err.Error())
	}

	msgs, err := rmqChn.Consume(q.Name, configs.GnCfg.RMQUserExchange, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("registering the consumer: %s", err.Error())
	}

	for d := range msgs {
		ch <- &Event{DSN: DSN, Data: d.Body}
	}

	return nil
}
