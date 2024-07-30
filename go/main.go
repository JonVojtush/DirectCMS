package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall/js"
	"time"
)

var (
	err             error
	document        = js.Global().Get("document")
	postList        []Post
	imageExtensions = []string{"jpg", "jpeg", "png", "gif", "webp"}
	videoExtensions = []string{"mp4", "avi", "mov", "webm"}
	mediaExtensions = append(imageExtensions, videoExtensions...)
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func isMediaFile(fileName string) bool {
	for _, ext := range mediaExtensions {
		if strings.HasSuffix(fileName, "."+ext) {
			return true
		}
	}
	return false
}

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
	MetaData    *MetaData  `json:"MetaData"`
	Content     *string    `json:"Content"`
	Media       []*string  `json:"Media"`
}

func serveCustomResources(w http.ResponseWriter, r *http.Request) {
	basePath := "/custom/"
	files := []string{"logo.*", "sitemap.xml", "custom.css", "custom.js", "favicon.*"}

	for _, file := range files {
		filePath := filepath.Join(basePath, file)
		if exists, err := pathExists(filePath); err == nil && exists {
			http.ServeFile(w, r, filePath)
			return
		} else {
			log.Printf("Error checking existence of %s: %v", filePath, err)
		}
	}

	log.Println("None of the custom files exist")
	http.NotFound(w, r)
}

func newPost(postTitle string) Post {
	var (
		post           Post
		mediaFileNames []*string
		err            error
		metaFile       []byte
		metaData       MetaData
		contentBytes   []byte
	)

	// post.Title
	post.Title = &postTitle

	// post.ID
	postID := strings.ReplaceAll(strings.ToLower(*post.Title), " ", "%20")
	post.ID = &postID

	// post.LastUpdated
	if info, err := os.Stat("posts/" + *post.Title); err == nil {
		lastUpdated := info.ModTime()
		post.LastUpdated = &lastUpdated
	} else {
		log.Fatal("Could not read the content file: " + err.Error())
	}

	// post.MetaData
	metaFilePath := filepath.Join("posts", *post.ID, "meta.json")
	if metaFile, err = os.ReadFile(metaFilePath); err != nil {
		log.Fatalf("Could not read meta file for post %s: %v", *post.Title, err)
	}
	if err = json.Unmarshal(metaFile, &metaData); err != nil {
		log.Fatalf("Could not unmarshal meta data for post %s: %v", *post.Title, err)
	}
	post.MetaData = &metaData

	// post.Content
	contentPath := filepath.Join("posts", *post.ID, "content.md")
	if contentBytes, err = os.ReadFile(contentPath); err != nil {
		fmt.Println("Error reading file:", err)
	}
	contentString := string(contentBytes)
	post.Content = &contentString

	// post.Media
	mediaDirPath := filepath.Join("posts", *post.ID)
	if files, err := os.ReadDir(mediaDirPath); err == nil {
		for _, file := range files {
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

	return post
}

func buildPostList() {
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

func servePost(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		mediaFiles []fs.DirEntry
		parts      []string
		postIndex  int
	)

	if parts = strings.Split(r.URL.Path, "/"); len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	postId := parts[2]

	for i, post := range postList {
		if *post.ID == postId {
			postIndex = i
			break
		}
	}

	if postIndex == -1 {
		http.NotFound(w, r)
		return
	}

	mediaDir := filepath.Join("posts", postId)

	if mediaFiles, err = os.ReadDir(mediaDir); err == nil {
		for _, file := range mediaFiles {
			if !file.IsDir() && isMediaFile(file.Name()) {
				http.ServeFile(w, r, filepath.Join(mediaDir, file.Name()))
			}
		}
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	displayPost(postList[postIndex])
}

// TODO: buildNav()

func displayPost(post Post) {
	var (
		postContainer    js.Value
		hasFeaturedMedia bool
		displayedContent string
	)

	if postContainer = document.Call("getElementById", "post-container"); postContainer.IsUndefined() {
		fmt.Println("No container to display the post.")
		return
	}

	for _, ext := range mediaExtensions {
		re := regexp.MustCompile(`featured\.(\w+)`)
		if re.MatchString(*post.ID) {
			hasFeaturedMedia = true
			log.Println("Featured media found with " + ext + " extension.")
			break
		}
	}

	if hasFeaturedMedia {
		re := regexp.MustCompile(`featured\.(\w+)`)
		matches := re.FindStringSubmatch(*post.ID)
		if len(matches) > 1 {
			featuredImage := "featured." + matches[1]
			displayedContent = `<div id="featured"><img src="/posts/` + *post.ID + `/` + featuredImage + `" alt="Featured Media"></div>`
		}
	}
	if len(displayedContent) == 0 {
		displayedContent = `<div id="content">` + *post.Content + `</div>`
	} else {
		displayedContent += `<div id="content">` + *post.Content + `</div>`
	}
	postContainer.Set("innerHTML", displayedContent)
}

func main() {
	buildPostList()
	http.HandleFunc("/", serveCustomResources)
	http.HandleFunc("/posts/", servePost)
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	select {}
}
