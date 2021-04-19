package transactions

import "sync"

var once sync.Once

var transactions *Transactions

func GetTransactions() *Transactions {
  once.Do(func() {
    transactions = GetDefaultValue()
  })

  return transactions
}

func GetDefaultValue() *Transactions {
  return &Transactions{
    Values: make([]Transaction, 0),
  }
}

func Reset() {
  transactions = GetDefaultValue()
}
