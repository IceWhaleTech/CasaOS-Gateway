package route

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

type StaticRoute struct {
	state *service.State
}

var startTime = time.Now()

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

	// serve /index.html to sovle 304 cache problem by 'If-Modified-Since: Wed, 21 Oct 2015 07:28:00 GMT' from web browser
	e.GET("/", func(ctx echo.Context) error {
		f, err := os.Open(filepath.Join(s.state.GetWWWPath(), "index.html"))
		if err != nil {
			return err
		}
		defer f.Close()
		http.ServeContent(ctx.Response(), ctx.Request(), "index.html", startTime, f)
		return nil
	})

	e.Static("/", s.state.GetWWWPath())

	return e
}
