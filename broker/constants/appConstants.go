package constants

import (
  "broker/utils"
  "fmt"
  "strconv"
  "sync"
)

type Exchange struct {
  Name string `json:"name"`
  Type string `json:"type"`
}

type AppConstants struct {
  BrokerName       string   `json:"brokerName"`
  RabbitMQAddress  string   `json:"rabbitMQAddress"`
  ProducerExchange Exchange `json:"producerExchange"`
  ConsumerExchange Exchange `json:"consumerExchange"`
}

var once sync.Once

var appConstants *AppConstants

func GetAppConstants() *AppConstants {
  once.Do(func() {
    user := utils.GetEnv("RABBITMQ_USER", "guest")
    pass := utils.GetEnv("RABBITMQ_PASS", "guest")
    host := utils.GetEnv("RABBITMQ_HOST", "localhost")
    port := utils.GetEnv("RABBITMQ_PORT", "5672")

    randomInt := strconv.Itoa(utils.RandomInt(1, 999))

    appConstants = &AppConstants{
      BrokerName:      fmt.Sprintf("BRK%s", utils.GetEnv("REPLICA_ID", randomInt)),
      RabbitMQAddress: fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port),
      ProducerExchange: Exchange{
        Name: "BROKER",
        Type: "topic",
      },
      ConsumerExchange: Exchange{
        Name: "BOLSADEVALORES",
        Type: "topic",
      },
    }
  })

  return appConstants
}
