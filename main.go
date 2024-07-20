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
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall/js"
	"time"
)

type Post struct {
	Title       string    `json:"Title"`
	ID          string    `json:"ID"`
	LastUpdated time.Time `json:"Updated"`
}

func newPost(postTitle string) Post {
	var (
		post            Post
		contentFileInfo fs.FileInfo
		err             error
	)

	post.Title = postTitle
	post.ID = strings.ReplaceAll(post.Title, " ", "_")

	if contentFileInfo, err = os.Stat("posts/" + post.Title); err != nil {
		log.Fatal("Could not read the file: " + err.Error())
	}
	post.LastUpdated = contentFileInfo.ModTime()

	return post
}

func convertPost2JS(post interface{}) js.Value {
	jsObj := js.Global().Get("Object").New()
	for key, value := range post.(map[string]interface{}) {
		jsObj.Set(key, js.ValueOf(value))
	}
	return jsObj
}

func builPostList() []Post {
	var (
		postList []Post
		err      error
	)

	if err = filepath.Walk("posts", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			postTitle := filepath.Base(path)
			postList = append(postList, newPost(postTitle))
		}

		return nil
	}); err != nil {
		log.Fatal("Error walking the posts directory: ", err)
	}

	return postList
}

func fetchPostList(postList []Post) js.Value {
	jsPostList := js.Global().Get("Array").New() // Create an array to hold the objects

	for _, post := range postList {
		jsPost := convertPost2JS(post)
		jsPostList.SetIndex(jsPostList.Length(), jsPost)
	}
	return jsPostList
}

func servePost(w http.ResponseWriter, r *http.Request) {
	var (
		content    []byte
		err        error
		mediaFiles []fs.DirEntry
	)

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
	if _, err = os.Stat(contentPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Serve the markdown content
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	if content, err = os.ReadFile(contentPath); err != nil {
		http.Error(w, "Failed to read content file", http.StatusInternalServerError)
		return
	}
	w.Write(content)

	// List and serve media files in the post directory
	if mediaFiles, err = os.ReadDir(mediaDir); err == nil {
		for _, file := range mediaFiles {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".mp4") {
				http.ServeFile(w, r, filepath.Join(mediaDir, file.Name()))
			}
		}
	}
}

func main() {
	postList := builPostList()
	js.Global().Set("fetchPostList", func() js.Value { return fetchPostList(postList) }) // Allow Javascript to call fetchPostList() which will return an array
	http.HandleFunc("/posts/", servePost)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	/*
		! Advised by AI: select {} // a `select` statement at the end of the `main()` function. This is necessary to prevent the Go program from exiting, as the WebAssembly binary will be terminated when the Go program exits.
		? May not be necessary since Go is compiled to WASM. There would be no need to keep go running... I am not sure.
	*/
}
