package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS-Gateway/route"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

const (
	ConfigKeyGatewayPort = "gateway.Port"
	ConfigKeyRuntimePath = "common.RuntimePath"
)

func main() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}

	err = checkPrequisites()
	if err != nil {
		panic(err)
	}

	pidFilename, err := writePidFile()
	if err != nil {
		panic(err)
	}

	defer cleanupFiles(
		pidFilename,
		service.RoutesFile,
		common.GatewayURLFilename,
		common.ManagementURLFilename,
	)

	ctx, cancel := context.WithCancel(context.Background())
	kill := make(chan os.Signal, 1)
	signal.Notify(kill)

	go func() {
		<-kill
		cancel()
	}()

	app := fx.New(
		fx.Provide(func() *service.Management {
			return service.NewManagementService(viper.GetString(ConfigKeyRuntimePath))
		}),
		fx.Provide(route.NewRoutes),
		fx.Invoke(run),
	)

	if err := app.Start(ctx); err != nil {
		log.Println(err)
	}
}

func run(lifecycle fx.Lifecycle, route *gin.Engine, management *service.Management) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				var g errgroup.Group

				// management server
				g.Go(func() error {
					return serve(common.ManagementURLFilename, "127.0.0.1:0", route)
				})

				// gateway server
				g.Go(func() error {
					gatewayMux := http.NewServeMux()
					gatewayMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						proxy := management.GetProxy(r.URL.Path)

						if proxy == nil {
							w.WriteHeader(http.StatusNotFound)
							return
						}

						proxy.ServeHTTP(w, r)
					})

					port := viper.GetString(ConfigKeyGatewayPort)
					addr := net.JoinHostPort("", port)

					return serve(common.GatewayURLFilename, addr, gatewayMux)
				})

				return g.Wait()
			},
		})
}

func writePidFile() (string, error) {
	runtimePath := viper.GetString(ConfigKeyRuntimePath)

	filename := "gateway.pid"
	filepath := filepath.Join(runtimePath, filename)
	return filename, ioutil.WriteFile(filepath, []byte(fmt.Sprintf("%d", os.Getpid())), 0o600)
}

func writeAddressFile(filename string, address string) (string, error) {
	path := viper.GetString(ConfigKeyRuntimePath)

	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return "", err
	}

	filepath := filepath.Join(path, filename)
	return filepath, ioutil.WriteFile(filepath, []byte(address), 0o600)
}

func cleanupFiles(filenames ...string) {
	RuntimePath := viper.GetString(ConfigKeyRuntimePath)

	for _, filename := range filenames {
		err := os.Remove(filepath.Join(RuntimePath, filename))
		if err != nil {
			log.Println(err)
		}
	}
}

func checkPrequisites() error {
	path := viper.GetString(ConfigKeyRuntimePath)

	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return fmt.Errorf("please ensure the owner of this service has write permission to that path %s", path)
	}

	return nil
}

func loadConfig() error {
	viper.SetDefault(ConfigKeyGatewayPort, "8080")
	viper.SetDefault(ConfigKeyRuntimePath, "/var/run/casaos") // See https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch05s13.html

	viper.SetConfigName("gateway")
	viper.SetConfigType("ini")

	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if configPath, success := os.LookupEnv("CASAOS_CONFIG_PATH"); success {
		viper.AddConfigPath(configPath)
	}

	viper.AddConfigPath(currentDirectory)
	viper.AddConfigPath(filepath.Join(currentDirectory, "conf"))
	viper.AddConfigPath(filepath.Join("/", "etc", "casaos"))

	return viper.ReadInConfig()
}

func serve(urlFilename string, addr string, route http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	url := "http://" + listener.Addr().String()

	urlFilePath, err := writeAddressFile(urlFilename, url)
	if err != nil {
		return err
	}

	log.Printf("listening on %s (saved to %s)", url, urlFilePath)

	s := &http.Server{
		Handler:           route,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s.Serve(listener)
}
