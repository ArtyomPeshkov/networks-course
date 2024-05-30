package main

import (
	"log"
	"net/http"
)

/*
{
	"email": "artemyipeshkov@gmail.com",
	"password": "12345"
}

{
	"name": "Milk",
	"description": "Fresh"
}

{
	"name": "Candy",
	"description": "Sweet"
}


http://localhost:8081/user/sign-up
http://localhost:8081/user/sign-in

http://localhost:8081/product
http://localhost:8081/product?token=YXJ0ZW15aXBlc2hrb3ZAZ21haWwuY29tMTIzNDU+

http://localhost:8081/product/1
http://localhost:8081/product/2
http://localhost:8081/product/2?token=YXJ0ZW15aXBlc2hrb3ZAZ21haWwuY29tMTIzNDU+

http://localhost:8081/products
http://localhost:8081/products?token=YXJ0ZW15aXBlc2hrb3ZAZ21haWwuY29tMTIzNDU+

http://localhost:8081/product/2
http://localhost:8081/product/2?token=YXJ0ZW15aXBlc2hrb3ZAZ21haWwuY29tMTIzNDU+
http://localhost:8081/products?token=YXJ0ZW15aXBlc2hrb3ZAZ21haWwuY29tMTIzNDU+
*/

func main() {
	mux := http.NewServeMux()
	server := NewStoreServer()
	mux.HandleFunc("/", server.TaskHandler)

	if err := http.ListenAndServe("localhost:"+"8081", mux); err != http.ErrServerClosed {
		log.Print("Failed to run server at port 8080\nError: " + err.Error() + "\n")
	}
}
