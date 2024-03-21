package main

import (
	"io"
	"net/smtp"
)

func sendEmail(from, to, password, subject, body, contentType string) {
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"Content-Type: " + contentType + "; charset=UTF-8\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil && err != io.EOF {
		panic(err)
	}
}
