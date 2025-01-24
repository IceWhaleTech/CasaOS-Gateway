package route

import (
	"net/http"
	"strings"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type GatewayRoute struct {
	management *service.Management
}

func NewGatewayRoute(management *service.Management) *GatewayRoute {
	return &GatewayRoute{
		management: management,
	}
}

// rewriteRequestSourceIP ensures the request source IP is correct.
// we may receive two kinds of requests:
// 1. a request from reverse proxy
// 2. a request from client

// in reverse proxy, X-Forwarded-For will like:
// - `X-Forwarded-For:[192.168.6.102]`(normal)
// - `X-Forwarded-For:[::1, 192.168.6.102]`(hacked) Note: the ::1 is inject by attacker.
// - `X-Forwarded-For:[::1]`(normal or hacked) local request. But it from browser have JWT. So we can and need to verify it
// - `X-Forwarded-For:[::1,::1]`(normal or hacked) attacker can build the request to bypass the verification.
// But in the case. the remoteAddress should be the real ip. So we can use remoteAddress to verify it.
func rewriteRequestSourceIP() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			ipList := []string{}

			// when r.Header.Get("X-Forwarded-For") is "". the ipList should be empty.
			// fix https://github.com/IceWhaleTech/CasaOS/issues/1247
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				ipList = strings.Split(xff, ",")
				// when req.Header.Get("X-Forwarded-For") is "". to clean the ipList.
				// fix https://github.com/IceWhaleTech/CasaOS/issues/1247
				if len(ipList) == 1 && ipList[0] == "" {
					ipList = []string{}
				}
			}

			r.Header.Del("X-Forwarded-For")
			r.Header.Del("X-Real-IP")

			// Note: the X-Forwarded-For depend the correct config from reverse proxy.
			// otherwise the X-Forwarded-For may be empty.
			remoteIP := c.RealIP()
			if strings.Contains(remoteIP, ":") {
				remoteIP = remoteIP[:strings.LastIndex(remoteIP, ":")]
			}

			if len(ipList) > 0 && (remoteIP == "127.0.0.1" || remoteIP == "::1") {
				// to process the request from reverse proxy
				// in reverse proxy, X-Forwarded-For will container multiple IPs.
				// if the request is from reverse proxy, the remoteIP will be 127.0.0.1.
				// So we need get ip from X-Forwarded-For
				r.Header.Add("X-Forwarded-For", ipList[len(ipList)-1])
			}
			// to process the request from client.
			// the gateway will add the X-Forwarded-For to request header.
			// So we didn't need to add it.
			return next(c)
		}
	}
}

// GetEcho returns the configured Echo instance with all routes and middleware set up
func (g *GatewayRoute) GetRoute() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())

	e.GET("/ping", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "pong from gateway service")
	})

	// to fix https://github.com/IceWhaleTech/CasaOS/security/advisories/GHSA-32h8-rgcj-2g3c#event-102885
	// API V1 and V2 both read ip from request header. So the fix is effective for v1 and v2.
	e.Use(rewriteRequestSourceIP())

	// Handle all other requests through proxy
	e.Any("/*", func(ctx echo.Context) error {
		proxy := g.management.GetProxy(ctx.Request().URL.Path)

		if proxy == nil {
			return ctx.String(http.StatusNotFound, "not found")
		}

		proxy.ServeHTTP(ctx.Response(), ctx.Request())

		return nil
	})

	return e
}
