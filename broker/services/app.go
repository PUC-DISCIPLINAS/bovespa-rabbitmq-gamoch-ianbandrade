package services

import (
  "broker/constants"
  "broker/data/app"
  "broker/middlewares/rabbitMQMiddleware"
  "broker/middlewares/webSocketMiddleware"
  "broker/shared/structs"
  "encoding/json"
  "fmt"
  "github.com/streadway/amqp"
  "log"
)

func AppHandler(pageData app.PageData) app.PageData {
  switch action := pageData.Action; action {
  case "RESET":
    pageData.Action = ""
    pageData = app.GetDefaultData()
    rabbitMQ := rabbitMQMiddleware.GetRabbitMQ()
    rabbitMQ.Connection.Close()

  case "CONNECT":
    pageData.Action = ""
    rabbitMQ := rabbitMQMiddleware.GetRabbitMQ()
    amqpURI := pageData.Configuration.RabbitMQAddress

    if err := rabbitMQ.SetAmqpURI(amqpURI); err != nil {
      pageData.Error = err.Error()
      break
    }

    pageData.Control.IsConnected = true

  case "ORDER":
    pageData.Action = ""
    rabbitMQ := rabbitMQMiddleware.GetRabbitMQ()
    appConstants := constants.GetAppConstants()

    clientOrder := pageData.Order
    order := structs.Order{
      Quantity: clientOrder.Quantity,
      Value:    clientOrder.Value,
      Broker:   appConstants.BrokerName,
    }

    exchangeName := appConstants.ProducerExchange.Name
    exchangeType := appConstants.ProducerExchange.Type
    routingKey := fmt.Sprintf("%s.%s", clientOrder.Operation, clientOrder.ActiveStock)
    contentType := "application/json"

    orderMessage, err := json.Marshal(order)
    if err != nil {
      log.Println("falha ao serializar os dados do pedido")
      break
    }

    if err = rabbitMQ.Publish(exchangeName, exchangeType, routingKey, contentType, orderMessage); err != nil {
      log.Println(err.Error())
      break
    }

  default:
    topicsHandler(pageData.Topics.Actives)
  }

  return pageData
}

func topicsHandler(activeStocks map[string]bool) {
  rabbitMQ := rabbitMQMiddleware.GetRabbitMQ()
  appConstants := constants.GetAppConstants()

  exchangeName := appConstants.ConsumerExchange.Name
  exchangeType := appConstants.ConsumerExchange.Type

  for activeStock, isSubscribed := range activeStocks {
    queue := fmt.Sprintf("%s_%s", appConstants.BrokerName, activeStock)
    routingKey := fmt.Sprintf("*.%s", activeStock)

    if isSubscribed {
      if err := rabbitMQ.Consume(queue, exchangeName, exchangeType, routingKey, consumerHandler); err != nil {
        log.Println(err.Error())
        continue
      }
    } else {
      if err := rabbitMQ.DeleteQueue(queue); err != nil {
        log.Println(err.Error())
        continue
      }
    }
  }
}

func consumerHandler(message amqp.Delivery) {
  appData := app.GetAppData()
  events := &appData.Data.Topics.Events
  *events = append(*events, string(message.Body))

  pageDataJSON, err := json.Marshal(appData.Data)
  if err != nil {
    log.Println("falha ao serializar os dados da página")
    return
  }

  webSocket := webSocketMiddleware.GetWebSocket()
  webSocket.Broadcast(pageDataJSON)
}
