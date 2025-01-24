package route

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/samber/lo"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
)

type CaddyGateway struct {
	management *service.Management
}

func NewCaddyGateway(management *service.Management) *CaddyGateway {
	return &CaddyGateway{
		management: management,
	}
}

// ConfigureCaddy sets up Caddy with similar functionality to the Echo gateway
func (g *CaddyGateway) ConfigureCaddy() (*caddy.Config, error) {
	httpApp := caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{
			"zimaos": {
				Listen: []string{":80", ":443"},
				Routes: []caddyhttp.Route{
					// Automatic HTTPS redirect
					{
						HandlersRaw: []json.RawMessage{
							// caddyhttp.MatchExpression{
							// 			Match: `{http.request.scheme} == "http"`,
							// 		}.RawMsg(),
							// 		caddyhttp.StaticResponse{
							// 			StatusCode: http.StatusPermanentRedirect,
							// 			Headers:    map[string][]string{"Location": {"{http.request.scheme}s://{http.request.host}{http.request.uri}"}},
							// 		}.RawMsg(),
						},
					},

					// Health check route
					{
						HandlersRaw: []json.RawMessage{
							// caddyhttp.MatchPathRE{
							// 	Pattern: "/ping",
							// }.RawMsg(),
							// caddyhttp.StaticResponse{
							// 	StatusCode: http.StatusOK,
							// 	Body:       "pong from gateway service",
							// }.RawMsg(),
						},
					},

					// Main proxy route
					{
						HandlersRaw: []json.RawMessage{
							// 	// Custom middleware for IP handling
							// 	caddyhttp.MiddlewareHandler{
							// 		Handler: g.createIPMiddleware(),
							// 	}.RawMsg(),
							// 	// Reverse proxy handler
							// 	g.createReverseProxyHandler(),
						},
					},
				},
				AutoHTTPS: &caddyhttp.AutoHTTPSConfig{},
				// Enable HTTP/2 and HTTP/3
				Protocols: []string{"h1", "h2", "h3"},
			},
		},
	}

	// Create a new Caddy config
	config := &caddy.Config{
		Admin: &caddy.AdminConfig{
			Disabled: false,
			Config: &caddy.ConfigSettings{
				Persist: lo.ToPtr(false),
			},
		},
		AppsRaw: caddy.ModuleMap{
			"http": caddyconfig.JSON(httpApp, nil),
		},
	}

	return config, nil
}

// createIPMiddleware creates a middleware that handles IP forwarding similar to the Echo implementation
func (g *CaddyGateway) createIPMiddleware() caddyhttp.MiddlewareHandler {
	return &ipMiddleware{}
}

type ipMiddleware struct{}

func (m *ipMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	ipList := []string{}

	// when r.Header.Get("X-Forwarded-For") is "". the ipList should be empty.
	// fix https://github.com/IceWhaleTech/CasaOS/issues/1247
	if xff := r.Header.Get("X-Forwarded-For");if xff != "" {
		ipList = strings.Split(xff, ",")
		// when req.Header.Get("X-Forwarded-For") is "". to clean the ipList.
		// fix https://github.com/IceWhaleTech/CasaOS/issues/1247
		if len(ipList) == 1 && ipList[0] == "" {
			ipList = []string{}
		}
	}

	r.Header.Del("X-Forwarded-For")
	r.Header.Del("X-Real-IP")

	// Get the remote address
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}

	// Add the remote IP to the list if it's not empty
	if remoteIP != "" {
		ipList = append(ipList, remoteIP)
	}

	// Join the IPs back together
	if len(ipList) > 0 {
		r.Header.Set("X-Forwarded-For", strings.Join(ipList, ","))
	}

	return next.ServeHTTP(w, r)
}
