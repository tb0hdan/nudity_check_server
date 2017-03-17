package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/koyachi/go-nude"
)

var (
	// in memory cache
	memCache = map[string]bool{}

	// command line flags
	listenTo = flag.String("listen-to", "localhost:8000", "bind address and port")
)

// fetchLink download image by link and returns it
func fetchLink(link string) (img image.Image, err error) {
	r, err := http.DefaultClient.Get(link)
	if err != nil {
		return
	}

	defer r.Body.Close()

	img, _, err = image.Decode(r.Body)
	return
}

// checkLinknudity checks that passed link is nude
func checkLinknudity(link string) (isNude bool) {
	img, err := fetchLink(link)
	if err != nil {
		log.Printf("fetch url error %s: %s", link, err)
		return
	}

	isNude, err = nude.IsImageNude(img)
	if err != nil {
		log.Printf("check nudity error: %s", err)
	}
	return
}

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var isNude, ok bool

		w.Header().Set("Content-type", "application/json")

		u := r.URL.Query().Get("u")
		if u != "" {
			link, err := base64.URLEncoding.DecodeString(u)
			if err == nil {
				key := fmt.Sprintf("%x", sha256.Sum256(link))
				isNude, ok = memCache[key]
				if !ok {
					isNude = checkLinknudity(string(link))
					memCache[key] = isNude
				}
			} else {
				log.Printf("url decoding error %s: %s", u, err)
			}
		}

		fmt.Fprintf(w, "{\"isNude\": \"%v\"}", isNude)
	})
	http.ListenAndServe(*listenTo, nil)
}
