package route

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

type StaticRoute struct {
	state *service.State
}

var RouteCache = make(map[string]string)

func NewStaticRoute(state *service.State) *StaticRoute {
	return &StaticRoute{
		state: state,
	}
}

func (s *StaticRoute) GetRoute() http.Handler {
	e := echo.New()

	e.Use(echo_middleware.Gzip())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if _, ok := RouteCache[ctx.Request().URL.Path]; !ok {
				ctx.Response().Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate,proxy-revalidate, max-age=0")
				RouteCache[ctx.Request().URL.Path] = ctx.Request().URL.Path
			}
			return next(ctx)
		}
	})

	e.Static("/", s.state.GetWWWPath())

	return e
}
