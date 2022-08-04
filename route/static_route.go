package route

import (
	"net/http"
	"os"
	"text/template"

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

	r.Static("/ui", s.state.GetWWWPath())
	r.GET("/", func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		index, err := template.ParseFiles(s.state.GetWWWPath() + "/index.html")
		if err != nil {
			ctx.Status(http.StatusNoContent)
			return
		}

		if err := index.Execute(ctx.Writer, nil); err != nil {
			ctx.Status(http.StatusInternalServerError)
		}
	})

	return r
}
