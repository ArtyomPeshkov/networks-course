package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	server := "smtp.gmail.com"
	port := "465"
	from := os.Args[1]     // from
	password := os.Args[2] // google application password
	to := os.Args[3]       // to
	subject := "Test Email"
	body := "This is a test email sent using a custom Go client."
	imagePath := "resources/icon.png"

	message := getRegularMessage(from, to, subject, body)
	messageImg, err := getImageMessage(imagePath, from, to, subject)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading file "+imagePath+":", err)
		os.Exit(1)
	}

	sendEmail(server, port, from, password, to, message)
	sendEmail(server, port, from, password, to, messageImg)

}
