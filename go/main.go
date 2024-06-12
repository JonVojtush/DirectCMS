package main

import (
	"log"
	"os"
	"path/filepath"
	"syscall/js"
)

var (
	document = js.Global().Get("document")
	console  = js.Global().Get("console")
	err      error
)

func getData(this js.Value, args []js.Value) interface{} {
	type Post struct {
		root *string
		name *string
		id   *uint
	}

	// Loop through the items that exist in the specified location. If its a directory, load data into objects, otherwise skip.
	if err = filepath.Walk("../posts", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			var post Post

			*post.root = path

			log.Print("Found post directory: " + path)
			console.Call("log", "Found post directory: "+path)
			return nil
		} else {
			log.Println("Skipping " + path + " as it is not a post directory.")
			console.Call("log", "Found post directory: "+path)
			return nil
		}
	}); err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	select {} // a `select` statement at the end of the `main()` function. This is necessary to prevent the Go program from exiting, as the WebAssembly binary will be terminated when the Go program exits.
}
