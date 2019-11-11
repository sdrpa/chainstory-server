package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Client - Lisk blockchain client
type Client struct {
	NodeURL          string
	RecipientAddress string
}

// GetValidTransactions - Returns valid transactions
func (c Client) GetValidTransactions(offset int) []Transaction {
	xs, err := c.getTransactions(offset)
	if err != nil {
		return []Transaction{}
	}
	ys := filter(xs, validate)
	return ys
}

// https://lisk.io/documentation/lisk-core/api#/Transactions/getTransactions
// curl -i -X GET https://testnet.lisk.io/api/transactions?recipientId=4389113358205759328L&offset=0
func (c Client) getTransactions(offset int) ([]Transaction, error) {
	url := c.NodeURL + "/api/transactions?recipientId=" + c.RecipientAddress + "&offset=" + strconv.Itoa(offset) + "&limit=100"
	res, err := http.Get(url)
	if err != nil {
		return []Transaction{}, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Transaction{}, err
	}
	type JSON struct {
		Data []Transaction `json:"data"`
	}
	jsonRes := new(JSON)
	unmarshalErr := json.Unmarshal(body, &jsonRes)
	if unmarshalErr != nil {
		return []Transaction{}, err
	}
	return jsonRes.Data, nil
}
