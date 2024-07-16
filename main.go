/*
	TODO: **Sitemap**: Generate a sitemap that lists all post URLs and media files. You can use tools like `go-sitemap-generator` to generate a sitemap dynamically.
	TODO: **Robots.txt**: Configure your server's `robots.txt` file to disallow crawling of media files but allow indexing of post pages.
			User-agent: *
			Disallow: /media/
			Allow: /posts/
			Sitemap: https://yourdomain.com/sitemap.xml
*/

package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func servePost(w http.ResponseWriter, r *http.Request) {
	// Extract postId from the request URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 { // BaseURL, posts, postid
		http.NotFound(w, r)
		return
	}
	postId := parts[2]

	// Construct the full path to the content and media files
	contentPath := filepath.Join("posts", postId, "content.md")
	mediaDir := filepath.Join("posts", postId)

	// Check if the content file exists
	_, err := os.Stat(contentPath)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Serve the markdown content
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	content, err := os.ReadFile(contentPath)
	if err != nil {
		http.Error(w, "Failed to read content file", http.StatusInternalServerError)
		return
	}
	w.Write(content)

	// List and serve media files in the post directory
	mediaFiles, err := os.ReadDir(mediaDir)
	if err == nil {
		for _, file := range mediaFiles {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".mp4") {
				http.ServeFile(w, r, filepath.Join(mediaDir, file.Name()))
			}
		}
	}
}

func main() {
	http.HandleFunc("/posts/", servePost)
	http.ListenAndServe(":8080", nil)
	select {} // a `select` statement at the end of the `main()` function. This is necessary to prevent the Go program from exiting, as the WebAssembly binary will be terminated when the Go program exits.
}
