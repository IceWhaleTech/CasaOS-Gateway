// The commmon package provides structs and functions for external code to interact with this gateway service.
package common

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

const MANAGEMENT_URL_FILENAME = "management.url"
const GATEWAY_URL_FILENAME = "gateway.url"

type Route struct {
	Path   string `json:"path" binding:"required"`
	Target string `json:"target" binding:"required"`
}

type ManagementService interface {
}

type managementService struct {
	address string
}

func (m *managementService) CreateRoute(route *Route) {

}

func NewManagementService(runtimeVariablesPath string) (ManagementService, error) {

	managementAddresFile := filepath.Join(runtimeVariablesPath, MANAGEMENT_URL_FILENAME)

	buf, err := ioutil.ReadFile(managementAddresFile)
	if err != nil {
		return nil, err
	}

	address := string(buf)

	response, err := http.Get(address + "/ping")
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("failed to ping management service")
	}

	return &managementService{
		address: address,
	}, nil
}
