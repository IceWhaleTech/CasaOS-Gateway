package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	route := map[string]string{
		"/v1/user": "http://localhost:8080/v1/user",
	}

	// loop thru the route map
	for path, upstreamURL := range route {
		// create a proxy for the upstreamURL
		target, err := url.Parse(upstreamURL)

		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(target)

		http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
			proxy.ServeHTTP(w, req)
		})

		log.Printf("path: %s, upstream: %s", path, upstreamURL)
	}

	log.Printf("listening on port %s", "3000")
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		panic(err)
	}
}
