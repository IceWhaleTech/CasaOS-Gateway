package service

import (
	"encoding/json"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/IceWhaleTech/CasaOS-Common/model"
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"go.uber.org/zap"
)

const RoutesFile = "routes.json"

type Management struct {
	pathTargetMap       map[string]string
	pathReverseProxyMap map[string]*httputil.ReverseProxy
	state               *State
}

func NewManagementService(state *State) *Management {
	routesFilepath := filepath.Join(state.GetRuntimePath(), RoutesFile)

	// try to load routes from routes.json
	pathTargetMap, err := loadPathTargetMapFrom(routesFilepath)
	if err != nil {
		logger.Error("Failed to load routes", zap.Any("error", err), zap.Any("filepath", routesFilepath))
		pathTargetMap = make(map[string]string)
	}

	pathReverseProxyMap := make(map[string]*httputil.ReverseProxy)

	for path, target := range pathTargetMap {
		targetURL, err := url.Parse(target)
		if err != nil {
			logger.Error("Failed to parse target", zap.Any("error", err), zap.String("target", target))
			continue
		}
		pathReverseProxyMap[path] = httputil.NewSingleHostReverseProxy(targetURL)
	}

	return &Management{
		pathTargetMap:       pathTargetMap,
		pathReverseProxyMap: pathReverseProxyMap,
		state:               state,
	}
}

func (g *Management) CreateRoute(route *model.Route) error {
	url, err := url.Parse(route.Target)
	if err != nil {
		return err
	}

	g.pathTargetMap[route.Path] = route.Target
	g.pathReverseProxyMap[route.Path] = httputil.NewSingleHostReverseProxy(url)

	routesFilePath := filepath.Join(g.state.GetRuntimePath(), RoutesFile)

	err = savePathTargetMapTo(routesFilePath, g.pathTargetMap)
	if err != nil {
		return err
	}

	return nil
}

func (g *Management) GetRoutes() []*model.Route {
	routes := make([]*model.Route, 0)

	for path, target := range g.pathTargetMap {
		routes = append(routes, &model.Route{
			Path:   path,
			Target: target,
		})
	}

	return routes
}

func (g *Management) GetProxy(path string) *httputil.ReverseProxy {
	// sort paths by length in descending order
	// (without this step, a path like "/abcd" can potentially be matched with "/ab")
	paths := getSortedKeys(g.pathReverseProxyMap)

	for _, p := range paths {
		if strings.HasPrefix(path, p) {
			return g.pathReverseProxyMap[p]
		}
	}
	return nil
}

func (g *Management) GetGatewayPort() string {
	return g.state.GetGatewayPort()
}

func (g *Management) SetGatewayPort(port string) error {
	if err := g.state.SetGatewayPort(port); err != nil {
		return err
	}

	return nil
}

func getSortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })

	return keys
}

func loadPathTargetMapFrom(routesFilepath string) (map[string]string, error) {
	content, err := os.ReadFile(routesFilepath)
	if err != nil {
		return nil, err
	}

	pathTargetMap := make(map[string]string)
	err = json.Unmarshal(content, &pathTargetMap)
	if err != nil {
		return nil, err
	}

	return pathTargetMap, nil
}

func savePathTargetMapTo(routesFilepath string, pathTargetMap map[string]string) error {
	content, err := json.Marshal(pathTargetMap)
	if err != nil {
		return err
	}

	return os.WriteFile(routesFilepath, content, 0o600)
}
