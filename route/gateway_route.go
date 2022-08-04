package route

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
)

type GatewayRoute struct {
	management *service.Management
}

func NewGatewayRoute(management *service.Management) *GatewayRoute {
	return &GatewayRoute{
		management: management,
	}
}

func (g *GatewayRoute) GetRoute() *http.ServeMux {
	gatewayMux := http.NewServeMux()
	gatewayMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy := g.management.GetProxy(r.URL.Path)

		if proxy == nil {
			if r.URL.Path == "/" {
				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusNotFound)
			return
		}

		proxy.ServeHTTP(w, r)
	})

	return gatewayMux
}
