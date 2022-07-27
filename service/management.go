package service

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
)

type Management struct {
	pathTargetMap       map[string]string
	pathReverseProxyMap map[string]*httputil.ReverseProxy
}

func NewManagementService() *Management {
	return &Management{
		pathTargetMap:       make(map[string]string),
		pathReverseProxyMap: make(map[string]*httputil.ReverseProxy),
	}
}

func (g *Management) CreateRoute(route *common.Route) {
	url, err := url.Parse(route.Target)

	if err != nil {
		log.Fatalln(err)
	}

	g.pathTargetMap[route.Path] = route.Target
	g.pathReverseProxyMap[route.Path] = httputil.NewSingleHostReverseProxy(url)
}

func (g *Management) GetRoutes() []*common.Route {
	routes := make([]*common.Route, 0)

	for path, target := range g.pathTargetMap {
		routes = append(routes, &common.Route{
			Path:   path,
			Target: target,
		})
	}

	return routes
}

func (g *Management) GetProxy(route string) *httputil.ReverseProxy {
	return g.pathReverseProxyMap[route]
}
