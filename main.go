package main

import (
	"github.com/babadee08/inventoryservice/product"
	"net/http"
)

const apiBasePath = "/api"

func main() {
	product.SetupRoutes(apiBasePath)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		return
	}
}
