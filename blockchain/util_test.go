package blockchain

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func transactions() []Transaction {
	return []Transaction{
		Transaction{"90000000", Asset{"A"}, 108167942},
		Transaction{"80000000", Asset{"B"}, 108167941},
		Transaction{"70000000", Asset{"C"}, 108167940},
	}
}
func TestFilter(t *testing.T) {
	ts := transactions()
	xs := filter(ts, validate)
	fmt.Println(xs)
}
