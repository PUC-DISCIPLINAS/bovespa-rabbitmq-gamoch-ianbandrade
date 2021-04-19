package app

import (
  "stock-exchange/constants"
  "sync"
)

var once sync.Once

var appData *AppData

func GetAppData() *AppData {
  once.Do(func() {
    appData = &AppData{
      Data: GetDefaultData(),
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
    Events: make([]string, 0),
  }
}
