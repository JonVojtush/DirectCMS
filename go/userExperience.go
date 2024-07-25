package main

import (
	"fmt"
	"net/http"
	"regexp"
	"syscall/js"
)

var document = js.Global().Get("document")

func serveCustom(w http.ResponseWriter, r *http.Request) {
	// TODO: Serve logo.*, sitemap.xml, custom.css, custom.js & logo.* to /root (web)
}

// TODO: function buildNav();

// Function to check for featured image or video and display it at the top of the post content
func displayPost(post Post) {
	var (
		postContainer    js.Value
		hasFeaturedMedia bool
		displayedContent string
		matched          bool
	)

	if postContainer = document.Call("getElementById", "post-container"); postContainer.IsUndefined() {
		fmt.Println("No container to display the post.")
		return
	}

	// Check if the postId contains featured media
	if matched, _ = regexp.MatchString(`featured\.(jpg|jpeg|png|gif|webp|mp4|avi|mov|webm)`, post.ID); matched {
		hasFeaturedMedia = true
	}

	if hasFeaturedMedia {
		// Extract the featured media file name
		re := regexp.MustCompile(`featured\.(jpg|jpeg|png|gif|webp|mp4|avi|mov|webm)`)
		featuredImage := re.FindString(post.ID)
		displayedContent = `<div id="post-media"><img src="/posts/` + post.ID + `/` + featuredImage + `" alt="Featured Media"></div>`
	}
	// Append the content to displayedContent
	displayedContent += `<div id="post-content">` + post.Content + `</div>`
	postContainer.Set("innerHTML", displayedContent)
}
