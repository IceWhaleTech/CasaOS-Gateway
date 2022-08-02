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

	state1 := NewState()
	if err := state1.SetRuntimePath(tmpdir1); err != nil {
		t.Fatal(err)
	}

	state2 := NewState()
	if err := state2.SetRuntimePath(tmpdir2); err != nil {
		t.Fatal(err)
	}

	management := NewManagementService(state1)

	route := &common.Route{
		Path:   "/test",
		Target: "http://localhost:8080",
	}

	management.CreateRoute(route)

	management = NewManagementService(state2)
	routes := management.GetRoutes()
	assert.Equal(t, 0, len(routes))

	management = NewManagementService(state1)
	routes = management.GetRoutes()
	assert.Equal(t, 1, len(routes))
	assert.Equal(t, "/test", routes[0].Path)
	assert.Equal(t, "http://localhost:8080", routes[0].Target)
}
