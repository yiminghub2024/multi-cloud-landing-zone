package handlers

import (
	"github.com/gin-gonic/gin"
)

// ProviderHandler 处理云服务提供商相关的请求
type ProviderHandler interface {
	GetProviders(c *gin.Context)
	GetRegions(c *gin.Context)
	GetAvailabilityZones(c *gin.Context)
	GetComponents(c *gin.Context)
}

// DeploymentHandler 处理部署相关的请求
type DeploymentHandler interface {
	StartDeployment(c *gin.Context)
	GetDeploymentStatus(c *gin.Context)
}
