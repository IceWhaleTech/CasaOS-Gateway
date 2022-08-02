package route

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var router *gin.Engine

func setup(t *testing.T) func(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "casaos-gateway-route-test")

	state := service.NewState()
	if err := state.SetRuntimePath(tmpdir); err != nil {
		t.Fatal(err)
	}

	management := service.NewManagementService(state)
	router = NewRoutes(management)

	return func(t *testing.T) {
		management = nil
		router = nil
		os.RemoveAll(tmpdir)
	}
}

func TestPing(t *testing.T) {
	defer setup(t)(t)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRoute(t *testing.T) {
	defer setup(t)(t)

	route := &common.Route{
		Path:   "test",
		Target: "http://localhost:8080",
	}

	body, err := json.Marshal(route)
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/v1/gateway/routes", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/v1/gateway/routes", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var routes []*common.Route

	decoder := json.NewDecoder(w.Body)

	err = decoder.Decode(&routes)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(routes))
	assert.Equal(t, route.Path, routes[0].Path)
	assert.Equal(t, route.Target, routes[0].Target)
}
