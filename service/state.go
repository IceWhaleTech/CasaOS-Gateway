package service

type State struct {
	gatewayPort string
	runtimePath string
	onChange    []func(*State) error
}

func NewState() *State {
	return &State{
		gatewayPort: "",
		runtimePath: "",
		onChange:    make([]func(*State) error, 0),
	}
}

func (c *State) SetGatewayPort(port string) error {
	c.gatewayPort = port
	return c.change()
}

func (c *State) GetGatewayPort() string {
	return c.gatewayPort
}

func (c *State) SetRuntimePath(path string) error {
	c.runtimePath = path
	return c.change()
}

func (c *State) GetRuntimePath() string {
	return c.runtimePath
}

func (c *State) OnChange(f func(*State) error) {
	c.onChange = append(c.onChange, f)
}

func (c *State) change() error {
	for _, f := range c.onChange {
		if err := f(c); err != nil {
			return err
		}
	}

	return nil
}
