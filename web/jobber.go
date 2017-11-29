package web

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Consumer task handler
type Consumer func(id string, body []byte) error

// NewJobber new jobbber
func NewJobber(url, queue string) (*Jobber, error) {
	it := &Jobber{url: url, queue: queue, consumers: make(map[string]Consumer)}
	if err := it.open(func(ch *amqp.Channel) error {
		_, err := ch.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return it, nil
}

// Jobber jobber
type Jobber struct {
	url       string
	queue     string
	consumers map[string]Consumer
}

// Status status
func (p *Jobber) Status() map[string]interface{} {
	rst := make(map[string]interface{})
	for k, v := range p.consumers {
		rst[k] = FuncName(v)
	}
	return rst
}

// Register register handler
func (p *Jobber) Register(_type string, hnd Consumer) {
	if _, ok := p.consumers[_type]; ok {
		log.Warn("handler for ", _type, " already exists, will override it")
	}
	p.consumers[_type] = hnd
}

// Receive receive
func (p *Jobber) Receive(consumer string) error {
	log.Info("waiting for messages, to exit press CTRL+C")
	return p.open(func(ch *amqp.Channel) error {
		if err := ch.Qos(1, 0, false); err != nil {
			return err
		}
		msgs, err := ch.Consume(p.queue, consumer, false, false, false, false, nil)
		if err != nil {
			return err
		}
		for d := range msgs {
			d.Ack(false)
			log.Info("receive message ", d.MessageId, "@", d.Type)
			now := time.Now()
			hnd, ok := p.consumers[d.Type]
			if !ok {
				return fmt.Errorf("unknown message type %s", d.Type)
			}
			if err := hnd(d.MessageId, d.Body); err != nil {
				return err
			}
			log.Infof("done %s %s", d.MessageId, time.Now().Sub(now))
		}
		return nil
	})
}

// Send send job
func (p *Jobber) Send(_type string, priority uint8, body interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(body)
	if err != nil {
		return err
	}
	return p.open(func(ch *amqp.Channel) error {
		return ch.Publish("", p.queue, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Priority:     priority,
			Body:         buf.Bytes(),
			Timestamp:    time.Now(),
			Type:         _type,
		})
	})
}
func (p *Jobber) open(f func(*amqp.Channel) error) error {
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return f(ch)
}
