package main

import (
	"log"
	"net/http"
)

func main() {
	// http.HandleFunc("/posts/", servePost)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	select {} // a `select` statement at the end of the `main()` function. This is necessary to prevent the Go program from exiting, as the WebAssembly binary will be terminated when the Go program exits.
}
