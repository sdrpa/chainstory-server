package blockchain

import "strconv"

// Asset - Lisk transaction asset
type Asset struct {
	Data string `json:"data"`
}

// Transaction - Lisk transaction
type Transaction struct {
	Amount    string `json:"amount"`
	Asset     Asset  `json:"asset"`
	Timestamp int    `json:"timestamp"`
}

// MakeTransaction - Creates a default transaction
func MakeTransaction() Transaction {
	return Transaction{"0", Asset{""}, 0}
}

func validate(t Transaction) bool {
	required := 90000000 // Beddows
	amount, err := strconv.Atoi(t.Amount)
	if err != nil {
		return false
	}
	if amount < required {
		return false
	}
	return true
}
