package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/babadee08/inventoryservice/database"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("loading products...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v products loaded... \n", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"

	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%v] does not exist", fileName)
	}
	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductId] = productList[i]
	}
	return prodMap, nil
}

func getProduct(productId int) (*Product, error) {
	row := database.DbConn.QueryRow(`SELECT productId, 
       manufacturer, 
       sku, 
       upc, 
       pricePerUnit, 
       quantityOnHand, 
       productName FROM products WHERE productId = ?`, productId)
	product := &Product{}
	err := row.Scan(&product.ProductId,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return product, nil
	/*productMap.RLock()
	defer productMap.RUnlock()
	if product, ok := productMap.m[productId]; ok {
		return &product
	}
	return nil*/
}

func removeProduct(productId int) error {
	_, err := database.DbConn.Query(`DELETE FROM products WHERE productId = ?`, productId)
	if err != nil {
		return err
	}
	return nil
	/*productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, productId)*/
}

func getProductList() ([]Product, error) {
	results, err := database.DbConn.Query(`SELECT productId, 
       manufacturer, 
       sku, 
       upc, 
       pricePerUnit, 
       quantityOnHand, 
       productName FROM products`)
	if err != nil {
		return nil, err
	}
	defer func(results *sql.Rows) {
		err := results.Close()
		if err != nil {
			panic(err)
		}
	}(results)
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		err := results.Scan(&product.ProductId,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
	/*productMap.RLock()
	products := make([]Product, 0, len(productMap.m))
	for _, value := range productMap.m {
		products = append(products, value)
	}
	productMap.RUnlock()
	return products*/
}

func getProductIds() []int {
	productMap.RLock()
	productIds := []int{}
	for key := range productMap.m {
		productIds = append(productIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

func getNextProductId() int {
	productIds := getProductIds()
	return productIds[len(productIds)-1] + 1
}

func updateProduct(product Product) error {
	_, err := database.DbConn.Exec(`UPDATE products SET 
 	manufacturer=?,
 	sku=?,
 	upc=?,
 	pricePerUnit=CAST(? AS DECIMAL(13,2)),
 	quantityOnHand=?,
 	productName=? WHERE productId=?`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductId,
	)
	if err != nil {
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	result, err := database.DbConn.Exec(`INSERT INTO products 
 	(manufacturer,
 	sku,
 	upc,
 	pricePerUnit,
 	quantityOnHand,
 	productName) VALUES (?, ?, ?, ?, ?, ?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
	)
	if err != nil {
		return 0, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(insertId), nil
}

func addOrUpdateProduct(product Product) (int, error) {
	addOrUpdateID := -1
	if product.ProductId > 0 {
		oldProduct, _ := getProduct(product.ProductId)
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%v] doesn't exist", product.ProductId)
		}
		addOrUpdateID = product.ProductId
	} else {
		addOrUpdateID = getNextProductId()
		product.ProductId = addOrUpdateID
	}

	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}
