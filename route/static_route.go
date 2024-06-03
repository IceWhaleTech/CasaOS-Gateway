package route

import (
	"os"
	"path"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
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

func (s *StaticRoute) GetRoute() *gin.Engine {
	// check if environment variable is set
	if ginMode, success := os.LookupEnv("GIN_MODE"); success {
		gin.SetMode(ginMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(func(ctx *gin.Context) {
		// Extract the file name from the path
		_, file := path.Split(ctx.Request.URL.Path)
		// If the file name contains a dot, it's likely a file
		if path.Ext(file) == "" {
			if _, ok := RouteCache[ctx.Request.URL.Path]; !ok {
				ctx.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate,proxy-revalidate, max-age=0")
				RouteCache[ctx.Request.URL.Path] = ctx.Request.URL.Path
			}
		}
		ctx.Next()
	})

	r.Static("/", s.state.GetWWWPath())

	return r
}
