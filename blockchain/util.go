package blockchain

func filter(vs []Transaction, f func(Transaction) bool) []Transaction {
	xs := make([]Transaction, 0)
	for _, v := range vs {
		if f(v) {
			xs = append(xs, v)
		}
	}
	return xs
}
