// The commmon package provides structs and functions for external code to interact with this gateway service.
package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
)

const MANAGEMENT_URL_FILENAME = "management.url"
const GATEWAY_URL_FILENAME = "gateway.url"

const API_PATH = "/v1/routes"

type Route struct {
	Path   string `json:"path" binding:"required"`
	Target string `json:"target" binding:"required"`
}

type ManagementService interface {
}

type managementService struct {
	address string
}

func (m *managementService) CreateRoute(route *Route) error {
	url := path.Join(m.address, API_PATH)
	body, err := json.Marshal(route)
	if err != nil {
		return err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return errors.New("failed to create route (status code: " + fmt.Sprint(response.StatusCode) + ")")
	}

	return nil
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
