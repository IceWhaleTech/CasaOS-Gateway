package route

import (
	"net/http"
	"strings"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"go.uber.org/zap"
)

type GatewayRoute struct {
	management *service.Management
}

func NewGatewayRoute(management *service.Management) *GatewayRoute {
	return &GatewayRoute{
		management: management,
	}
}

// the function is to ensure the request source IP is correct.
func rewriteRequestSourceIP(r *http.Request) {
	// we may receive two kinds of requests. a request from reverse proxy. a request from client.

	// in reverse proxy, X-Forwarded-For will like
	// `X-Forwarded-For:[192.168.6.102]`(normal)
	// `X-Forwarded-For:[::1, 192.168.6.102]`(hacked) Note: the ::1 is inject by attacker.
	// `X-Forwarded-For:[::1]`(normal or hacked) local request. But it from browser have JWT. So we can and need to verify it
	// `X-Forwarded-For:[::1,::1]`(normal or hacked) attacker can build the request to bypass the verification.
	// But in the case. the remoteAddress should be the real ip. So we can use remoteAddress to verify it.

	ipList := strings.Split(r.Header.Get("X-Forwarded-For"), ",")

	// when r.Header.Get("X-Forwarded-For") is "". to clean the ipList.
	// fix https://github.com/IceWhaleTech/CasaOS/issues/1247
	if len(ipList) == 1 && ipList[0] == "" {
		ipList = []string{}
	}

	r.Header.Del("X-Forwarded-For")
	r.Header.Del("X-Real-IP")

	// Note: the X-Forwarded-For depend the correct config from reverse proxy.
	// otherwise the X-Forwarded-For may be empty.
	remoteIP := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
	if len(ipList) > 0 && (remoteIP == "127.0.0.1" || remoteIP == "::1") {
		// to process the request from reverse proxy

		// in reverse proxy, X-Forwarded-For will container multiple IPs.
		// if the request is from reverse proxy, the r.RemoteAddr will be 127.0.0.1.
		// So we need get ip from X-Forwarded-For
		r.Header.Add("X-Forwarded-For", ipList[len(ipList)-1])
	}
	// to process the request from client.
	// the gateway will add the X-Forwarded-For to request header.
	// So we didn't need to add it.
}

func (g *GatewayRoute) GetRoute() *http.ServeMux {
	gatewayMux := http.NewServeMux()
	gatewayMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("pong from gateway service")); err != nil {
				logger.Error("Failed to `pong` in resposne to `ping`", zap.Any("error", err))
			}
			return
		}

		proxy := g.management.GetProxy(r.URL.Path)

		if proxy == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// to fix https://github.com/IceWhaleTech/CasaOS/security/advisories/GHSA-32h8-rgcj-2g3c#event-102885
		// API V1 and V2 both read ip from request header. So the fix is effective for v1 and v2.
		rewriteRequestSourceIP(r)

		proxy.ServeHTTP(w, r)
	})

	return gatewayMux
}
