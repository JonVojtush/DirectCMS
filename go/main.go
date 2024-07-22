package main

import (
	"log"
	"net/http"
	"syscall/js"
)

func main() {
	postList := builPostList()
	js.Global().Set("fetchPostList", func() js.Value { return fetchPostList(postList) }) // Allow Javascript to call fetchPostList() which will return an array
	http.HandleFunc("/posts/", servePost)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	select {} // a `select` statement at the end of the `main()` function. This is necessary to prevent the Go program from exiting, as the WebAssembly binary will be terminated when the Go program exits.
}
