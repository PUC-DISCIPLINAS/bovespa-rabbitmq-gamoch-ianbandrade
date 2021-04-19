package app

type PageDataControl struct {
  IsConnected bool `json:"isConnected"`
}

type PageDataConfiguration struct {
  RabbitMQAddress string `json:"rabbitMQAddress"`
}

type PageData struct {
  Action        string                `json:"action"`
  Error         string                `json:"error"`
  Control       PageDataControl       `json:"control"`
  Configuration PageDataConfiguration `json:"configuration"`
  Events        []string              `json:"events"`
}

type AppData struct {
  Data PageData `json:"data"`
}
