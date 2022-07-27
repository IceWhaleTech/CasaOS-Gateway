package service

type Management struct {
	routes map[string]string
}

func NewManagementService() *Management {
	return &Management{
		routes: make(map[string]string),
	}
}

func (g *Management) CreateRoute(route string, target string) {
	g.routes[route] = target
}

func (g *Management) GetRoutes() map[string]string {
	return g.routes
}

func (g *Management) GetRoute(route string) string {
	return g.routes[route]
}
