package services

import (
  "encoding/json"
  "fmt"
  "github.com/streadway/amqp"
  "log"
  "stock-exchange/constants"
  "stock-exchange/data/app"
  "stock-exchange/data/offerBook"
  "stock-exchange/data/transactions"
  "stock-exchange/middlewares/rabbitMQMiddleware"
  "stock-exchange/middlewares/webSocketMiddleware"
  "stock-exchange/shared/structs"
  "strings"
  "time"
)

func AppHandler(pageData app.PageData) app.PageData {
  switch action := pageData.Action; action {
  case "RESET":
    pageData.Action = ""
    pageData = app.GetDefaultData()
    rabbitMQ := rabbitMQMiddleware.GetRabbitMQ()
    rabbitMQ.Connection.Close()
    transactions.Reset()
    offerBook.Reset()

  case "CONNECT":
    pageData.Action = ""
    appConstants := constants.GetAppConstants()
    rabbitMQ := rabbitMQMiddleware.GetRabbitMQ()
    amqpURI := pageData.Configuration.RabbitMQAddress

    if err := rabbitMQ.SetAmqpURI(amqpURI); err != nil {
      pageData.Error = err.Error()
      break
    }

    queue := appConstants.Queue
    exchangeName := appConstants.ConsumerExchange.Name
    exchangeType := appConstants.ConsumerExchange.Type

    if err := rabbitMQ.Consume(queue, exchangeName, exchangeType, "*.*", orderHandler); err != nil {
      pageData.Error = err.Error()
      break
    }

    pageData.Control.IsConnected = true
  }

  return pageData
}

func orderHandler(message amqp.Delivery) {
  appConstants := constants.GetAppConstants()
  rabbitMQ := rabbitMQMiddleware.GetRabbitMQ()

  var order structs.Order
  if err := json.Unmarshal(message.Body, &order); err != nil {
    log.Println("falha ao desserializar os dados da fila")
    return
  }

  formattedMessage := fmt.Sprintf("%s\t%s", message.RoutingKey, order.String())
  sendData(formattedMessage)

  exchangeName := appConstants.ProducerExchange.Name
  exchangeType := appConstants.ProducerExchange.Type

  routingKey := message.RoutingKey
  contentType := "text/plain"
  body := []byte(formattedMessage)

  if err := rabbitMQ.Publish(exchangeName, exchangeType, routingKey, contentType, body); err != nil {
    log.Println(err.Error())
    return
  }

  orderCopy := order
  routingKeyData := strings.Split(routingKey, ".")
  offerBookHandler(routingKeyData[0], routingKeyData[1], &orderCopy)
}

func offerBookHandler(operation, activeStock string, order *structs.Order) {
  book := offerBook.GetOfferBook()

  if operation == "compra" {
    if _, ok := book.Buys[activeStock]; !ok {
      book.Buys[activeStock] = make(map[*structs.Order]bool)
    }

    book.Buys[activeStock][order] = true

    for sellOrder := range book.Sells[activeStock] {
      if order.Value >= sellOrder.Value {
        var transactionQuantity int

        if order.Quantity > sellOrder.Quantity {
          transactionQuantity = sellOrder.Quantity
          order.Quantity -= sellOrder.Quantity
          delete(book.Sells[activeStock], sellOrder)
        } else if sellOrder.Quantity > order.Quantity {
          transactionQuantity = order.Quantity
          sellOrder.Quantity -= order.Quantity
          delete(book.Buys[activeStock], order)
        } else {
          transactionQuantity = order.Quantity
          delete(book.Sells[activeStock], sellOrder)
          delete(book.Buys[activeStock], order)
        }

        transactionHandler(
          activeStock,
          sellOrder.Broker,
          order.Broker,
          transactionQuantity,
          sellOrder.Value,
        )
      }
    }
  } else if operation == "venda" {
    if _, ok := book.Sells[activeStock]; !ok {
      book.Sells[activeStock] = make(map[*structs.Order]bool)
    }

    book.Sells[activeStock][order] = true

    for buyOrder := range book.Buys[activeStock] {
      if order.Value <= buyOrder.Value {
        var transactionQuantity int

        if buyOrder.Quantity > order.Quantity {
          transactionQuantity = order.Quantity
          buyOrder.Quantity -= order.Quantity
          delete(book.Sells[activeStock], order)
        } else if order.Quantity > buyOrder.Quantity {
          transactionQuantity = buyOrder.Quantity
          order.Quantity -= buyOrder.Quantity
          delete(book.Buys[activeStock], buyOrder)
        } else {
          transactionQuantity = buyOrder.Quantity
          delete(book.Sells[activeStock], order)
          delete(book.Buys[activeStock], order)
        }

        transactionHandler(
          activeStock,
          order.Broker,
          buyOrder.Broker,
          transactionQuantity,
          order.Value,
        )
      }
    }
  }
}

func transactionHandler(activeStock, sellBroker, buyBroker string, quantity int, value float64) {
  transactionsStorage := transactions.GetTransactions()

  location, _ := time.LoadLocation("America/Sao_Paulo")
  timeNow := time.Now().In(location)

  transaction := transactions.Transaction{
    DateTime:   timeNow.Format("02/01/2006 15:04"),
    SellBroker: sellBroker,
    BuyBroker:  buyBroker,
    Quantity:   quantity,
    Value:      value,
  }

  transactionsStorage.Values = append(transactionsStorage.Values, transaction)

  formattedMessage := fmt.Sprintf("transacao.%s\t%s", activeStock, transaction.String())
  sendData(formattedMessage)

  appConstants := constants.GetAppConstants()
  rabbitMQ := rabbitMQMiddleware.GetRabbitMQ()

  exchangeName := appConstants.ProducerExchange.Name
  exchangeType := appConstants.ProducerExchange.Type

  routingKey := fmt.Sprintf("transacao.%s", activeStock)
  contentType := "text/plain"
  body := []byte(formattedMessage)

  if err := rabbitMQ.Publish(exchangeName, exchangeType, routingKey, contentType, body); err != nil {
    log.Println(err.Error())
    return
  }
}

func sendData(message string) {
  appData := app.GetAppData()
  events := &appData.Data.Events
  *events = append(*events, message)

  pageDataJSON, err := json.Marshal(appData.Data)
  if err != nil {
    log.Println("falha ao serializar os dados da página")
    return
  }

  webSocket := webSocketMiddleware.GetWebSocket()
  webSocket.Broadcast(pageDataJSON)
}
