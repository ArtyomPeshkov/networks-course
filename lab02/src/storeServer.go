package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type RequsetProductStructure struct {
	Name        string
	Description string
	Icon        string
}

type StoreServer struct {
	store *Store
}

func NewStoreServer() *StoreServer {
	store := new(Store)
	store.products = map[int]Product{}
	return &StoreServer{store: store}
}

func (ts *StoreServer) processRequest(w http.ResponseWriter, req *http.Request) (map[string]string, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body.", http.StatusBadRequest)
		return map[string]string{}, err
	}

	var readJson map[string]string

	err = json.Unmarshal(data, &readJson)
	if err != nil {
		http.Error(w, "Incorrect product format in request.", http.StatusBadRequest)
		return map[string]string{}, err
	}

	return readJson, nil
}

func checkProductFields(obtainedProduct map[string]string) (bool, bool, error) {
	_, hasName := obtainedProduct["name"]
	_, hasDescription := obtainedProduct["description"]
	hasUnpleasentField := len(obtainedProduct) > 3

	if !hasDescription && !hasName {
		return hasName, hasDescription, fmt.Errorf("product does not contain both of name and description fields")
	}

	if hasUnpleasentField {
		return hasName, hasDescription, fmt.Errorf("unexpected fields in product json")
	}

	return hasName, hasDescription, nil
}

func (ts *StoreServer) addProductHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Add product request handler\n")

	obtainedProduct, err := ts.processRequest(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hasName, hasDescription, err := checkProductFields(obtainedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !(hasName && hasDescription) {
		http.Error(w, "not enougth information about product", http.StatusBadRequest)
		return
	}

	iconPath, ok := obtainedProduct["icon"]
	if !ok {
		iconPath = ""
	}

	product := ts.store.AddProduct(obtainedProduct["name"], obtainedProduct["description"], iconPath)

	renderJSON(w, product)
}

func (ts *StoreServer) getAllProductsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Get all products request handler\n")

	allProducts := ts.store.GetAllProducts()
	renderJSON(w, allProducts)
}

func (ts *StoreServer) getProductHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("Get product request handler\n")

	product, err := ts.store.GetProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, product)

}

func (ts *StoreServer) updateProductHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("Update product request handler\n")

	obtainedProduct, err := ts.processRequest(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _, err = checkProductFields(obtainedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var product Product
	for k, v := range obtainedProduct {
		product, err = ts.store.UpdateProductData(id, k, v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	renderJSON(w, product)
}

func (ts *StoreServer) deleteProductHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("Delete product request handler\n")

	product, err := ts.store.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, product)

}

func (ts *StoreServer) getIconFromRequest(w http.ResponseWriter, req *http.Request, path string) {
	givenFile, _, err := req.FormFile("icon")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	localFile, err := os.Create("./" + path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, givenFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func (ts *StoreServer) addProductIconHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("Add product icon request handler\n")

	var iconPath = "resources/icon.png"
	ts.getIconFromRequest(w, req, iconPath)

	_, err := ts.store.AddProductIcon(id, iconPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, "")
}

func (ts *StoreServer) getProductIconHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("Get product icon request handler\n")

	product, err := ts.store.GetProductIcon(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendImage(w, product.Icon)
}

func (ts *StoreServer) TaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Started handling")

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	var id int
	var err error
	if len(pathParts) > 1 {
		id, err = strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	switch {
	case len(pathParts) == 1 && pathParts[0] == "product" && req.Method == http.MethodPost:
		ts.addProductHandler(w, req)
	case len(pathParts) == 1 && pathParts[0] == "products" && req.Method == http.MethodGet:
		ts.getAllProductsHandler(w, req)
	case len(pathParts) == 2 && pathParts[0] == "product" && req.Method == http.MethodGet:
		ts.getProductHandler(w, req, id)
	case len(pathParts) == 2 && pathParts[0] == "product" && req.Method == http.MethodPut:
		ts.updateProductHandler(w, req, id)
	case len(pathParts) == 2 && pathParts[0] == "product" && req.Method == http.MethodDelete:
		ts.deleteProductHandler(w, req, id)
	case len(pathParts) == 3 && pathParts[0] == "product" && pathParts[2] == "image" && req.Method == http.MethodPost:
		ts.addProductIconHandler(w, req, id)
	case len(pathParts) == 3 && pathParts[0] == "product" && pathParts[2] == "image" && req.Method == http.MethodGet:
		ts.getProductIconHandler(w, req, id)
	default:
		http.Error(w, fmt.Sprintf("Unexpected method for specified path. Path: %v, method: %v.", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
		return
	}
}
