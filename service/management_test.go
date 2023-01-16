package service

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/IceWhaleTech/CasaOS-Common/model"
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"gotest.tools/assert"
)

func init() {
	logger.LogInitConsoleOnly()
}

func TestRoutesPersistence(t *testing.T) {
	tmpdir1, _ := os.MkdirTemp("", "casaos-gateway-route-test")
	tmpdir2, _ := os.MkdirTemp("", "casaos-gateway-route-test")

	defer func() {
		os.RemoveAll(tmpdir1)
		os.RemoveAll(tmpdir2)
	}()

	state1 := NewState()
	state1.SetRuntimePath(tmpdir1)

	state2 := NewState()
	state2.SetRuntimePath(tmpdir2)

	management := NewManagementService(state1)

	route := &model.Route{
		Path:   "/test",
		Target: "http://localhost:8080",
	}

	if err := management.CreateRoute(route); err != nil {
		t.Fatal(err)
	}

	management = NewManagementService(state2)
	routes := management.GetRoutes()
	assert.Equal(t, 0, len(routes))

	management = NewManagementService(state1)
	routes = management.GetRoutes()
	assert.Equal(t, 1, len(routes))
	assert.Equal(t, "/test", routes[0].Path)
	assert.Equal(t, "http://localhost:8080", routes[0].Target)
}

func TestPathSorting(t *testing.T) {
	tmpdir, _ := os.MkdirTemp("", "casaos-gateway-route-test")

	defer func() {
		os.RemoveAll(tmpdir)
	}()

	state := NewState()
	state.SetRuntimePath(tmpdir)

	management := NewManagementService(state)

	routes := map[string]string{
		"/test":         "http://localhost:8080/",
		"/":             "http://localhost:8081/",
		"/testtesttest": "http://localhost:8082/",
		"/testtest":     "http://localhost:8083/",
	}

	for path, target := range routes {
		if err := management.CreateRoute(&model.Route{
			Path:   path,
			Target: target,
		}); err != nil {
			t.Fatal(err)
		}
	}

	for path, target := range routes {
		req := &http.Request{
			URL:    &url.URL{},
			Header: http.Header{},
		}

		proxy := management.GetProxy(path)
		assert.Assert(t, proxy != nil)

		proxy.Director(req)
		assert.Equal(t, target, req.URL.String())
	}
}
