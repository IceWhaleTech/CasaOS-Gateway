package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRewriteRequestSourceIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		xff        string
		xRealIP    string
		wantXFF    string
	}{
		{
			name:       "direct remote client strips spoofed headers",
			remoteAddr: "203.0.113.10:54321",
			xff:        "127.0.0.1",
			xRealIP:    "127.0.0.1",
			wantXFF:    "",
		},
		{
			name:       "ipv4 loopback reverse proxy preserves trimmed last hop",
			remoteAddr: "127.0.0.1:41234",
			xff:        "127.0.0.1, 198.51.100.77 ",
			xRealIP:    "127.0.0.1",
			wantXFF:    "198.51.100.77",
		},
		{
			name:       "ipv6 loopback reverse proxy with brackets preserves last hop",
			remoteAddr: "[::1]:41234",
			xff:        "::1, 192.0.2.45",
			xRealIP:    "::1",
			wantXFF:    "192.0.2.45",
		},
		{
			name:       "trailing comma in xff does not emit empty entry",
			remoteAddr: "127.0.0.1:41234",
			xff:        "198.51.100.12, ",
			wantXFF:    "198.51.100.12",
		},
		{
			name:       "remote addr without port does not panic",
			remoteAddr: "127.0.0.1",
			xff:        "198.51.100.99",
			wantXFF:    "198.51.100.99",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://casaos.local/", nil)
			req.RemoteAddr = tc.remoteAddr
			if tc.xff != "" {
				req.Header.Set("X-Forwarded-For", tc.xff)
			}
			if tc.xRealIP != "" {
				req.Header.Set("X-Real-IP", tc.xRealIP)
			}

			rewriteRequestSourceIP(req)

			assert.Equal(t, req.Header.Get("X-Forwarded-For"), tc.wantXFF)
			assert.Equal(t, req.Header.Get("X-Real-IP"), "")
		})
	}
}
