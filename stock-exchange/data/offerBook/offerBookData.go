package offerBook

import (
  "stock-exchange/shared/structs"
  "sync"
)

var once sync.Once

var offerBook *OfferBook

func GetOfferBook() *OfferBook {
  once.Do(func() {
    offerBook = GetDefaultValue()
  })

  return offerBook
}

func GetDefaultValue() *OfferBook {
  return &OfferBook{
    Sells: map[string]map[*structs.Order]bool{},
    Buys: map[string]map[*structs.Order]bool{},
  }
}

func Reset() {
  offerBook = GetDefaultValue()
}
