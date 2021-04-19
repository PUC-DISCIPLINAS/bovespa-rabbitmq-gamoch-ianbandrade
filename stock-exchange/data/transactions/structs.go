package transactions

import "fmt"

type Transaction struct {
  DateTime   string  `json:"dateTime"`
  SellBroker string  `json:"sellBroker"`
  BuyBroker  string  `json:"buyBroker"`
  Quantity   int     `json:"quantity"`
  Value      float64 `json:"value"`
}

type Transactions struct {
  Values []Transaction `json:"values"`
}

func (transaction *Transaction) String() string {
  dateTime := transaction.DateTime
  sellBroker := transaction.SellBroker
  buyBroker := transaction.BuyBroker
  quantity := transaction.Quantity
  value := transaction.Value

  return fmt.Sprintf(
    "<data-hora: %s, corretora-venda: %s, "+
      "corretora-compra: %s, quantidade: %d, valor: %.2f>",
    dateTime, sellBroker, buyBroker, quantity, value,
  )
}
