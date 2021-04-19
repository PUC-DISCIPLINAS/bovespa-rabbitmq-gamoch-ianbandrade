package app

import (
  "broker/constants"
  "sync"
)

var once sync.Once

var appData *AppData

func GetAppData() *AppData {
  once.Do(func() {
    appConstants := constants.GetAppConstants()

    appData = &AppData{
      Broker: appConstants.BrokerName,
      Data:   GetDefaultData(),
    }
  })

  return appData
}

func GetDefaultData() PageData {
  appConstants := constants.GetAppConstants()

  return PageData{
    Configuration: PageDataConfiguration{
      RabbitMQAddress: appConstants.RabbitMQAddress,
    },
    Order: PageDataOrder{
      Quantity:    1,
      Value:       1,
      ActiveStock: "ABEV3",
      Operation:   "compra",
    },
    Topics: PageDataTopics{
      Events: make([]string, 0),
      Actives: map[string]bool{
        "ABEV3": false,
        "PETR4": false,
        "VALE5": false,
        "iTUB4": false,
        "BBDC4": false,
        "BBAS3": false,
        "CiEL3": false,
        "PETR3": false,
        "HYPE3": false,
        "VALE3": false,
        "BBSE3": false,
        "CTiP3": false,
        "GGBR4": false,
        "FiBR3": false,
        "RADL3": false,
      },
    },
  }
}
