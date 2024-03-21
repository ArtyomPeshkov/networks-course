package main

import "os"

func main() {
	// args[1] - from
	// args[2] - to
	// args[3] - gmail application password
	sendEmail(os.Args[1], os.Args[2], os.Args[3], "Halova?", "I have some text 4u.", "text/plain")
	sendEmail(os.Args[1], os.Args[2], os.Args[3], "Oh-oh, wait a min", "<html><body><h1>I also have some HTML :).</h1></body></html>", "text/html")
}
