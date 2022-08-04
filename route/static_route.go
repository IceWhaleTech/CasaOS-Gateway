package route

import (
	"os"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type StaticRoute *gin.Engine

func NewStaticRoutes(state *service.State) StaticRoute {
	// check if environment variable is set
	if ginMode, success := os.LookupEnv("GIN_MODE"); success {
		gin.SetMode(ginMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Static("/", state.GetWWWPath())
	r.Static("/UI", state.GetWWWPath())

	return r
}
