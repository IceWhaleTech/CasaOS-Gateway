package service

type State struct {
	gatewayPort         string
	onGatewayPortChange []func(string) error

	runtimePath string
	wwwPath     string
}

func NewState() *State {
	return &State{
		gatewayPort:         "",
		onGatewayPortChange: make([]func(string) error, 0),

		runtimePath: "",
		wwwPath:     "",
	}
}

func (c *State) SetGatewayPort(port string) error {
	c.gatewayPort = port
	return c.notifiyOnGatewayPortChange()
}

func (c *State) GetGatewayPort() string {
	return c.gatewayPort
}

func (c *State) OnGatewayPortChange(f func(string) error) {
	c.onGatewayPortChange = append(c.onGatewayPortChange, f)
}

func (c *State) notifiyOnGatewayPortChange() error {
	for _, f := range c.onGatewayPortChange {
		if err := f(c.gatewayPort); err != nil {
			return err
		}
	}

	return nil
}

func (c *State) SetRuntimePath(path string) error {
	c.runtimePath = path
	return nil
}

func (c *State) GetRuntimePath() string {
	return c.runtimePath
}

func (c *State) SetWWWPath(path string) error {
	c.wwwPath = path
	return nil
}

func (c *State) GetWWWPath() string {
	return c.wwwPath
}
