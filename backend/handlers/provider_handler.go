package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/multi-cloud-landing-zone/backend/controllers"
)

// ControllerProviderHandler 使用controllers包实现ProviderHandler接口
type ControllerProviderHandler struct{}

// NewProviderHandler 创建一个新的ProviderHandler实例
func NewProviderHandler() ProviderHandler {
	return &ControllerProviderHandler{}
}

// GetProviders 获取云服务提供商列表
func (h *ControllerProviderHandler) GetProviders(c *gin.Context) {
	controllers.GetProviders(c)
}

// GetRegions 获取区域列表
func (h *ControllerProviderHandler) GetRegions(c *gin.Context) {
	controllers.GetRegions(c)
}

// GetAvailabilityZones 获取可用区列表
func (h *ControllerProviderHandler) GetAvailabilityZones(c *gin.Context) {
	controllers.GetAvailabilityZones(c)
}

// GetComponents 获取云组件列表
func (h *ControllerProviderHandler) GetComponents(c *gin.Context) {
	controllers.GetComponents(c)
}
