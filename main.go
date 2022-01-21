package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func findProductByID(productID int) (*Product, int) {
	for i, product := range productList {
		if product.ProductId == productID {
			return &product, i
		}
	}
	return nil, 0
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "product/")
	productID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product, listItemIndex := findProductByID(productID)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		productJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(productJSON)
		if err != nil {
			return
		}
	case http.MethodPut:
		// update product in the list
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedProduct.ProductId != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product = &updatedProduct
		productList[listItemIndex] = *product
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
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

func main() {
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/product/", productHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		return
	}
}
