package route

import (
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var gateway = service.NewGateway()

func BuildManagementRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	buildV1Group(r)

	return r
}

func buildV1Group(r *gin.Engine) {
	v1Group := r.Group("/v1")
	buildV1GatewayGroup(v1Group)
}

func buildV1GatewayGroup(v1Group *gin.RouterGroup) {
	v1GatewayGroup := v1Group.Group("/gateway")

	v1GatewayGroup.POST("/routes", func(ctx *gin.Context) {
		json := make(map[string]interface{})
		ctx.BindJSON(&json)

		route := json["route"].(string)
		target := json["target"].(string)

		gateway.Register(route, target)

		ctx.JSON(200, gin.H{})
	})
}
