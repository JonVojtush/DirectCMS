package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"syscall/js"
)

var (
	document        = js.Global().Get("document")
	mediaExtensions = []string{"jpg", "jpeg", "png", "gif", "webp", "mp4", "avi", "mov", "webm"}
)

// Function to check for featured image or video and display it at the top of the post content
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

	// Check if the postId contains featured media
	for _, ext := range mediaExtensions {
		re := regexp.MustCompile(`featured\.(\w+)`)
		if re.MatchString(*post.ID) {
			hasFeaturedMedia = true
			log.Println("Featured media found with" + ext + "extension.")
			break
		}
	}

	if hasFeaturedMedia {
		// Extract the featured media file name
		re := regexp.MustCompile(`featured\.(\w+)`)
		matches := re.FindStringSubmatch(*post.ID)
		if len(matches) > 1 {
			featuredImage := "featured." + matches[1]
			displayedContent = `<div id="post-media"><img src="/posts/` + *post.ID + `/` + featuredImage + `" alt="Featured Media"></div>`
		}
	}
	// Append the content to displayedContent
	if len(displayedContent) == 0 {
		displayedContent = `<div id="post-content">` + *post.Content + `</div>`
	} else {
		displayedContent += `<div id="post-content">` + *post.Content + `</div>`
	}
	postContainer.Set("innerHTML", displayedContent)
}

// isMediaFile checks if a file has one of the specified media extensions
func isMediaFile(fileName string) bool {
	for _, ext := range mediaExtensions {
		if strings.HasSuffix(fileName, "."+ext) {
			return true
		}
	}
	return false
}
