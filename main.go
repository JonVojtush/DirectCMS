/*
	https://marketplace.visualstudio.com/items?itemName=aaron-bond.better-comments
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
	"syscall/js"
)

type Post struct {
	Title string `json:"Title"`
	ID    string `json:"ID"`
}

// Helper function to convert Post to JS Object
func goStruct2jsObject(post Post) js.Value {
	obj := js.Global().Get("Object").New()
	obj.Set("Title", js.ValueOf(post.Title))
	obj.Set("ID", js.ValueOf(post.ID))
	return obj
}

func fetchPostList() js.Value {
	array := js.Global().Get("Array").New() // Create an array to hold the objects
	//! populate array of structs with details pertaining to posts that exist in server-side posts directory
	//! for each post in posts: array.SetIndex(i, goStruct2jsObject(post))
	return array
}

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
	js.Global().Set("fetchPostList", func() js.Value { return fetchPostList() }) // Allow Javascript to call fetchPostList() which will return an array
	http.HandleFunc("/posts/", servePost)
	http.ListenAndServe(":8080", nil)
	/*
		! Advised by AI: select {} // a `select` statement at the end of the `main()` function. This is necessary to prevent the Go program from exiting, as the WebAssembly binary will be terminated when the Go program exits.
		? May not be necessary since Go is compiled to WASM. There would be no need to keep go running... not sure.
	*/
}
