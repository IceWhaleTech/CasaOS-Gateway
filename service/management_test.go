package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"gotest.tools/assert"
)

func TestRoutesPersistence(t *testing.T) {
	tmpdir1, _ := ioutil.TempDir("", "casaos-gateway-route-test")
	tmpdir2, _ := ioutil.TempDir("", "casaos-gateway-route-test")

	defer func() {
		os.RemoveAll(tmpdir1)
		os.RemoveAll(tmpdir2)
	}()

	management := NewManagementService(tmpdir1)

	route := &common.Route{
		Path:   "/test",
		Target: "http://localhost:8080",
	}

	management.CreateRoute(route)

	management = NewManagementService(tmpdir2)
	routes := management.GetRoutes()
	assert.Equal(t, 0, len(routes))

	management = NewManagementService(tmpdir1)
	routes = management.GetRoutes()
	assert.Equal(t, 1, len(routes))
	assert.Equal(t, "/test", routes[0].Path)
	assert.Equal(t, "http://localhost:8080", routes[0].Target)
}
