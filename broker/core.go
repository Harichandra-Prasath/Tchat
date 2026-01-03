package broker

import (
	"fmt"

	"github.com/Harichandra-Prasath/Tchat/configs"
	"github.com/Harichandra-Prasath/Tchat/logging"
	ampq "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Type string
	DSN  string
	Data []byte
}

var rmqConn *ampq.Connection

func IntialiseBroker() error {
	conn, err := ampq.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", configs.GnCfg.RMQUser, configs.GnCfg.RMQPassword, configs.GnCfg.RMQHost, configs.GnCfg.RMQPort))
	if err != nil {
		return fmt.Errorf("creating broker conn: %s", err.Error())
	}
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("creating conn channel: %s", err.Error())
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(configs.GnCfg.RMQUserExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declaring the exchange: %s", err.Error())
	}

	rmqConn = conn

	return nil
}

func PublishEvents(event *Event) error {
	ch, err := rmqConn.Channel()
	if err != nil {
		return fmt.Errorf("creating conn channel: %s", err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("users."+event.DSN, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declaring queue: %s", err.Error())
	}

	err = ch.QueueBind(q.Name, event.DSN, configs.GnCfg.RMQUserExchange, false, nil)
	if err != nil {
		return fmt.Errorf("binding queue: %s", err.Error())
	}

	err = ch.Publish(configs.GnCfg.RMQUserExchange, event.DSN, false, false, ampq.Publishing{ContentType: "application/json", Body: event.Data, DeliveryMode: 2, Headers: ampq.Table{"type": event.Type}})
	if err != nil {
		return fmt.Errorf("publishing message: %s", err.Error())
	}

	return nil
}

func ConsumeEvents(evch chan *Event, DSN string) error {
	ch, err := rmqConn.Channel()
	if err != nil {
		return fmt.Errorf("creating conn channel: %s", err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("users."+DSN, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declaring queue: %s", err.Error())
	}

	for {
		msg, ok, err := ch.Get(q.Name, false)
		if err != nil {
			logging.Logger.Error("Message Consuming", "err", err.Error())
			continue
		}
		if !ok {
			break
		}
		evch <- &Event{Type: msg.Headers["type"].(string), Data: msg.Body, DSN: DSN}
		msg.Ack(false)
	}

	return nil
}
