package offerBook

import "stock-exchange/shared/structs"

type OfferBook struct {
  Buys  map[string]map[*structs.Order]bool `json:"buys"`
  Sells map[string]map[*structs.Order]bool `json:"sells"`
}
