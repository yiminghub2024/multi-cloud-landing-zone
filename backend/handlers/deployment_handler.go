package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/multi-cloud-landing-zone/backend/controllers"
)

// ControllerDeploymentHandler 使用controllers包实现DeploymentHandler接口
type ControllerDeploymentHandler struct{}

// NewDeploymentHandler 创建一个新的DeploymentHandler实例
func NewDeploymentHandler() DeploymentHandler {
	return &ControllerDeploymentHandler{}
}

// StartDeployment 开始部署过程
func (h *ControllerDeploymentHandler) StartDeployment(c *gin.Context) {
	controllers.StartDeployment(c)
}

// GetDeploymentStatus 获取部署状态
func (h *ControllerDeploymentHandler) GetDeploymentStatus(c *gin.Context) {
	controllers.GetDeploymentStatus(c)
}
