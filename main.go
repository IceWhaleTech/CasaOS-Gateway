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
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	configKeyGatewayPort = "gateway.Port"
	configKeyWWWPath     = "gateway.WWWPath"
	configKeyRuntimePath = "common.RuntimePath"

	localhost          = "127.1"
	defaultGatewayPort = "80"
)

var (
	_state   *service.State
	_gateway *http.Server
)

func init() {
	_state = service.NewState()

	if err := loadConfig(_state); err != nil {
		panic(err)
	}

	if err := checkPrequisites(_state); err != nil {
		panic(err)
	}

	_state.OnGatewayPortChange(func(s string) error {
		return saveConfig(_state)
	})
}

func main() {
	pidFilename, err := writePidFile(_state.GetRuntimePath())
	if err != nil {
		panic(err)
	}

	defer cleanupFiles(
		_state.GetRuntimePath(),
		pidFilename, common.ManagementURLFilename, common.StaticURLFilename,
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
		fx.Provide(func() *service.State { return _state }),
		fx.Provide(service.NewManagementService),
		fx.Provide(route.NewManagementRoute),
		fx.Provide(route.NewGatewayRoute),
		fx.Provide(route.NewStaticRoute),
		fx.Invoke(run),
	)

	if err := app.Start(ctx); err != nil {
		log.Println(err)
	}
}

func run(
	lifecycle fx.Lifecycle,
	management *service.Management,
	managementRoute *route.ManagementRoute,
	gatewayRoute *route.GatewayRoute,
	staticRoute *route.StaticRoute,
) {
	// management server
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				listener, err := net.Listen("tcp", net.JoinHostPort(localhost, "0"))
				if err != nil {
					return err
				}

				managementServer := &http.Server{
					Handler:           managementRoute.GetRoute(),
					ReadHeaderTimeout: 5 * time.Second,
				}

				urlFilePath, err := writeAddressFile(_state.GetRuntimePath(), common.ManagementURLFilename, "http://"+listener.Addr().String())
				if err != nil {
					return err
				}

				go func() {
					log.Printf("management listening on %s (saved to %s)", listener.Addr().String(), urlFilePath)
					err := managementServer.Serve(listener)
					if err != nil {
						log.Fatalln(err)
					}
				}()

				return nil
			},
		},
	)

	// gateway service
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				route := gatewayRoute.GetRoute()

				if _state.GetGatewayPort() == "" {
					if err := _state.SetGatewayPort("80"); err != nil {
						return err
					}
				}

				_state.OnGatewayPortChange(func(port string) error {
					return reloadGateway(port, route)
				})

				if err := reloadGateway(_state.GetGatewayPort(), route); err != nil {
					return err
				}

				return management.CreateRoute(&common.Route{
					Path:   "/v1/gateway/port",
					Target: net.JoinHostPort(localhost, _state.GetGatewayPort()),
				})
			},
		})

	// static web
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", net.JoinHostPort(localhost, "0"))
			if err != nil {
				return err
			}

			staticServer := &http.Server{
				Handler:           staticRoute.GetRoute(),
				ReadHeaderTimeout: 5 * time.Second,
			}

			target := "http://" + listener.Addr().String()

			urlFilePath, err := writeAddressFile(_state.GetRuntimePath(), common.StaticURLFilename, target)
			if err != nil {
				return err
			}

			if err := management.CreateRoute(&common.Route{
				Path:   "/",
				Target: target,
			}); err != nil {
				return err
			}

			log.Printf("static server listening on %s (saved to %s)", listener.Addr().String(), urlFilePath)
			return staticServer.Serve(listener)
		},
	})
}

func reloadGateway(port string, route *http.ServeMux) error {
	addr := net.JoinHostPort("", port)

	if _gateway != nil && _gateway.Addr == addr {
		log.Println("port is the same as current running gateway - no change is required")
		return nil
	}

	// start new gateway
	gatewayNew := &http.Server{
		Addr:              addr,
		Handler:           route,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		err := gatewayNew.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	// test if gateway is running
	url := "http://" + gatewayNew.Addr
	if err := checkURL(url); err != nil {
		return err
	}

	log.Printf("gateway listening on %s", gatewayNew.Addr)

	// stop old gateway
	if _gateway != nil {
		log.Printf("stopping current gateway on %s", _gateway.Addr)
		if err := _gateway.Shutdown(context.Background()); err != nil {
			return err
		}
	}

	_gateway = gatewayNew

	return nil
}

func checkURL(url string) error {
	var response *http.Response
	var err error

	count := 3

	for count >= 0 {
		response, err = http.Get(url) //nolint:gosec

		if err == nil && response.StatusCode == http.StatusOK {
			break
		}

		time.Sleep(time.Second)

		count--
	}

	if (response == nil) || (response.StatusCode != http.StatusOK) {
		return errors.New("gateway not running")
	}

	return nil
}

func writePidFile(runtimePath string) (string, error) {
	filename := "gateway.pid"
	filepath := filepath.Join(runtimePath, filename)
	return filename, ioutil.WriteFile(filepath, []byte(fmt.Sprintf("%d", os.Getpid())), 0o600)
}

func writeAddressFile(runtimePath string, filename string, address string) (string, error) {
	err := os.MkdirAll(runtimePath, 0o755)
	if err != nil {
		return "", err
	}

	filepath := filepath.Join(runtimePath, filename)
	return filepath, ioutil.WriteFile(filepath, []byte(address), 0o600)
}

func cleanupFiles(runtimePath string, filenames ...string) {
	for _, filename := range filenames {
		err := os.Remove(filepath.Join(runtimePath, filename))
		if err != nil {
			log.Println(err)
		}
	}
}

func checkPrequisites(state *service.State) error {
	path := state.GetRuntimePath()

	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return fmt.Errorf("please ensure the owner of this service has write permission to that path %s", path)
	}

	return nil
}

func loadConfig(state *service.State) error {
	viper.SetDefault(configKeyGatewayPort, defaultGatewayPort)
	viper.SetDefault(configKeyWWWPath, "/var/lib/casaos/www")
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

	if err := state.SetWWWPath(viper.GetString(configKeyWWWPath)); err != nil {
		return err
	}

	return nil
}

func saveConfig(state *service.State) error {
	viper.Set(configKeyGatewayPort, state.GetGatewayPort())
	viper.Set(configKeyRuntimePath, state.GetRuntimePath())

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
