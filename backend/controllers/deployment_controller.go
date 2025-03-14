package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings" 
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/multi-cloud-landing-zone/backend/models"
	"github.com/multi-cloud-landing-zone/backend/utils"

)

// 全局部署状态对象
var (
	deploymentStatus = models.DeploymentStatus{
		Status:   "idle",
		Progress: 0,
		Message:  "",
		Logs:     []string{},
		Result:   nil,
		Topology: nil,
	}
	deploymentMutex sync.Mutex
)

// StartDeployment 开始部署过程
func StartDeployment(c *gin.Context) {
	utils.LogInfo("收到新的部署请求")
	
	// 读取请求体原始数据用于日志记录
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.LogError(fmt.Sprintf("读取请求体失败: %v", err))
		c.JSON(400, gin.H{
			"success": false,
			"message": "无法读取请求数据: " + err.Error(),
		})
		return
	}
	
	// 记录原始请求数据
	utils.LogInfo(fmt.Sprintf("收到前端传送的原始参数:\n%s", string(requestBody)))
	
	// 重新设置请求体，因为已经被读取
	c.Request.Body = ioutil.NopCloser(strings.NewReader(string(requestBody)))
	
	var deploymentConfig models.DeploymentConfig
	if err := c.ShouldBindJSON(&deploymentConfig); err != nil {
		utils.LogError(fmt.Sprintf("解析部署配置失败: %v", err))
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的部署配置: " + err.Error(),
		})
		return
	}

	// 记录解析后的部署配置
	configJSON, _ := json.MarshalIndent(deploymentConfig, "", "  ")
	utils.LogInfo(fmt.Sprintf("解析后的部署配置:\n%s", string(configJSON)))

	// 更新部署状态
	deploymentMutex.Lock()
	deploymentStatus.Status = "preparing"
	deploymentStatus.Progress = 0
	deploymentStatus.Message = "正在准备部署资源..."
	deploymentStatus.Logs = []string{"开始部署过程..."}
	deploymentStatus.Result = nil
	deploymentStatus.Topology = nil
	deploymentMutex.Unlock()

	// 异步处理部署
	deploymentID := fmt.Sprintf("%d", time.Now().Unix())
	utils.LogInfo(fmt.Sprintf("开始异步处理部署，部署ID: %s", deploymentID))
	go processDeploy(deploymentConfig, deploymentID)

	// 立即返回响应，不等待部署完成
	c.JSON(200, gin.H{
		"success":      true,
		"message":      "部署已开始",
		"deploymentId": deploymentID,
	})
	utils.LogInfo(fmt.Sprintf("已返回部署开始响应，部署ID: %s", deploymentID))
}

// GetDeploymentStatus 获取部署状态
func GetDeploymentStatus(c *gin.Context) {
	utils.LogInfo("收到获取部署状态请求")
	
	deploymentMutex.Lock()
	status := deploymentStatus
	deploymentMutex.Unlock()

	c.JSON(200, gin.H{
		"success": true,
		"data":    status,
	})
	
	utils.LogInfo(fmt.Sprintf("已返回部署状态: %s, 进度: %d%%", status.Status, status.Progress))
}

// processDeploy 异步处理部署过程
func processDeploy(config models.DeploymentConfig, deploymentID string) {
	utils.LogInfo(fmt.Sprintf("开始处理部署 ID: %s", deploymentID))
	
	defer func() {
		if r := recover(); r != nil {
			utils.LogError(fmt.Sprintf("部署过程崩溃: %v", r))
			deploymentMutex.Lock()
			deploymentStatus.Status = "failed"
			deploymentStatus.Message = fmt.Sprintf("部署过程崩溃: %v", r)
			deploymentStatus.Logs = append(deploymentStatus.Logs, fmt.Sprintf("错误: %v", r))
			deploymentMutex.Unlock()
		}
	}()

	try := func(action func() error) {
		if err := action(); err != nil {
			utils.LogError(fmt.Sprintf("部署操作失败: %v", err))
			panic(err)
		}
	}

	try(func() error {
		// 创建部署工作目录
		workDir := filepath.Join("terraform", "deployments", deploymentID)
		
		if err := os.MkdirAll(workDir, 0755); err != nil {
			utils.LogError(fmt.Sprintf("创建部署工作目录失败: %v", err))
			return fmt.Errorf("创建部署工作目录失败: %w", err)
		}

		utils.LogInfo(fmt.Sprintf("创建部署工作目录: %s", workDir))
		deploymentMutex.Lock()
		deploymentStatus.Logs = append(deploymentStatus.Logs, fmt.Sprintf("创建部署工作目录: %s", workDir))
		deploymentStatus.Progress = 10
		deploymentStatus.Message = "正在生成Terraform配置..."
		deploymentMutex.Unlock()

		// 生成Terraform配置文件
		terraformConfig := utils.GenerateTerraformConfig(config)
		mainTfPath := filepath.Join(workDir, "main.tf")
		
		// 保存Terraform配置文件
		if err := utils.SaveTerraformConfig(terraformConfig, mainTfPath); err != nil {
			utils.LogError(fmt.Sprintf("保存Terraform配置文件失败: %v", err))
			return fmt.Errorf("保存Terraform配置文件失败: %w", err)
		}

		utils.LogInfo(fmt.Sprintf("Terraform配置文件已保存到: %s", mainTfPath))
		deploymentMutex.Lock()
		deploymentStatus.Logs = append(deploymentStatus.Logs, "生成Terraform配置文件完成")
		deploymentStatus.Logs = append(deploymentStatus.Logs, fmt.Sprintf("Terraform配置文件路径: %s", mainTfPath))
		deploymentStatus.Progress = 20
		deploymentStatus.Message = "正在初始化Terraform..."
		deploymentMutex.Unlock()

		// 初始化Terraform
		utils.LogInfo("开始初始化Terraform")
		cmd := exec.Command("terraform", "init")
		cmd.Dir = workDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			utils.LogError(fmt.Sprintf("Terraform初始化失败: %v, 输出: %s", err, string(output)))
			return fmt.Errorf("Terraform初始化失败: %w, 输出: %s", err, string(output))
		}

		utils.LogInfo(fmt.Sprintf("Terraform初始化完成，输出:\n%s", string(output)))
		deploymentMutex.Lock()
		deploymentStatus.Logs = append(deploymentStatus.Logs, "Terraform初始化完成")
		deploymentStatus.Logs = append(deploymentStatus.Logs, string(output))
		deploymentStatus.Progress = 30
		deploymentStatus.Message = "正在验证Terraform配置..."
		deploymentMutex.Unlock()

		// 验证Terraform配置
		utils.LogInfo("开始验证Terraform配置")
		cmd = exec.Command("terraform", "validate")
		cmd.Dir = workDir
		output, err = cmd.CombinedOutput()
		if err != nil {
			utils.LogError(fmt.Sprintf("Terraform配置验证失败: %v, 输出: %s", err, string(output)))
			return fmt.Errorf("Terraform配置验证失败: %w, 输出: %s", err, string(output))
		}

		utils.LogInfo(fmt.Sprintf("Terraform配置验证通过，输出:\n%s", string(output)))
		deploymentMutex.Lock()
		deploymentStatus.Logs = append(deploymentStatus.Logs, "Terraform配置验证通过")
		deploymentStatus.Logs = append(deploymentStatus.Logs, string(output))
		deploymentStatus.Progress = 40
		deploymentStatus.Message = "正在生成Terraform执行计划..."
		deploymentMutex.Unlock()

		// 生成执行计划
		utils.LogInfo("开始生成Terraform执行计划")
		cmd = exec.Command("terraform", "plan", "-out=tfplan")
		cmd.Dir = workDir
		output, err = cmd.CombinedOutput()
		if err != nil {
			utils.LogError(fmt.Sprintf("Terraform计划生成失败: %v, 输出: %s", err, string(output)))
			return fmt.Errorf("Terraform计划生成失败: %w, 输出: %s", err, string(output))
		}

		utils.LogInfo(fmt.Sprintf("Terraform执行计划生成完成，输出:\n%s", string(output)))
		deploymentMutex.Lock()
		deploymentStatus.Logs = append(deploymentStatus.Logs, "Terraform执行计划生成完成")
		deploymentStatus.Logs = append(deploymentStatus.Logs, string(output))
		deploymentStatus.Progress = 60
		deploymentStatus.Message = "正在执行Terraform部署..."
		deploymentMutex.Unlock()

		// 执行部署
		utils.LogInfo("开始执行Terraform部署")
		cmd = exec.Command("terraform", "apply", "-auto-approve", "tfplan")
		cmd.Dir = workDir
		output, err = cmd.CombinedOutput()
		if err != nil {
			utils.LogError(fmt.Sprintf("Terraform部署失败: %v, 输出: %s", err, string(output)))
			return fmt.Errorf("Terraform部署失败: %w, 输出: %s", err, string(output))
		}

		utils.LogInfo(fmt.Sprintf("Terraform部署执行完成，输出:\n%s", string(output)))
		deploymentMutex.Lock()
		deploymentStatus.Logs = append(deploymentStatus.Logs, "Terraform部署执行完成")
		deploymentStatus.Logs = append(deploymentStatus.Logs, string(output))
		deploymentStatus.Progress = 90
		deploymentStatus.Message = "正在生成资源拓扑图..."
		deploymentMutex.Unlock()

		// 生成拓扑图
		utils.LogInfo("开始生成资源拓扑图")
		topology := utils.GenerateTopology(config)

		// 完成部署
		utils.LogInfo("部署完成，更新最终状态")
		deploymentMutex.Lock()
		deploymentStatus.Status = "completed"
		deploymentStatus.Progress = 100
		deploymentStatus.Message = "部署完成"
		deploymentStatus.Result = map[string]interface{}{
			"deploymentId":   deploymentID,
			"cloudProvider":  config.CloudProvider,
			"region":         config.Region,
			"az":             config.AZ,
			"vpc":            config.VPC,
			"subnet":         config.Subnet,
			"components":     config.Components,
			"terraformPath":  mainTfPath,
		}
		deploymentStatus.Topology = topology
		deploymentMutex.Unlock()

		utils.LogInfo(fmt.Sprintf("部署 ID: %s 已成功完成", deploymentID))
		return nil
	})
}
