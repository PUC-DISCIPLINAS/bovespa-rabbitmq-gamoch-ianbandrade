package structs

import "fmt"

type Order struct {
  Quantity int     `json:"quantity"`
  Value    float64 `json:"value"`
  Broker   string  `json:"broker"`
}

func (order *Order) String() string {
  quantity := order.Quantity
  value := order.Value
  broker := order.Broker

  return fmt.Sprintf("<quantidade: %d, valor: %.2f, corretora: %s>", quantity, value, broker)
}
