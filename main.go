package main

import (
	"github.com/babadee08/inventoryservice/database"
	"github.com/babadee08/inventoryservice/product"
	"github.com/babadee08/inventoryservice/receipt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		return
	}
}
