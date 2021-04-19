package app

type PageDataControl struct {
  IsConnected bool `json:"isConnected"`
}

type PageDataConfiguration struct {
  RabbitMQAddress string `json:"rabbitMQAddress"`
}

type PageDataOrder struct {
  Quantity    int     `json:"quantity,string"`
  Value       float64 `json:"value,string"`
  ActiveStock string  `json:"activeStock"`
  Operation   string  `json:"operation"`
}

type PageDataTopics struct {
  Actives map[string]bool `json:"actives"`
  Events  []string        `json:"events"`
}

type PageData struct {
  Action        string                `json:"action"`
  Error         string                `json:"error"`
  Control       PageDataControl       `json:"control"`
  Configuration PageDataConfiguration `json:"configuration"`
  Order         PageDataOrder         `json:"order"`
  Topics        PageDataTopics        `json:"topics"`
}

type AppData struct {
  Broker string   `json:"broker"`
  Data   PageData `json:"data"`
}
