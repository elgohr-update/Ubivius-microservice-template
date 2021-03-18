package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Ubivius/microservice-template/pkg/data"
)

// GetProducts returns the full list of products
func (productHandler *ProductsHandler) GetProducts(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Println("Handle GET products")
	productList := productHandler.db.GetProducts()
	err := json.NewEncoder(responseWriter).Encode(productList)
	if err != nil {
		productHandler.logger.Println("[ERROR] serializing product", err)
		http.Error(responseWriter, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// GetProductByID returns a single product from the database
func (productHandler *ProductsHandler) GetProductByID(responseWriter http.ResponseWriter, request *http.Request) {
	id := getProductID(request)

	productHandler.logger.Println("[DEBUG] getting id", id)

	product, err := productHandler.db.GetProductByID(id)

	switch err {
	case nil:
		err = json.NewEncoder(responseWriter).Encode(product)
		if err != nil {
			productHandler.logger.Println("[ERROR] serializing product", err)
		}
	case data.ErrorProductNotFound:
		productHandler.logger.Println("[ERROR] fetching product", err)
		http.Error(responseWriter, "Product not found", http.StatusBadRequest)
		return
	default:
		productHandler.logger.Println("[ERROR] fetching product", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

}
