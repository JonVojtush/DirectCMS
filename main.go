package main

import (
	_ "encoding/json"
	_ "fmt"
	_ "io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	_ "regexp"
	_ "strings"
	_ "syscall/js"
	_ "time"
)

var (
	err     error
	logFile *os.File
	// document = js.Global().Get("document")
	// postList        []Post
	// imageExtensions = []string{"jpg", "jpeg", "png", "gif", "webp"}
	// videoExtensions = []string{"mp4", "avi", "mov", "webm"}
	// mediaExtensions = append(imageExtensions, videoExtensions...)
)

/* type MetaData struct {
	Title       *string   `json:"Title"`
	Description *string   `json:"Description"`
	Keywords    []*string `json:"Keywords"`
	Author      *string   `json:"Author"`
} */

/* type Post struct {
	Title       *string    `json:"Title"`
	ID          *string    `json:"ID"`
	LastUpdated *time.Time `json:"Updated"`
	MetaData    *MetaData  `json:"MetaData"`
	Content     *string    `json:"Content"`
	Media       []*string  `json:"Media"`
} */

func serveCustomResources(w http.ResponseWriter, r *http.Request) {
	basePath := "/custom/"
	files := []string{"logo.*", "sitemap.xml", "custom.css", "custom.js", "favicon.*"}

	for _, pattern := range files {
		matches, err := filepath.Glob(filepath.Join(basePath, pattern))
		if err != nil {
			log.Println("Error checking existence of pattern " + pattern + ": " + err.Error())
			continue
		}

		for _, filePath := range matches {
			http.ServeFile(w, r, filePath)
			return
		}
	}

	log.Println("None of the custom files exist.")
	http.NotFound(w, r)
}

/* func newPost(postTitle string) Post {
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
	if info, err := os.Stat("/user/posts/" + *post.Title); err == nil {
		lastUpdated := info.ModTime()
		post.LastUpdated = &lastUpdated
	} else {
		log.Fatal("Could not read the post directory file: " + err.Error())
	}

	// post.MetaData
	metaFilePath := filepath.Join("posts", *post.ID, "meta.json")
	if metaFile, err = os.ReadFile(metaFilePath); err != nil {
		log.Println("Could not read meta file for " + *post.Title + ": " + err.Error())
	}
	if err = json.Unmarshal(metaFile, &metaData); err != nil {
		log.Println("Could not unmarshal meta data for " + *post.Title + ": " + err.Error())
	}
	post.MetaData = &metaData

	// post.Content
	contentPath := filepath.Join("posts", *post.ID, "content.html")
	if contentBytes, err = os.ReadFile(contentPath); err != nil {
		fmt.Println("Error reading file:", err)
	}
	contentString := string(contentBytes)
	post.Content = &contentString

	// post.Media
	mediaDirPath := filepath.Join("posts", *post.ID)
	// Read media files associated with the post from the media directory
	if files, err := os.ReadDir(mediaDirPath); err == nil {
		for _, file := range files {
			// Check if the file is not a directory and has a valid media file extension
			if !file.IsDir() && (strings.HasSuffix(file.Name(), ".jpg") ||
				strings.HasSuffix(file.Name(), ".png") ||
				strings.HasSuffix(file.Name(), ".mp4")) {
				fileName := file.Name()                            // Get the file name
				mediaFileNames = append(mediaFileNames, &fileName) // Append the file name to the media file names slice
			}
		}
		// Handle featured media files if they exist
		if len(mediaFileNames) > 0 {
			featuredIndex := -1 // Initialize index for featured media
			for i, fileName := range mediaFileNames {
				// Check if the file is a featured media file
				if *fileName == "featured.jpg" || *fileName == "featured.png" {
					featuredIndex = i // Store the index of the featured file
					break
				}
			}
			// If a featured file is found, move it to the front of the slice
			if featuredIndex != -1 {
				temp := *mediaFileNames[featuredIndex]
				mediaFileNames[0], mediaFileNames[featuredIndex] = &temp, nil // Swap featured file with the first element
			} else {
				mediaFileNames[0] = nil // If no featured file, set the first element to nil
			}
		} else {
			mediaFileNames = []*string{nil} // If no media files, initialize with nil
		}
	} else {
		fmt.Println("Error reading directory:", err)
	}
	post.Media = mediaFileNames

	return post
} */

/* func buildPostList() {
	// Walk through the posts directory and create a post object for each directory
	if err = filepath.WalkDir("posts", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Println("Error(1) walking the posts directory: " + err.Error())
		}

		if entry.IsDir() {
			postTitle := filepath.Base(path)
			postList = append(postList, newPost(postTitle))
		}

		return nil
	}); err != nil {
		log.Println("Error(2) walking the posts directory: " + err.Error())
	}
} */

/* func servePost(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		mediaFiles []fs.DirEntry
		parts      []string
		postIndex  int
	)

	// Split the URL path to extract the post ID
	if parts = strings.Split(r.URL.Path, "/"); len(parts) < 3 {
		http.NotFound(w, r) // Return 404 if the URL path is invalid
		return
	}
	postId := parts[2] // Extract the post ID from the URL

	// Find the index of the post in the post list based on the post ID
	for i, post := range postList {
		if *post.ID == postId {
			postIndex = i // Store the index if the post ID matches
			break
		}
	}

	if postIndex == -1 {
		http.NotFound(w, r) // Return 404 if the post is not found
		return
	}

	// Read the media files from the post's media directory
	mediaDir := filepath.Join("posts", postId)
	if mediaFiles, err = os.ReadDir(mediaDir); err == nil {
		for _, file := range mediaFiles {
			// Check if the file is not a directory and is a valid media file
			if !file.IsDir() && isMediaFile(file.Name()) {
				http.ServeFile(w, r, filepath.Join(mediaDir, file.Name())) // Serve the media file
			}
		}
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError) // Return 500 if the media directory cannot be read
		return
	}

	displayPost(postList[postIndex]) // Display the post content
} */

// TODO: buildNav()

/* func displayPost(post Post) {
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
} */

func main() {
	if logFile, err = os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666); err != nil {
		log.Println("Failed to open or create log file: " + err.Error())
		log.SetOutput(os.Stdout) // Fallback to stdout if log file can't be opened
		return
	}
	defer logFile.Close()  // Ensure the file is closed when done
	log.SetOutput(logFile) // Set the log output to the log file

	// buildPostList()

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Failed to start server: " + err.Error())
	}
	http.HandleFunc("/", serveCustomResources)
	// http.HandleFunc("/posts/", servePost)
	select {}
}

// ---------- UTILITIES ----------

/* func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
} */

/* func isMediaFile(fileName string) bool {
	for _, ext := range mediaExtensions {
		if strings.HasSuffix(fileName, "."+ext) {
			return true
		}
	}
	return false
} */
