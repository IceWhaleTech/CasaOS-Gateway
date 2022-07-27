package service

import (
	"log"
	"net/http/httputil"
	"net/url"
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

func (g *Management) CreateRoute(route string, target string) {
	url, err := url.Parse(target)

	if err != nil {
		log.Fatalln(err)
	}

	g.pathTargetMap[route] = target
	g.pathReverseProxyMap[route] = httputil.NewSingleHostReverseProxy(url)
}

func (g *Management) GetRoutes() map[string]string {
	return g.pathTargetMap
}

func (g *Management) GetProxy(route string) *httputil.ReverseProxy {
	return g.pathReverseProxyMap[route]
}
