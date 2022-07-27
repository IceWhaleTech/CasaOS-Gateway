package route

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var _management *service.Management

func Build(management *service.Management) *gin.Engine {
	_management = management

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	buildV1Group(r)

	return r
}

func buildV1Group(r *gin.Engine) {
	v1Group := r.Group("/v1")

	v1Group.Use()
	{
		buildV1RouteGroup(v1Group)
	}
}

func buildV1RouteGroup(v1Group *gin.RouterGroup) {
	v1RoutesGroup := v1Group.Group("/routes")

	v1RoutesGroup.Use()
	{
		v1RoutesGroup.GET("", func(ctx *gin.Context) {
			ctx.JSON(200, _management.GetRoutes())
		})

		v1RoutesGroup.POST("", func(ctx *gin.Context) {
			var route common.Route
			err := ctx.ShouldBindJSON(&route)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}

			_management.CreateRoute(route.Path, route.Target)

			ctx.JSON(200, route)
		})
	}
}
