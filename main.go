package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/IceWhaleTech/CasaOS-Gateway/route"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	kill := make(chan os.Signal, 1)
	signal.Notify(kill)

	go func() {
		<-kill
		cancel()
	}()

	app := fx.New(
		fx.Provide(service.NewManagementService),
		fx.Provide(route.Build),
		fx.Invoke(run),
	)

	if err := app.Start(ctx); err != nil {
		log.Fatalln(err)
	}
}

func run(lifecycle fx.Lifecycle, route *gin.Engine, management *service.Management) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				var g errgroup.Group

				// management server
				g.Go(func() error {
					managementServer := &http.Server{
						Addr:         ":8081",
						Handler:      route,
						ReadTimeout:  10 * time.Second,
						WriteTimeout: 10 * time.Second,
					}
					return managementServer.ListenAndServe()
				})

				// gateway server
				g.Go(func() error {
					http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						route := r.URL.Path[1:]
						target := management.GetRoute(route)
						log.Println("route:", route, "target:", target)
					})
					return http.ListenAndServe(":8080", nil)
				})

				return g.Wait()
			},
		})
}
