package route

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/IceWhaleTech/CasaOS-Common/model"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
)

const (
	ROUTES_FILE = "/var/run/casaos/routes.json"
	STATIC_FILE = "/usr/share/casaos/www"
)

type CaddyGateway struct {
	management *service.Management
	state      *service.State
}

func NewCaddyGateway(management *service.Management, state *service.State) *CaddyGateway {
	return &CaddyGateway{
		management: management,
		state:      state,
	}
}

func (g *CaddyGateway) handleStatic() string {
	return `
handle {
	root * /usr/share/casaos/www
	encode gzip

	@index {
		path_regexp /($|modules/[^\/]*/($|(index\.(html?|aspx?|cgi|do|jsp))|((default|index|home)\.php)))
	}
	header @index Cache-Control "no-cache, no-store, must-revalidate, proxy-revalidate, max-age=0"

	file_server
}
`
}

func (g *CaddyGateway) handleReverseProxy() string {
	data, err := os.ReadFile(ROUTES_FILE)
	if err != nil {
		return ""
	}

	var routes map[string]string
	err = json.Unmarshal(data, &routes)
	if err != nil {
		return ""
	}

	routesList := make([]model.Route, 0)
	for path, target := range routes {
		target = strings.TrimPrefix(target, "http://")

		// skip root and docs
		if path == "/" || strings.HasPrefix(path, "/doc") || strings.HasPrefix(path, "/v1doc") {
			continue
		}

		routesList = append(routesList, model.Route{
			Path:   path,
			Target: target,
		})
	}

	result := ""
	for _, route := range routesList {
		result += `
handle ` + route.Path + "*" + ` {
  reverse_proxy ` + route.Target + `
}
`
	}

	return result
}

// ConfigureCaddy sets up Caddy with similar functionality to the Echo gateway
func (g *CaddyGateway) BuildZimaCaddyfile() string {
	domain := g.state.GetGatewayDomain()

	return `` + domain + `
` + g.handleStatic() + g.handleReverseProxy() + ``
}

// ConfigureCaddy sets up Caddy with similar functionality to the Echo gateway
func (g *CaddyGateway) BuildCaddyfile() string {
	return `
	{
	}

	import *.caddyfile
	`
}
