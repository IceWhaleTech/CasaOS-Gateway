package main

import (
	"context"
	"errors"
	"flag"
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
	"go.uber.org/fx"
)

const localhost = "127.0.0.1"

var (
	_state   *service.State
	_gateway *http.Server
)

func init() {
	versionFlag := flag.Bool("v", false, "version")
	flag.Parse()

	if *versionFlag {
		fmt.Println(common.Version)
		os.Exit(0)
	}

	_state = service.NewState()

	config, err := common.LoadConfig()
	if err != nil {
		panic(err)
	}

	if err := _state.SetRuntimePath(config.GetString(common.ConfigKeyRuntimePath)); err != nil {
		panic(err)
	}

	if err := _state.SetGatewayPort(config.GetString(common.ConfigKeyGatewayPort)); err != nil {
		panic(err)
	}

	if err := _state.SetWWWPath(config.GetString(common.ConfigKeyWWWPath)); err != nil {
		panic(err)
	}

	if err := checkPrequisites(_state); err != nil {
		panic(err)
	}

	_state.OnGatewayPortChange(func(s string) error {
		config.Set(common.ConfigKeyGatewayPort, _state.GetGatewayPort())
		config.Set(common.ConfigKeyRuntimePath, _state.GetRuntimePath())

		return config.WriteConfig()
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

				return management.CreateRoute(&common.Route{
					Path:   "/v1/gateway/port",
					Target: "http://" + listener.Addr().String(),
				})
			},
		},
	)

	// gateway service
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				route := gatewayRoute.GetRoute()

				if _state.GetGatewayPort() == "" {
					// check if a port is availble starting from port 80/8080
					portsToCheck := []int{}
					for i := 80; i < 90; i++ {
						portsToCheck = append(portsToCheck, i)
					}

					for i := 8080; i < 8090; i++ {
						portsToCheck = append(portsToCheck, i)
					}

					port := ""
					for _, p := range portsToCheck {
						port = fmt.Sprintf("%d", p)
						log.Printf("checking if port %s is available...", port)
						if listener, err := net.Listen("tcp", net.JoinHostPort("", port)); err == nil {
							if err = listener.Close(); err != nil {
								log.Printf("failed to close listener: %s", err)
								continue
							}
							break
						}
					}

					if port == "" {
						log.Fatalln("no port available for gateway to use")
					}

					if err := _state.SetGatewayPort(port); err != nil {
						return err
					}
				}

				_state.OnGatewayPortChange(func(port string) error {
					return reloadGateway(port, route)
				})

				return reloadGateway(_state.GetGatewayPort(), route)
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
	listener, err := net.Listen("tcp", net.JoinHostPort("", port))
	if err != nil {
		return err
	}

	addr := listener.Addr().String()

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
		err := gatewayNew.Serve(listener)
		if err != nil {
			log.Println(err)
		}
	}()

	// test if gateway is running
	url := "http://" + addr + "/ping"
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

	count := 10

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
