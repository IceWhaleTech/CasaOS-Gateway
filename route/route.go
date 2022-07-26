package route

import (
	"encoding/json"
	"net/http"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var gateway = service.NewManagementService()

func Build() *gin.Engine {
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
			ctx.JSON(200, gateway.GetRoutes())
		})

		v1RoutesGroup.POST("", func(ctx *gin.Context) {
			decoder := json.NewDecoder(ctx.Request.Body)

			var request common.CreateRouteRequest
			err := decoder.Decode(&request)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			gateway.Register(request.Route, request.Target)

			ctx.JSON(200, gin.H{})
		})
	}
}
