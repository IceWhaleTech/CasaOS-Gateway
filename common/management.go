// The commmon package provides structs and functions for external code to interact with this gateway service.
package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	ManagementURLFilename = "management.url"
	StaticURLFilename     = "static.url"
	APIGatewayRoutes      = "/v1/gateway/routes"
	APIGatewayPort        = "/v1/gateway/port"
)

type ManagementService interface {
	CreateRoute(route *Route) error
	ChangePort(request *ChangePortRequest) error
}

type managementService struct {
	address string
}

func (m *managementService) CreateRoute(route *Route) error {
	url := strings.TrimSuffix(m.address, "/") + "/" + strings.TrimPrefix(APIGatewayRoutes, "/")
	body, err := json.Marshal(route)
	if err != nil {
		return err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body)) //nolint:gosec
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return errors.New("failed to create route (status code: " + fmt.Sprint(response.StatusCode) + ")")
	}

	return nil
}

func (m *managementService) ChangePort(request *ChangePortRequest) error {
	url := strings.TrimSuffix(m.address, "/") + "/" + strings.TrimPrefix(APIGatewayPort, "/")
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body)) //nolint:gosec
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("failed to change port (status code: " + fmt.Sprint(response.StatusCode) + ")")
	}

	return nil
}

func NewManagementService(RuntimePath string) (ManagementService, error) {
	managementAddressFile := filepath.Join(RuntimePath, ManagementURLFilename)

	buf, err := ioutil.ReadFile(managementAddressFile)
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
