package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	ProductId      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product

func init() {
	productsJSON := `[
		{
			"productId": 9,
			"manufacturer": "Ut Nulla Corporation",
			"sku": "RXI73TND7UB",
			"upc": "36919",
			"pricePerUnit": "$5.13",
			"quantityOnHand": 8,
			"productName": "Trevor Solomon"
		},
		{
			"productId": 6,
			"manufacturer": "Est Arcu LLC",
			"sku": "BJF72STR2NP",
			"upc": "28594",
			"pricePerUnit": "$33.11",
			"quantityOnHand": 47,
			"productName": "Nicholas Santana"
		},
		{
			"productId": 8,
			"manufacturer": "Donec Corporation",
			"sku": "LVT31LHK8VJ",
			"upc": "27582",
			"pricePerUnit": "$41.27",
			"quantityOnHand": 59,
			"productName": "Rebecca Chambers"
		},
		{
			"productId": 2,
			"manufacturer": "Malesuada Limited",
			"sku": "FFO28CDN6GW",
			"upc": "43710",
			"pricePerUnit": "$97.50",
			"quantityOnHand": 97,
			"productName": "Miriam Holcomb"
		},
		{
			"productId": 10,
			"manufacturer": "Tincidunt Tempus Risus Ltd",
			"sku": "SJM20YBF7RH",
			"upc": "38977",
			"pricePerUnit": "$48.21",
			"quantityOnHand": 28,
			"productName": "Iris Hebert"
		}
	]`
	err := json.Unmarshal([]byte(productsJSON), &productList)
	if err != nil {
		log.Fatalln(err)
	}
}

func getNextId() int {
	highestId := -1
	for _, product := range productList {
		if highestId < product.ProductId {
			highestId = product.ProductId
		}
	}
	return highestId + 1
}

type fooHandler struct {
	Message string
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(productsJson)
		if err != nil {
			return
		}
	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductId != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newProduct.ProductId = getNextId()
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
		return
	}

}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(f.Message))
	if err != nil {
		return
	}
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("bar called"))
	if err != nil {
		return
	}
}

func main() {
	http.Handle("/foo", &fooHandler{Message: "foo called"})
	http.HandleFunc("/bar", barHandler)
	http.HandleFunc("/products", productsHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		return
	}
}
