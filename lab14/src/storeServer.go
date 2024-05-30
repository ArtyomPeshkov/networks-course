package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	Token string
}

type RequsetProductStructure struct {
	Name        string
	Description string
	Icon        string
}

type StoreServer struct {
	currentId       int
	email_storage   map[string]string
	timer_storage   map[string]*time.Timer
	already_running map[string]bool
	main_store      *Store
	personal_stores map[string]*Store
}

func NewStoreServer() *StoreServer {
	store := new(Store)
	store.products = map[int]Product{}
	stores := map[string]*Store{}
	e_storage := map[string]string{}
	t_storage := map[string]*time.Timer{}
	a_running := map[string]bool{}
	return &StoreServer{main_store: store, personal_stores: stores, email_storage: e_storage, timer_storage: t_storage, already_running: a_running}
}

func sendEmail(to string) {
	from := "artemyipeshkov@gmail.com"
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Greetings from shop\n" +
		"Content-Type: text/plain; charset=UTF-8\n\n" +
		"We happy to see you in our shop!"

	apppass := "" // provide correct app password
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, apppass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil && err != io.EOF {
		panic(err)
	}
}

func (ts *StoreServer) runSender(ip string) {
	fmt.Println(ip)
	ts.timer_storage[ip].Reset(5 * time.Minute)
	if ts.already_running[ip] {
		return
	}
	fmt.Println("Started sending")
	ts.already_running[ip] = true
	go func() {
		<-ts.timer_storage[ip].C
		ts.already_running[ip] = false
		sendEmail(ts.email_storage[ip])
	}()
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
		http.Error(w, "Incorrect format in request.", http.StatusBadRequest)
		return map[string]string{}, err
	}

	return readJson, nil
}

func (ts *StoreServer) signHandler(w http.ResponseWriter, req *http.Request, action string) {
	log.Printf("Sign request handler\n")
	data, err := ts.processRequest(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, pass := data["email"], data["password"]
	token := base64.StdEncoding.EncodeToString([]byte(email + pass))
	token = strings.Replace(token, "?", "!", -1)
	token = strings.Replace(token, "=", "+", -1)
	ip := strings.Split(req.RemoteAddr, ":")[0]
	print(ip)

	if action == "sign-up" {
		_, ok := ts.personal_stores[token]
		if !ok {
			ts.personal_stores[token] = &Store{map[int]Product{}}
		}
		_, ok = ts.already_running[ip]
		if !ok {
			ts.already_running[ip] = false
		}
		_, ok = ts.timer_storage[ip]
		if !ok {
			ts.timer_storage[ip] = time.NewTimer(time.Second)
			ts.timer_storage[ip].Stop()
		}
		_, ok = ts.email_storage[ip]
		if !ok {
			ts.email_storage[ip] = email
		}
	} else if action == "sign-in" {
		_, ok := ts.personal_stores[token]
		if ok {
			renderJSON(w, Token{token})
		} else {
			http.Error(w, "User not found", http.StatusBadRequest)
		}
	}
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

func (ts *StoreServer) addProductHandler(w http.ResponseWriter, req *http.Request, token string) {
	log.Printf("Add product request handler\n")
	log.Printf(token + "\n")

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

	if token == "" {
		renderJSON(w, ts.main_store.AddProduct(obtainedProduct["name"], obtainedProduct["description"], iconPath, ts.currentId))
	} else {
		renderJSON(w, ts.personal_stores[token].AddProduct(obtainedProduct["name"], obtainedProduct["description"], iconPath, ts.currentId))
	}
	ts.currentId += 1
}

func (ts *StoreServer) getAllProductsHandler(w http.ResponseWriter, req *http.Request, token string) {
	log.Printf("Get all products request handler\n")
	if token == "" {
		ts.runSender(strings.Split(req.RemoteAddr, ":")[0])
	}

	allProducts := ts.main_store.GetAllProducts()
	if token != "" {
		allProducts = append(allProducts, ts.personal_stores[token].GetAllProducts()...)
	}
	renderJSON(w, allProducts)
}

func (ts *StoreServer) getProductHandler(w http.ResponseWriter, req *http.Request, id int, token string) {
	log.Printf("Get product request handler\n")
	if token == "" {
		ts.runSender(strings.Split(req.RemoteAddr, ":")[0])
	}
	var errPriv error = nil
	product, errPub := ts.main_store.GetProduct(id)
	if token != "" && errPub != nil {
		product, errPriv = ts.personal_stores[token].GetProduct(id)
	}
	if errPub != nil && (errPriv != nil || token == "") {
		http.Error(w, errPub.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, product)

}

func (ts *StoreServer) updateProductHandler(w http.ResponseWriter, req *http.Request, id int, token string) {
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
	var errPub error = nil
	var errPriv error = nil
	for k, v := range obtainedProduct {
		product, errPub = ts.main_store.UpdateProductData(id, k, v)
		if token != "" && errPub != nil {
			product, errPriv = ts.personal_stores[token].UpdateProductData(id, k, v)
		}
		if errPub != nil && (errPriv != nil || token == "") {
			http.Error(w, errPub.Error(), http.StatusBadRequest)
			return
		}
	}

	renderJSON(w, product)
}

func (ts *StoreServer) deleteProductHandler(w http.ResponseWriter, id int, token string) {
	log.Printf("Delete product request handler\n")

	var errPriv error = nil
	product, errPub := ts.main_store.DeleteProduct(id)
	if token != "" && errPub != nil {
		product, errPriv = ts.personal_stores[token].DeleteProduct(id)
	}
	if errPub != nil && (errPriv != nil || token == "") {
		http.Error(w, errPub.Error(), http.StatusBadRequest)
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

	_, err := ts.main_store.AddProductIcon(id, iconPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, "")
}

func (ts *StoreServer) getProductIconHandler(w http.ResponseWriter, id int) {
	log.Printf("Get product icon request handler\n")

	product, err := ts.main_store.GetProductIcon(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendImage(w, product.Icon)
}

func (ts *StoreServer) TaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Started handling\n")
	path := strings.Trim(req.RequestURI, "/")

	var token = ""
	print(req.URL.Path + "\n")
	if strings.Contains(path, "?token=") {
		pathAndToken := strings.Split(path, "?token=")
		path = pathAndToken[0]
		if len(pathAndToken) > 1 {
			token = pathAndToken[1]
		}
	}

	pathParts := strings.Split(path, "/")

	var id int
	var err error
	if len(pathParts) > 1 && pathParts[0] != "user" {
		id, err = strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	print(req.Method)
	switch {
	case len(pathParts) == 1 && pathParts[0] == "product" && req.Method == http.MethodPost:
		ts.addProductHandler(w, req, token)
	case len(pathParts) == 1 && pathParts[0] == "products" && req.Method == http.MethodGet:
		ts.getAllProductsHandler(w, req, token)
	case len(pathParts) == 2 && pathParts[0] == "user" && req.Method == http.MethodPost:
		ts.signHandler(w, req, pathParts[1])
	case len(pathParts) == 2 && pathParts[0] == "product" && req.Method == http.MethodGet:
		ts.getProductHandler(w, req, id, token)
	case len(pathParts) == 2 && pathParts[0] == "product" && req.Method == http.MethodPut:
		ts.updateProductHandler(w, req, id, token)
	case len(pathParts) == 2 && pathParts[0] == "product" && req.Method == http.MethodDelete:
		ts.deleteProductHandler(w, id, token)
	case len(pathParts) == 3 && pathParts[0] == "product" && pathParts[2] == "image" && req.Method == http.MethodPost:
		ts.addProductIconHandler(w, req, id)
	case len(pathParts) == 3 && pathParts[0] == "product" && pathParts[2] == "image" && req.Method == http.MethodGet:
		ts.getProductIconHandler(w, id)
	default:
		http.Error(w, fmt.Sprintf("Unexpected method for specified path. Path: %v, method: %v.", req.URL.Path, req.Method), http.StatusMethodNotAllowed)
		return
	}
}
