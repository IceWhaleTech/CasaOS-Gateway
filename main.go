package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/IceWhaleTech/CasaOS-Gateway/route"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

func main() {

	err := loadConfig()
	if err != nil {
		panic(err)
	}

	settings := viper.AllSettings()
	log.Println("settings:", settings)

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
					return serve("management", ":0", route)
				})

				// gateway server
				g.Go(func() error {
					http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						proxy := management.GetProxy(r.URL.Path)
						proxy.ServeHTTP(w, r)
					})

					port := viper.GetString("gateway.port")
					addr := net.JoinHostPort("", port)

					return serve("gateway", addr, route)
				})

				return g.Wait()
			},
		})
}

func writeAddressFile(filename string, address string) error {
	path := viper.GetString("gateway.run-path")

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	filepath := filepath.Join(path, filename)
	return ioutil.WriteFile(filepath, []byte(address), 0644)
}

func loadConfig() error {
	viper.SetDefault("gateway.port", "8080")
	viper.SetDefault("gateway.run-path", "/var/run/casaos")

	viper.SetConfigName("gateway")
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("$HOME/.casaos")
	viper.AddConfigPath("/etc/casaos")

	return viper.ReadInConfig()
}

func serve(name string, addr string, route *gin.Engine) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	err = writeAddressFile(name+".address", listener.Addr().String())
	if err != nil {
		panic(err)
	}

	log.Println(name+" server listening on", listener.Addr().String())
	return http.Serve(listener, route)
}
