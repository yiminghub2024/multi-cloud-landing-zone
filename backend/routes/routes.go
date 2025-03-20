package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/multi-cloud-landing-zone/backend/handlers"
)

// SetupRoutes 配置API路由
func SetupRoutes(router *gin.Engine) {
	// 创建处理器实例
	providerHandler := handlers.NewProviderHandler()
	deploymentHandler := handlers.NewDeploymentHandler()

	// API路由组
	api := router.Group("/api")
	{
		// 获取云服务提供商列表
		api.GET("/providers", providerHandler.GetProviders)

		// 获取区域列表
		api.GET("/regions/:provider", providerHandler.GetRegions)

		// 获取可用区列表
		api.GET("/azs/:provider/:region", providerHandler.GetAvailabilityZones)

		// 获取云组件列表
		api.GET("/components/:provider/:region", providerHandler.GetComponents)

		// 执行部署
		api.POST("/deploy", deploymentHandler.StartDeployment)

		// 获取部署状态
		api.GET("/deployment/status", deploymentHandler.GetDeploymentStatus)
	}
}

// SetLogger 设置日志记录器（如果需要）
func SetLogger() {
	// 这个函数可以用于设置路由包的日志记录器
	// 如果main.go中有调用routes.SetLogger()，则需要保留此函数
}
