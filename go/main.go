/*
	! **Sitemap**: Generate a sitemap that lists all post URLs and media files. You can use tools like `go-sitemap-generator` to generate a sitemap dynamically.
	! **Robots.txt**: Configure your server's `robots.txt` file to disallow crawling of media files but allow indexing of post pages.
			User-agent: *
			Disallow: /media/
			Allow: /posts/
			Sitemap: https://yourdomain.com/sitemap.xml
*/

package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"syscall/js"
)

var (
	document = js.Global().Get("document")
	console  = js.Global().Get("console")
	err      error
)

func serveMedia(w http.ResponseWriter, r *http.Request) {
	// Extract postId and imageName from the request URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	postId := parts[2]
	imageName := parts[3]

	// Construct the full path to the media file
	mediaPath := filepath.Join("posts", postId, imageName)

	// Serve the media file securely
	http.ServeFile(w, r, mediaPath)
}

func servePostContent(w http.ResponseWriter, r *http.Request) {
	// Extract postId from the request URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	postId := parts[2]

	// Construct the full path to the content file
	contentPath := filepath.Join("posts", postId, "content.md")

	// Read and serve the markdown file
	content, err := ioutil.ReadFile(contentPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Write(content)
}

func main() {
	http.HandleFunc("/media/", serveMedia)
	http.HandleFunc("/posts/", servePostContent)
	http.ListenAndServe(":8080", nil)
	select {} // a `select` statement at the end of the `main()` function. This is necessary to prevent the Go program from exiting, as the WebAssembly binary will be terminated when the Go program exits.
}
