package route

import (
	"os"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type StaticRoute struct {
	state *service.State
}

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
		if ctx.Request.URL.Path == "/" {
			// disable caching for index.html (/) to fix blank page issue
			ctx.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate,proxy-revalidate, max-age=0")
		}
		ctx.Next()
	})

	r.Static("/", s.state.GetWWWPath())

	return r
}
