package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS-Gateway/route"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	configKeyGatewayPort = "gateway.Port"
	configKeyRuntimePath = "common.RuntimePath"
)

var (
	_state   *service.State
	_gateway *http.Server
)

func init() {
	_state = service.NewState()

	if err := load(_state); err != nil {
		panic(err)
	}

	if err := checkPrequisites(); err != nil {
		panic(err)
	}
}

func main() {
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

	defer func() {
		if _gateway != nil {
			if err := _gateway.Shutdown(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-kill
		cancel()
	}()

	app := fx.New(
		fx.Provide(func() *service.Management {
			return service.NewManagementService(_state)
		}),
		fx.Provide(route.NewRoutes),
		fx.Invoke(run),
	)

	if err := app.Start(ctx); err != nil {
		log.Println(err)
	}
}

func run(
	lifecycle fx.Lifecycle,
	route *gin.Engine,
	management *service.Management,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				// gateway service
				gatewayMux := http.NewServeMux()
				gatewayMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/" {
						if _, err := w.Write([]byte("TODO")); err != nil {
							log.Println(err)
						}
						return
					}

					proxy := management.GetProxy(r.URL.Path)

					if proxy == nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}

					proxy.ServeHTTP(w, r)
				})

				if _state.GetGatewayPort() == "" {
					if err := _state.SetGatewayPort("80"); err != nil {
						return err
					}
				}

				if err := reloadGateway(_state.GetGatewayPort(), gatewayMux); err != nil {
					return err
				}

				_state.OnGatewayPortChange(func(port string) error {
					return reloadGateway(port, gatewayMux)
				})

				// management server
				s := &http.Server{
					Handler:           route,
					ReadHeaderTimeout: 5 * time.Second,
				}

				return start(s, common.ManagementURLFilename, "127.0.0.1:0")
			},
		})
}

func start(s *http.Server, urlFilename string, addr string) error {
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

	return s.Serve(listener)
}

func reloadGateway(port string, route *http.ServeMux) error {
	addr := net.JoinHostPort("", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// start new gateway
	url := "http://" + listener.Addr().String()

	urlFilePath, err := writeAddressFile(common.GatewayURLFilename, url)
	if err != nil {
		return err
	}

	gatewayNew := &http.Server{
		Handler:           route,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s (saved to %s)", url, urlFilePath)
		err := gatewayNew.Serve(listener)
		if err != nil {
			log.Println(err)
		}
	}()

	// check if gatewayNew is running
	response, err := http.Get(url) //nolint:gosec
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("gatewayNew is not running")
	}

	// stop old gateway
	if _gateway != nil {
		if err := _gateway.Shutdown(context.Background()); err != nil {
			return err
		}
	}

	_gateway = gatewayNew

	return nil
}

func writePidFile() (string, error) {
	runtimePath := _state.GetRuntimePath()

	filename := "gateway.pid"
	filepath := filepath.Join(runtimePath, filename)
	return filename, ioutil.WriteFile(filepath, []byte(fmt.Sprintf("%d", os.Getpid())), 0o600)
}

func writeAddressFile(filename string, address string) (string, error) {
	runtimePath := _state.GetRuntimePath()

	err := os.MkdirAll(runtimePath, 0o755)
	if err != nil {
		return "", err
	}

	filepath := filepath.Join(runtimePath, filename)
	return filepath, ioutil.WriteFile(filepath, []byte(address), 0o600)
}

func cleanupFiles(filenames ...string) {
	runtimePath := _state.GetRuntimePath()

	for _, filename := range filenames {
		err := os.Remove(filepath.Join(runtimePath, filename))
		if err != nil {
			log.Println(err)
		}
	}
}

func checkPrequisites() error {
	path := _state.GetRuntimePath()

	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return fmt.Errorf("please ensure the owner of this service has write permission to that path %s", path)
	}

	return nil
}

func load(state *service.State) error {
	viper.SetDefault(configKeyGatewayPort, "80")
	viper.SetDefault(configKeyRuntimePath, "/var/run/casaos") // See https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch05s13.html

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

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := state.SetRuntimePath(viper.GetString(configKeyRuntimePath)); err != nil {
		return err
	}

	if err := state.SetGatewayPort(viper.GetString(configKeyGatewayPort)); err != nil {
		return err
	}

	return nil
}

func save(state *service.State) error {
	viper.Set(configKeyGatewayPort, state.GetGatewayPort())
	viper.Set(configKeyRuntimePath, state.GetRuntimePath())

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
