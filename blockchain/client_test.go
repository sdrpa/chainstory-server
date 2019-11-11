package blockchain

import (
	"fmt"
	"testing"
)

func TestGetTransactions(t *testing.T) {
	client := Client{"https://testnet.lisk.io", "4389113358205759328L"}
	xs, err := client.getTransactions(0)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(xs)
}

func TestGetValidTransactions(t *testing.T) {
	client := Client{"https://testnet.lisk.io", "4389113358205759328L"}
	xs := client.GetValidTransactions(0)
	fmt.Println(xs)
}
