package rabbitMQMiddleware

import (
  "errors"
  "fmt"
  "github.com/streadway/amqp"
  "sync"
)

type RabbitMQ struct {
  Connection *amqp.Connection
  Channel    *amqp.Channel
  Queues     []string
}

var once sync.Once

var rabbitMQ *RabbitMQ

func GetRabbitMQ() *RabbitMQ {
  once.Do(func() {
    rabbitMQ = &RabbitMQ{
      Queues: make([]string, 0),
    }
  })

  return rabbitMQ
}

func (rabbitMQ *RabbitMQ) SetAmqpURI(amqpURI string) error {
  connection, err := amqp.Dial(amqpURI)
  if err != nil {
    return errors.New("falha ao criar uma conexão no RabbitMQ")
  }
  rabbitMQ.Connection = connection

  channel, err := connection.Channel()
  if err != nil {
    return errors.New("falha ao criar um canal no RabbitMQ")
  }
  rabbitMQ.Channel = channel

  return nil
}

func (rabbitMQ *RabbitMQ) Publish(exchangeName, exchangeType, routingKey, contentType string, body []byte) error {
  if err := declareExchange(rabbitMQ.Channel, exchangeName, exchangeType); err != nil {
    return err
  }

  if err := rabbitMQ.Channel.Publish(
    exchangeName,
    routingKey,
    false,
    false,
    amqp.Publishing{
      ContentType: contentType,
      Body:        body,
    },
  ); err != nil {
    return errors.New("falha ao publicar uma mensagem no RabbitMQ")
  }

  return nil
}

func (rabbitMQ *RabbitMQ) Consume(queueName, exchangeName, exchangeType, routingKey string, handler func(amqp.Delivery)) error {
  if err := declareExchange(rabbitMQ.Channel, exchangeName, exchangeType); err != nil {
    return err
  }

  queue, err := rabbitMQ.Channel.QueueDeclare(
    queueName,
    false,
    true,
    false,
    false,
    nil,
  )
  if err != nil {
    return fmt.Errorf("falha ao declarar a fila \"%s\" no RabbitMQ", queueName)
  }

  rabbitMQ.Queues = append(rabbitMQ.Queues, queue.Name)

  if err = rabbitMQ.Channel.QueueBind(
    queue.Name,
    routingKey,
    exchangeName,
    false,
    nil,
  ); err != nil {
    return fmt.Errorf("falha ao conectar a fila \"%s\" no RabbitMQ", queueName)
  }

  messages, err := rabbitMQ.Channel.Consume(
    queue.Name,
    "",
    true,
    false,
    false,
    false,
    nil,
  )
  if err != nil {
    return fmt.Errorf("falha ao criar um consumidor para a fila \"%s\" no RabbitMQ", queueName)
  }

  go func() {
    for message := range messages {
      handler(message)
    }
  }()

  return nil
}

func (rabbitMQ *RabbitMQ) DeleteQueue(queueName string) error {
  for i, queue := range rabbitMQ.Queues {
    if queue == queueName {

      _, err := rabbitMQ.Channel.QueueDelete(
        queueName,
        false,
        false,
        false,
      )
      if err != nil {
        return fmt.Errorf("falha ao deletar a fila \"%s\" no RabbitMQ", queueName)
      }

      rabbitMQ.Queues[i] = rabbitMQ.Queues[len(rabbitMQ.Queues)-1]
      rabbitMQ.Queues[len(rabbitMQ.Queues)-1] = ""
      rabbitMQ.Queues = rabbitMQ.Queues[:len(rabbitMQ.Queues)-1]
    }
  }

  return nil
}

func declareExchange(channel *amqp.Channel, exchangeName, exchangeType string) error {
  if err := channel.ExchangeDeclare(
    exchangeName,
    exchangeType,
    true,
    false,
    false,
    false,
    nil,
  ); err != nil {
    return fmt.Errorf("falha ao declarar a exchangeName \"%s\" no RabbitMQ", exchangeName)
  }

  return nil
}
