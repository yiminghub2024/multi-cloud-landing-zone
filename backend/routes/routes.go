package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/multi-cloud-landing-zone/backend/controllers"
)

// SetupRoutes 配置API路由
func SetupRoutes(router *gin.Engine) {
	// API路由组
	api := router.Group("/api")
	{
		// 获取云服务提供商列表
		api.GET("/providers", controllers.GetProviders)

		// 获取区域列表
		api.GET("/regions/:provider", controllers.GetRegions)

		// 获取可用区列表
		api.GET("/azs/:provider/:region", controllers.GetAvailabilityZones)

		// 获取云组件列表
		api.GET("/components/:provider/:region", controllers.GetComponents)

		// 执行部署
		api.POST("/deploy", controllers.StartDeployment)

		// 获取部署状态
		api.GET("/deployment/status", controllers.GetDeploymentStatus)
	}
}
