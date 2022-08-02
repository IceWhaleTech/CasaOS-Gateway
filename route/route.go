package route

import (
	"net/http"
	"os"

	"github.com/IceWhaleTech/CasaOS-Common/model"
	"github.com/IceWhaleTech/CasaOS-Common/utils/err"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var _management *service.Management

func NewRoutes(management *service.Management) *gin.Engine {
	_management = management

	// check if environment variable is set
	if ginMode, success := os.LookupEnv("GIN_MODE"); success {
		gin.SetMode(ginMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
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
	v1RoutesGroup := v1Group.Group("/gateway")

	v1RoutesGroup.Use()
	{
		v1RoutesGroup.GET("/routes", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, _management.GetRoutes())
		})

		v1RoutesGroup.POST("/routes", func(ctx *gin.Context) {
			var route *common.Route
			err := ctx.ShouldBindJSON(&route)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}

			_management.CreateRoute(route)

			ctx.Status(http.StatusCreated)
		})

		v1RoutesGroup.GET("/port", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, &model.Result{
				Success: err.SUCCESS,
				Message: err.GetMsg(err.SUCCESS),
				Data:    _management.GetGatewayPort(),
			})
		})
	}
}
