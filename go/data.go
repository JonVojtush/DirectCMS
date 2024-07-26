// TODO: Load content to post in newPost(). Then fetch within servePost() rather than fetching manually.

package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var postList []Post

type MetaData struct {
	Title       *string   `json:"Title"`
	Description *string   `json:"Description"`
	Keywords    []*string `json:"Keywords"`
	Author      *string   `json:"Author"`
}

type Post struct {
	Title       *string    `json:"Title"`
	ID          *string    `json:"ID"`
	LastUpdated *time.Time `json:"Updated"`
	Media       []*string  `json:"Media"`
	MetaData    *MetaData  `json:"MetaData"`
	Content     *string    `json:"Content"`
}

func newPost(postTitle string) Post {
	var (
		post            Post
		contentFileInfo os.FileInfo
		err             error
		mediaFileNames  []*string
	)

	post.Title = &postTitle
	postID := strings.ReplaceAll(strings.ToLower(*post.Title), " ", "%20") // match URL
	post.ID = &postID

	if contentFileInfo, err = os.Stat("posts/" + *post.Title); err != nil {
		log.Fatal("Could not read the file: " + err.Error())
	}
	lastUpdated := contentFileInfo.ModTime()
	post.LastUpdated = &lastUpdated

	mediaDirPath := filepath.Join("posts", *post.ID, "content.md") // Corrected to join multiple strings

	if mediaFiles, err := os.ReadDir(mediaDirPath); err == nil { // Fixed the path in ReadDir
		for _, file := range mediaFiles {
			if !file.IsDir() && (strings.HasSuffix(file.Name(), ".jpg") ||
				strings.HasSuffix(file.Name(), ".png") ||
				strings.HasSuffix(file.Name(), ".mp4")) {
				fileName := file.Name()
				mediaFileNames = append(mediaFileNames, &fileName)
			}
		}

		if len(mediaFileNames) > 0 {
			featuredIndex := -1
			for i, fileName := range mediaFileNames {
				if *fileName == "featured.jpg" || *fileName == "featured.png" {
					featuredIndex = i
					break
				}
			}

			if featuredIndex != -1 {
				temp := *mediaFileNames[featuredIndex]
				mediaFileNames[0], mediaFileNames[featuredIndex] = &temp, nil
			} else {
				mediaFileNames[0] = nil
			}
		} else {
			mediaFileNames = []*string{nil}
		}
	} else {
		fmt.Println("Error reading directory:", err)
	}
	post.Media = mediaFileNames
	post.Content = filepath.Join("posts", *post.ID, "content.md")

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
