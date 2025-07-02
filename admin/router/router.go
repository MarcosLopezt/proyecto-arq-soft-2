package router

import (
	"admin/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
    router.GET("/admin/services", controllers.ListarInstancias)
    router.POST("/admin/services/create", controllers.CrearInstancia)
    router.POST("/admin/services/remove", controllers.EliminarInstancia)
    router.GET("/admin/services/balanceo", controllers.ProbarBalanceoCarga)
}
