package main

import (
	"log"
	"net/http"
	"time"

	"github.com/IceWhaleTech/CasaOS-Gateway/route"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	// management server
	g.Go(func() error {
		managementServer := &http.Server{
			Addr:         ":8081",
			Handler:      route.Build(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		return managementServer.ListenAndServe()
	})

	// gateway server
	g.Go(func() error {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			route := r.URL.Path[1:]
			log.Println("route:", route)
		})
		return http.ListenAndServe(":8080", nil)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
