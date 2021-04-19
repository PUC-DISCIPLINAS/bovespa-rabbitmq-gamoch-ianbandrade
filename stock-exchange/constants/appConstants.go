package constants

import (
  "fmt"
  "stock-exchange/utils"
  "sync"
)

type Exchange struct {
  Name string `json:"name"`
  Type string `json:"type"`
}

type AppConstants struct {
  RabbitMQAddress  string   `json:"rabbitMQAddress"`
  ProducerExchange Exchange `json:"producerExchange"`
  ConsumerExchange Exchange `json:"consumerExchange"`
  Queue            string   `json:"queue"`
}

var once sync.Once

var appConstants *AppConstants

func GetAppConstants() *AppConstants {
  once.Do(func() {
    user := utils.GetEnv("RABBITMQ_USER", "guest")
    pass := utils.GetEnv("RABBITMQ_PASS", "guest")
    host := utils.GetEnv("RABBITMQ_HOST", "localhost")
    port := utils.GetEnv("RABBITMQ_PORT", "5672")

    appConstants = &AppConstants{
      RabbitMQAddress: fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port),
      ProducerExchange: Exchange{
        Name: "BOLSADEVALORES",
        Type: "topic",
      },
      ConsumerExchange: Exchange{
        Name: "BROKER",
        Type: "topic",
      },
      Queue: "BROKER",
    }
  })

  return appConstants
}
