package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/koyachi/go-nude"
        "github.com/valyala/fasthttp"
)

type memCacheType struct {
	cache map[string]bool
	m     sync.RWMutex
}

func is_url_okay(url_string string) (result bool){
	parsed, _ := url.Parse(url_string)
	result = fmt.Sprintf("%v://%v%v", parsed.Scheme, parsed.Host, parsed.Path) == url_string
	return
}

func (mc *memCacheType) Get(key string) (isNude bool, ok bool) {
	mc.m.RLock()
	isNude, ok = mc.cache[key]
	mc.m.RUnlock()
	return
}

func (mc *memCacheType) Set(key string, value bool) {
	mc.m.Lock()
	mc.cache[key] = value
	mc.m.Unlock()
}

var (
	// in memory cache
	memCache = &memCacheType{cache: make(map[string]bool)}

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

type resultType struct {
	IsNude bool `json:"isNude"`
	Status string `json:"status"`
}

func check_b64(query_string string) (result string){
	var isNude, ok bool
	var status string
	if query_string != "" {
		status = "error"
		link, err := base64.URLEncoding.DecodeString(query_string)
		if err == nil {
			link = bytes.TrimRight(link, "\n")
			if is_url_okay(string(link)) {
				status = "ok"
				key := fmt.Sprintf("%x", sha256.Sum256(link))
				isNude, ok = memCache.Get(key)
				if !ok {
					isNude = checkLinknudity(string(link))
					memCache.Set(key, isNude)
				}
			} else {
				log.Printf("encoded url should not contain query string")
			}
		} else {
			log.Printf("url decoding error %s: %s", query_string, err)
		}
	} else {
            status = "empty_query"
        }
	buf, err := json.Marshal(resultType{IsNude: isNude, Status: status})
	if err != nil {
		log.Printf("json encode error: %s", err)
	}
        result = string(buf)
        return
}

func main() {
	flag.Parse()
        requestHandler := func(ctx *fasthttp.RequestCtx) {
            ctx.SetContentType("application/json")
            ctx.SetStatusCode(fasthttp.StatusOK)
            query_string := string(ctx.QueryArgs().Peek("u"))
            fmt.Fprintf(ctx, "%s", check_b64(query_string))
        }
        fasthttp.ListenAndServe(*listenTo, requestHandler)
}
