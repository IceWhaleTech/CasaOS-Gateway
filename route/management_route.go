package route

import (
	"net/http"
	"os"

	"github.com/IceWhaleTech/CasaOS-Common/model"
	"github.com/IceWhaleTech/CasaOS-Common/utils/common_err"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type ManagementRoute struct {
	management *service.Management
}

func NewManagementRoute(management *service.Management) *ManagementRoute {
	return &ManagementRoute{
		management: management,
	}
}

func (m *ManagementRoute) GetRoutes() *gin.Engine {
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

	m.buildV1Group(r)

	return r
}

func (m *ManagementRoute) buildV1Group(r *gin.Engine) {
	v1Group := r.Group("/v1")

	v1Group.Use()
	{
		m.buildV1RouteGroup(v1Group)
	}
}

func (m *ManagementRoute) buildV1RouteGroup(v1Group *gin.RouterGroup) {
	v1RoutesGroup := v1Group.Group("/gateway")

	v1RoutesGroup.Use()
	{
		v1RoutesGroup.GET("/routes", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, m.management.GetRoutes())
		})

		v1RoutesGroup.POST("/routes", func(ctx *gin.Context) {
			var route *common.Route
			err := ctx.ShouldBindJSON(&route)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, model.Result{
					Success: common_err.CLIENT_ERROR,
					Message: err.Error(),
				})
				return
			}

			if err := m.management.CreateRoute(route); err != nil {
				ctx.JSON(http.StatusInternalServerError, model.Result{
					Success: common_err.SERVICE_ERROR,
					Message: err.Error(),
				})
				return
			}

			ctx.Status(http.StatusCreated)
		})

		v1RoutesGroup.GET("/port", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, model.Result{
				Success: common_err.SUCCESS,
				Message: common_err.GetMsg(common_err.SUCCESS),
				Data:    m.management.GetGatewayPort(),
			})
		})

		v1RoutesGroup.PUT("/port", func(ctx *gin.Context) {
			var request *common.ChangePortRequest

			if err := ctx.ShouldBindJSON(&request); err != nil {
				ctx.JSON(http.StatusBadRequest, model.Result{
					Success: common_err.CLIENT_ERROR,
					Message: err.Error(),
				})
				return
			}

			if err := m.management.SetGatewayPort(request.Port); err != nil {
				ctx.JSON(http.StatusInternalServerError, model.Result{
					Success: common_err.SERVICE_ERROR,
					Message: err.Error(),
				})
				return
			}
		})
	}
}
