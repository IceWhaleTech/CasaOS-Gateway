package service

type management struct {
	routes map[string]string
}

func NewManagementService() *management {
	return &management{
		routes: make(map[string]string),
	}
}

func (g *management) CreateRoute(route string, target string) {
	g.routes[route] = target
}

func (g *management) GetRoutes() map[string]string {
	return g.routes
}
