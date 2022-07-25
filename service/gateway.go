package service

type gateway struct {
	routes map[string]string
}

func NewGateway() *gateway {
	return &gateway{
		routes: make(map[string]string),
	}
}

func (g *gateway) Register(route string, target string) {
	g.routes[route] = target
}
