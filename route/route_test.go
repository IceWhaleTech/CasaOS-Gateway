package route

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/go-playground/assert/v2"
)

func TestPing(t *testing.T) {
	router := Build()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRoute(t *testing.T) {
	router := Build()

	w := httptest.NewRecorder()

	createRouteRequest := &common.CreateRouteRequest{
		Route:  "test",
		Target: "http://localhost:8080",
	}

	body, err := json.Marshal(createRouteRequest)

	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/v1/routes", bytes.NewReader(body))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
