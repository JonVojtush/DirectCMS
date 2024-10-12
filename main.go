package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	err error
	logFile *os.File
)

const logFileName = "appSession.log"

func main() {
	if logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666); err != nil {
		log.SetOutput(os.Stdout)
		log.Println("Failed to open log file: " + err.Error())
	} else {
		log.SetOutput(logFile)
	}
	defer logFile.Close()
	log.Println("----------------------------------------------------------------------------------------------------\nApp session started.")

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Failed to start server: " + err.Error())
	}

	http.HandleFunc("/", serveUserResources)
	
	select {} // Keep session open to serve functions over WASM.
}

// ---------- ACTIONS ----------
func serveUserResources(w http.ResponseWriter, r *http.Request) {
	basePath := "/user/"
	files := []string{"logo.*", "sitemap.xml", "custom.css", "custom.js", "favicon.*"}

	for _, pattern := range files {
		var (err error
			matches []string)

		if matches, err = filepath.Glob(filepath.Join(basePath, pattern)); err != nil {
			log.Println("Error checking existence of pattern " + pattern + ": " + err.Error())
		}

		for _, filePath := range matches {
			http.ServeFile(w, r, filePath)
		}
	}

	log.Println("None of the custom files exist.")
	http.NotFound(w, r)
}

// ---------- UTILITIES ----------