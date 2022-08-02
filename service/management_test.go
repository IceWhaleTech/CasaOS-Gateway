package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS-Gateway/config"
	"gotest.tools/assert"
)

func TestRoutesPersistence(t *testing.T) {
	tmpdir1, _ := ioutil.TempDir("", "casaos-gateway-route-test")
	tmpdir2, _ := ioutil.TempDir("", "casaos-gateway-route-test")

	defer func() {
		os.RemoveAll(tmpdir1)
		os.RemoveAll(tmpdir2)
	}()

	cfg1 := config.NewConfig()
	if err := cfg1.SetRuntimePath(tmpdir1); err != nil {
		t.Fatal(err)
	}

	cfg2 := config.NewConfig()
	if err := cfg2.SetRuntimePath(tmpdir2); err != nil {
		t.Fatal(err)
	}

	management := NewManagementService(cfg1)

	route := &common.Route{
		Path:   "/test",
		Target: "http://localhost:8080",
	}

	management.CreateRoute(route)

	management = NewManagementService(cfg2)
	routes := management.GetRoutes()
	assert.Equal(t, 0, len(routes))

	management = NewManagementService(cfg1)
	routes = management.GetRoutes()
	assert.Equal(t, 1, len(routes))
	assert.Equal(t, "/test", routes[0].Path)
	assert.Equal(t, "http://localhost:8080", routes[0].Target)
}
