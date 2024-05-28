package route

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"

	"github.com/IceWhaleTech/CasaOS-Common/model"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"gotest.tools/v3/assert"
)

var (
	_router http.Handler
	_state  *service.State
)

func init() {
	logger.LogInitConsoleOnly()
}

func setup(t *testing.T) func(t *testing.T) {
	tmpdir, _ := os.MkdirTemp("", "casaos-gateway-route-test")

	_state = service.NewState()
	if err := _state.SetRuntimePath(tmpdir); err != nil {
		t.Fatal(err)
	}

	management := service.NewManagementService(_state)
	managementRoute := NewManagementRoute(management)
	_router = managementRoute.GetRoute()

	return func(t *testing.T) {
		management = nil
		_router = nil
		os.RemoveAll(tmpdir)
	}
}

func TestPing(t *testing.T) {
	defer setup(t)(t)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	_router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRoute(t *testing.T) {
	defer setup(t)(t)

	route := &model.Route{
		Path:   "test",
		Target: "http://localhost:8080",
	}

	body, err := json.Marshal(route)
	assert.NilError(t, err)

	req, _ := http.NewRequest(http.MethodPost, "/v1/gateway/routes", bytes.NewReader(body))
	req.RemoteAddr = "127.0.0.1:0"

	w := httptest.NewRecorder()
	_router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/v1/gateway/routes", nil)
	w = httptest.NewRecorder()
	_router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var routes []*model.Route

	decoder := json.NewDecoder(w.Body)

	err = decoder.Decode(&routes)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(routes))
	assert.Equal(t, route.Path, routes[0].Path)
	assert.Equal(t, route.Target, routes[0].Target)
}

func TestChangePort(t *testing.T) {
	defer setup(t)(t)

	actualPort := ""

	_state.OnGatewayPortChange(func(s string) error {
		actualPort = s
		return nil
	})

	expectedPort := "123"

	// set
	request := &model.ChangePortRequest{
		Port: expectedPort,
	}

	body, err := json.Marshal(request)
	assert.NilError(t, err)

	req, _ := http.NewRequest(http.MethodPut, "/v1/gateway/port", bytes.NewReader(body))
	req.RemoteAddr = "127.0.0.1:0"

	w := httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedPort, actualPort)

	// get
	req, _ = http.NewRequest(http.MethodGet, "/v1/gateway/port", nil)

	w = httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result *model.Result
	decoder := json.NewDecoder(w.Body)

	err = decoder.Decode(&result)
	assert.NilError(t, err)
	assert.Equal(t, expectedPort, result.Data)
}

func TestChangePortNegative(t *testing.T) {
	defer setup(t)(t)

	expectedPort := "123"

	// set
	request := &model.ChangePortRequest{
		Port: expectedPort,
	}

	body, err := json.Marshal(request)
	assert.NilError(t, err)

	req, _ := http.NewRequest(http.MethodPut, "/v1/gateway/port", bytes.NewReader(body))
	req.RemoteAddr = "127.0.0.1:0"

	w := httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedPort, "123")

	// get
	req, _ = http.NewRequest(http.MethodGet, "/v1/gateway/port", nil)

	w = httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result *model.Result
	decoder := json.NewDecoder(w.Body)

	err = decoder.Decode(&result)
	assert.NilError(t, err)
	assert.Equal(t, expectedPort, result.Data)

	// emulate error
	_state.OnGatewayPortChange(func(_ string) error {
		return errors.New("error")
	})

	// set
	request.Port = "456"

	body, err = json.Marshal(request)
	assert.NilError(t, err)

	req, _ = http.NewRequest(http.MethodPut, "/v1/gateway/port", bytes.NewReader(body))
	req.RemoteAddr = "127.0.0.1:0"

	w = httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedPort, "123")

	// get
	req, _ = http.NewRequest(http.MethodGet, "/v1/gateway/port", nil)

	w = httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	decoder = json.NewDecoder(w.Body)

	err = decoder.Decode(&result)
	assert.NilError(t, err)
	assert.Equal(t, expectedPort, result.Data)
}
