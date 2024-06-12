package main

import (
	_ "syscall/js"
)

/*var (
	document = js.Global().Get("document")
	console  = js.Global().Get("console")
	err error
)*/

func getData(this js.Value, args []js.Value) interface{} {
		//for posts in dir, build objecs
		//return postData
	}

func main() {
	select {} // a `select` statement at the end of the `main()` function. This is necessary to prevent the Go program from exiting, as the WebAssembly binary will be terminated when the Go program exits.
}
