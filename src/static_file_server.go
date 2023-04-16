package main

import (
	"bufio"
	"errors"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/exp/slices"
)

// StaticFileServer for serving static content. It looks for files using the URL's path, relative
// to the current directory.
//
// Compression is supported via files with a .br suffix, if Accept-Encoding indicates "br" is
// accepted.
type StaticFileServer struct{}

func (s StaticFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// No need to check the path for unsafe ".." chars, as this is already handled by the http
	// mixer.
	mime := mime.TypeByExtension(path.Ext(r.URL.Path))
	// If a mime-type could be found from the extensions table, use it.
	// ResponseWriter defaults to setting a mime type based on the file contents otherwise.
	if mime != "" {
		w.Header().Add("Content-Type", mime)
	}
	localPath := "." + r.URL.Path

	if couldBeBrotli(mime) && acceptsBrotli(r.Header) {
		brPath := localPath + ".br"
		f, err := os.Open(brPath)
		if err == nil {
			defer f.Close()

			i, err := f.Stat()
			if err == nil && !i.IsDir() {
				w.Header().Add("Content-Encoding", "br")
				writeFile(w, f)
				return
			} else if err != nil {
				log.Printf("failed to stat file under path %v: %v", brPath, err.Error())
			}
		} else {
			if !errors.Is(err, fs.ErrNotExist) {
				log.Printf("failed to read file under path %v: %v", brPath, err.Error())
			}
		}
		// If any part of returning the Brotli version failed, fall through to try the uncompressed
		// file.
	}

	f, err := os.Open(localPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		log.Printf("Failed to open file at path %v: %v", localPath, err.Error())
		ioError(w)
		return
	}
	defer f.Close()

	i, err := f.Stat()
	if err != nil {
		log.Printf("Failed to stat file at path %v: %v", localPath, err.Error())
		ioError(w)
		return
	}
	if i.IsDir() {
		http.NotFound(w, r)
		return
	}

	writeFile(w, f)
}

func acceptsBrotli(headers http.Header) bool {
	accepts := headers["Accept-Encoding"]
	for _, joinedEncodings := range accepts {
		encodings := strings.Split(joinedEncodings, ",")
		for _, encoding := range encodings {
			if strings.TrimSpace(encoding) == "br" {
				return true
			}
		}
	}
	return false
}

var compressibleMimeTypes = []string{"application/json", "application/ld+json", "application/xml", "image/svg+xml"}

func couldBeBrotli(mimeType string) bool {
	return strings.HasPrefix(mimeType, "text/") || slices.Contains(compressibleMimeTypes, mimeType)
}

// writeFile writes the files data into the ResponseWriter.
//
// The ResponseWriter should not be used again after this returns.
func writeFile(w http.ResponseWriter, f *os.File) {
	r := bufio.NewReader(f)
	_, err := r.WriteTo(w)
	if err != nil {
		log.Printf("Failed to write contents: %v", err.Error())
		ioError(w)
	}
}

func ioError(w http.ResponseWriter) {
	http.Error(w, "I/O error", http.StatusInternalServerError)
}
