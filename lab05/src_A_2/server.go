package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"os"
)

func setupConnection(server, port string) (*tls.Conn, error) {
	connTcp, err := net.Dial("tcp", server+":"+port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return &tls.Conn{}, err
	}

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn := tls.Client(connTcp, config)

	if conn.Handshake() != nil {
		fmt.Println("TLS handshake error:", err)
		return &tls.Conn{}, err
	}

	return conn, nil

}

func checkMailServerResponse(conn *tls.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading response from mail server:", err)
		os.Exit(1)
	}
	fmt.Println(string(buf[:n]))
}

func getRegularMessage(from, to, subject, body string) string {
	message := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n"
	return message
}

func getImageMessage(imagePath, from, to, subject string) (string, error) {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return "", err
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	message := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: image/png; name=\"icon.png\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"Content-Disposition: inline; filename=\"icon.png\"\r\n\r\n" +
		imageBase64 + "\r\n"

	return message, nil
}

func sendEmail(server, port, from, password, to, message string) {
	conn, err := setupConnection(server, port)
	if err != nil {
		return
	}
	defer conn.Close()

	checkMailServerResponse(conn)

	fmt.Fprintf(conn, "HELO localhost\r\n")
	checkMailServerResponse(conn)

	fmt.Fprintf(conn, "AUTH LOGIN\r\n")
	checkMailServerResponse(conn)

	fmt.Fprintf(conn, base64.StdEncoding.EncodeToString([]byte(from))+"\r\n")
	checkMailServerResponse(conn)

	fmt.Fprintf(conn, base64.StdEncoding.EncodeToString([]byte(password))+"\r\n")
	checkMailServerResponse(conn)

	fmt.Fprintf(conn, "MAIL FROM: <"+from+">\r\n")
	checkMailServerResponse(conn)

	recipients := []string{to}
	for _, recipient := range recipients {
		fmt.Fprintf(conn, "RCPT TO: <"+recipient+">\r\n")
		checkMailServerResponse(conn)
	}

	fmt.Fprintf(conn, "DATA\r\n")
	checkMailServerResponse(conn)

	fmt.Fprintf(conn, message+"\r\n.\r\n")
	checkMailServerResponse(conn)

	fmt.Fprintf(conn, "QUIT\r\n")
	checkMailServerResponse(conn)
}
