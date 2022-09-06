// The commmon package provides structs and functions for external code to interact with this gateway service.
package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
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

	retry := 10

	for retry > 0 {
		if _, err := os.Stat(managementAddressFile); err == nil {
			break
		}

		fmt.Printf("gateway management address file `%s` not found, retrying in 1 second...(%d)\n", managementAddressFile, retry)

		time.Sleep(1 * time.Second)

		retry--
	}

	buf, err := os.ReadFile(managementAddressFile)
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
