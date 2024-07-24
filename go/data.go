// TODO: Load content to post in buldpostList(). Then pass Post into servePost rather than fetching manually.

package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MetaData struct {
	Title       string   `json:"Title"`
	Description string   `json:"Description"`
	Keywords    []string `json:"Keywords"`
	Author      string   `json:"Author"`
}

type Post struct {
	Title       string    `json:"Title"`
	ID          string    `json:"ID"`
	LastUpdated time.Time `json:"Updated"`
	Media       []string  `json:"Media"`
	MetaData    MetaData  `json:"MetaData"`
	Content     string    `json:"Content"`
}

func newPost(postTitle string) Post {
	var (
		post            Post
		contentFileInfo fs.FileInfo
		err             error
	)

	post.Title = postTitle
	post.ID = strings.ReplaceAll(strings.ToLower(post.Title), " ", "%20") // match URL

	if contentFileInfo, err = os.Stat("posts/" + post.Title); err != nil {
		log.Fatal("Could not read the file: " + err.Error())
	}
	post.LastUpdated = contentFileInfo.ModTime()

	return post
}

func builPostList() []Post {
	var (
		postList []Post
		err      error
	)

	if err = filepath.Walk("posts", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal("Error(1) walking the posts directory: ", err)
		}

		if info.IsDir() {
			postTitle := filepath.Base(path)
			postList = append(postList, newPost(postTitle))
		}

		return nil
	}); err != nil {
		log.Fatal("Error(2) walking the posts directory: ", err)
	}

	return postList
}

func servePost(w http.ResponseWriter, r *http.Request) {
	var (
		content    []byte
		err        error
		mediaFiles []fs.DirEntry
		parts      []string
	)

	// Extract postId from the request URL
	if parts = strings.Split(r.URL.Path, "/"); len(parts) < 3 { // BaseURL, posts, postid
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
