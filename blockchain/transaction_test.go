package blockchain

import (
	"testing"
)

func TestValidTransaction(t *testing.T) {
	valid := Transaction{"90000000", Asset{"Initialize"}, 108167940}
	assertEqual(t, validate(valid), true, "")

	valid1 := Transaction{"100000000", Asset{"Initialize"}, 108167940}
	assertEqual(t, validate(valid1), true, "")

	invalid := Transaction{"89999999", Asset{"Initialize"}, 108167940}
	assertEqual(t, validate(invalid), false, "")
}
