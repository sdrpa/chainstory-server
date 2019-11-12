package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	artisan "./artisan"
	blockchain "./blockchain"
	cors "./cors"
	graphics "./graphics"
	ipfs "./ipfs"
)

// Test-net
// const (
// 	NodeURL          = "https://testnet.lisk.io"
// 	RecipientAddress = "4389113358205759328L"
// )

// Main-net
const (
	NodeURL          = "https://hub21.lisk.io"
	RecipientAddress = "16724144933475539753L"
)

func dumpRequest(r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(dump))
}

func home(w http.ResponseWriter, r *http.Request) {
	//dumpRequest(r)
	fmt.Fprintf(w, time.Now().Format(time.RFC850))
}

func add(w http.ResponseWriter, r *http.Request) {
	const maxFileSize = 1 * 1024
	if r.ContentLength > maxFileSize {
		http.Error(w, "Request too large", http.StatusExpectationFailed)
		return
	}
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		type JSON struct {
			Version int    `json:"version"`
			Index   int    `json:"index"`
			Message string `json:"message"`
			Data    string `json:"data"`
		}
		obj := new(JSON)
		objErr := json.Unmarshal(body, &obj)
		if objErr != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		//fmt.Println(obj.Version, obj.Index, obj.Data)
		//hash := "QmRAQB6YaCyidP37UdDnjFY5vQuiBrcqdyoW1CuDgwxkD4"
		hash := ipfs.Add(graphics.Drawing{obj.Version, obj.Index, obj.Message, obj.Data})

		type Response struct {
			Success   bool   `json:"success"`
			Recipient string `json:"recipient"`
			File      string `json:"file"`
		}
		res := Response{true, RecipientAddress, hash}
		resObj, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resObj)
	} else {
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed)
	}
}

func drawing(w http.ResponseWriter, r *http.Request) {
	indexes, base64String := artisan.DrawingBase64("./public/image.png")
	type Response struct {
		Success bool   `json:"success"`
		Indexes []int  `json:"indexes"` // array of colored indexes
		Base64  string `json:"base64"`
	}
	res := Response{true, indexes, base64String}
	resObj, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resObj)
}

func drawingAt(w http.ResponseWriter, r *http.Request) {
	last := strings.Replace(r.URL.Path, "/drawing/", "", 1)
	index, err := strconv.Atoi(last)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusInternalServerError)
		return
	}
	drawing, transaction, err := artisan.DrawingAt(index)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	type Response struct {
		Success   bool   `json:"success"`
		Index     int    `json:"index"`
		Message   string `json:"message"`
		Amount    int    `json:"amount"`
		Timestamp int    `json:"timestamp"`
	}
	amount, err := strconv.Atoi(transaction.Amount)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusInternalServerError)
		return
	}
	res := Response{true, drawing.Index, drawing.Message, amount, transaction.Timestamp}
	resObj, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resObj)
}

func node(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, NodeURL)
}

// Custom HTTP Routing in Go - https://gist.github.com/reagent/043da4661d2984e9ecb1ccb5343bf438
func main() {
	client := blockchain.Client{NodeURL, RecipientAddress}
	go artisan.Run(client)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/add", add)
	mux.HandleFunc("/drawing", drawing)
	mux.HandleFunc("/drawing/", drawingAt)
	mux.HandleFunc("/node", node)

	// cors.Default() setups the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	// See https://github.com/rs/cors for more options.
	// handler := cors.Default().Handler(mux)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	handler := c.Handler(mux)
	err := http.ListenAndServe(":3001", handler)
	if err != nil {
		log.Panic(err)
	}
}
