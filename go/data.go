// TODO: Load content to post in newPost(). Then fetch within servePost() rather than fetching manually.

package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var postList []Post

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

	// List and serve media files in the post directory
	if mediaFiles, err := os.ReadDir(filepath.Join("posts", post.ID)); err == nil {
		for _, file := range mediaFiles {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".mp4") {
				// TODO: Build an array of media file names with extensions. 0 should always be featured.jpg ifelse featured.png, if no make [0] nil.
			}
		}
	}

	// TODO: post.MetaData =
	post.Content = filepath.Join("posts", post.ID, "content.md")

	return post
}

func buildPostList() {
	var err error

	if err = filepath.WalkDir("posts", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal("Error(1) walking the posts directory: ", err)
		}

		if entry.IsDir() {
			postTitle := filepath.Base(path)
			postList = append(postList, newPost(postTitle))
		}

		return nil
	}); err != nil {
		log.Fatal("Error(2) walking the posts directory: ", err)
	}
}
