package route

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"

	"github.com/IceWhaleTech/CasaOS-Common/model"
	"github.com/IceWhaleTech/CasaOS-Common/utils/jwt"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
)

var (
	_router *gin.Engine
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
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/v1/gateway/routes", bytes.NewReader(body))
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

	if err != nil {
		t.Error(err)
	}

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
	if err != nil {
		t.Error(err)
	}

	token, err := jwt.GenerateToken("test", "test", 0, "casaos", 10*time.Second)
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodPut, "/v1/gateway/port", bytes.NewReader(body))
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedPort, actualPort)

	// get
	req, _ = http.NewRequest(http.MethodGet, "/v1/gateway/port", nil)
	req.Header.Set("Authorization", token)

	w = httptest.NewRecorder()
	_router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result *model.Result
	decoder := json.NewDecoder(w.Body)

	err = decoder.Decode(&result)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expectedPort, result.Data)
}
