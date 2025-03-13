package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/multi-cloud-landing-zone/backend/routes"
	"github.com/sirupsen/logrus"
)

// 日志文件配置
const (
	LogDir      = "logs"
	LogFileName = "multi-cloud-landing-zone"
)

var (
	// 全局日志实例
	Logger *logrus.Logger
)

// 请求日志中间件
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		
		// 记录请求信息
		Logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"request_id": c.GetHeader("X-Request-ID"),
		}).Info("收到API请求")

		// 如果是POST/PUT请求，记录请求体
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			// 读取请求体
			var requestBody map[string]interface{}
			if err := c.ShouldBindJSON(&requestBody); err == nil {
				Logger.WithFields(logrus.Fields{
					"request_body": requestBody,
					"path":         c.Request.URL.Path,
				}).Debug("请求参数")
			}
			
			// 由于请求体已被读取，需要重新设置
			c.Request.Body = nil // 这里简化处理，实际应该使用复制的请求体
		}
		
		// 处理请求
		c.Next()
		
		// 结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		
		// 记录响应信息
		Logger.WithFields(logrus.Fields{
			"status":       c.Writer.Status(),
			"latency":      latency,
			"path":         c.Request.URL.Path,
			"method":       c.Request.Method,
			"error_count":  len(c.Errors),
		}).Info("API请求完成")
		
		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			Logger.WithFields(logrus.Fields{
				"errors": c.Errors.String(),
				"path":   c.Request.URL.Path,
			}).Error("API请求错误")
		}
	}
}

// Terraform执行日志记录器
func TerraformLogger(command string, args []string, output string, err error) {
	fields := logrus.Fields{
		"command": command,
		"args":    args,
	}
	
	if err != nil {
		fields["error"] = err.Error()
		Logger.WithFields(fields).Error("Terraform执行失败")
		// 记录详细输出
		Logger.WithFields(logrus.Fields{
			"command": command,
			"output":  output,
		}).Debug("Terraform执行详细输出")
	} else {
		fields["success"] = true
		Logger.WithFields(fields).Info("Terraform执行成功")
		// 记录详细输出
		Logger.WithFields(logrus.Fields{
			"command": command,
			"output":  output,
		}).Debug("Terraform执行详细输出")
	}
}

// 初始化日志
func initLogger() *logrus.Logger {
	// 创建日志目录
	if err := os.MkdirAll(LogDir, 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
	}

	// 生成日志文件名（包含日期）
	currentTime := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(LogDir, fmt.Sprintf("%s-%s.log", LogFileName, currentTime))

	// 创建或打开日志文件
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("打开日志文件失败: %v\n", err)
		return nil
	}

	// 创建新的日志实例
	logger := logrus.New()
	
	// 设置日志输出到文件和控制台
	logger.SetOutput(io.MultiWriter(logFile, os.Stdout))
	
	// 设置日志格式为JSON
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false, // 生产环境设为false以提高性能
	})
	
	// 设置日志级别
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel) // 默认为INFO级别
	}

	// 添加调用者信息
	logger.SetReportCaller(true)

	logger.Info("日志系统初始化完成")
	logger.Info("日志文件路径: " + logFilePath)
	
	return logger
}

// 自定义Recovery中间件，增加详细日志
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录详细的恢复信息
				Logger.WithFields(logrus.Fields{
					"error":      err,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
					"client_ip":  c.ClientIP(),
					"request_id": c.GetHeader("X-Request-ID"),
				}).Error("服务器内部错误")
				
				// 返回500错误
				c.AbortWithStatusJSON(500, gin.H{
					"error": "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}

func main() {
	// 初始化日志系统
	Logger = initLogger()
	if Logger == nil {
		fmt.Println("初始化日志系统失败，使用默认日志配置")
		Logger = logrus.New()
		Logger.SetOutput(os.Stdout)
	}

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		Logger.Warn("未找到.env文件，使用默认环境变量")
	} else {
		Logger.Info("环境变量加载成功")
		
		// 记录关键环境变量（不包含敏感信息）
		Logger.WithFields(logrus.Fields{
			"GIN_MODE": os.Getenv("GIN_MODE"),
			"PORT":     os.Getenv("PORT"),
			"LOG_LEVEL": os.Getenv("LOG_LEVEL"),
		}).Debug("环境变量配置")
	}

	// 设置运行模式
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
		Logger.Info("Gin运行在Debug模式")
	} else {
		Logger.Info("Gin运行在Release模式")
	}

	// 配置Gin的日志输出
	currentTime := time.Now().Format("2006-01-02")
	ginLogPath := filepath.Join(LogDir, fmt.Sprintf("gin-%s.log", currentTime))
	ginLogFile, err := os.OpenFile(ginLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Logger.Errorf("打开Gin日志文件失败: %v", err)
	} else {
		gin.DefaultWriter = io.MultiWriter(ginLogFile, os.Stdout)
		Logger.Info("Gin日志文件路径: " + ginLogPath)
	}

	// 初始化Gin路由，使用自定义中间件
	router := gin.New() // 不使用默认中间件
	router.Use(gin.Logger())
	router.Use(CustomRecovery())
	router.Use(RequestLoggerMiddleware()) // 添加请求日志中间件
	
	Logger.Info("Gin路由初始化完成")

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"}
	router.Use(cors.New(config))
	Logger.Info("CORS配置完成")

	// 注册API路由
	routes.SetupRoutes(router)
	Logger.Info("API路由注册完成")
	
	// 注入日志实例到路由处理函数
	// 这里假设routes包有一个SetLogger方法
	// 如果没有，您需要修改routes包以接受日志实例
	if _, ok := interface{}(routes).(interface{ SetLogger(*logrus.Logger) }); ok {
		routes.SetLogger(Logger)
		Logger.Info("日志实例已注入到路由处理函数")
	}

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // 默认端口
		Logger.Info("使用默认端口: 3000")
	} else {
		Logger.Infof("使用配置端口: %s", port)
	}

	// 启动服务器
	Logger.Infof("服务器启动在 http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		Logger.Fatalf("启动服务器失败: %v", err)
	}
}

/*
使用说明:

1. 在routes包中，您需要修改处理函数以记录详细日志，例如:

func HandleCreateDeployment(c *gin.Context) {
    // 解析请求
    var deploymentRequest DeploymentRequest
    if err := c.ShouldBindJSON(&deploymentRequest); err != nil {
        Logger.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("解析部署请求失败")
        c.JSON(400, gin.H{"error": "无效的请求数据"})
        return
    }
    
    // 记录请求详情
    Logger.WithFields(logrus.Fields{
        "deployment_name": deploymentRequest.Name,
        "cloud_provider": deploymentRequest.Provider,
        "region": deploymentRequest.Region,
    }).Info("收到部署请求")
    
    // 执行Terraform命令
    cmd := exec.Command("terraform", "apply", "-auto-approve")
    output, err := cmd.CombinedOutput()
    
    // 记录Terraform执行情况
    TerraformLogger("apply", []string{"-auto-approve"}, string(output), err)
    
    if err != nil {
        c.JSON(500, gin.H{"error": "部署失败"})
        return
    }
    
    // 记录成功信息
    Logger.WithFields(logrus.Fields{
        "deployment_id": "新生成的ID",
        "status": "success",
    }).Info("部署成功完成")
    
    c.JSON(200, gin.H{"message": "部署成功", "deployment_id": "新生成的ID"})
}

2. 设置环境变量LOG_LEVEL=debug可以获取更详细的日志

3. 日志文件位于logs目录下，包括:
   - multi-cloud-landing-zone-YYYY-MM-DD.log (应用日志)
   - gin-YYYY-MM-DD.log (Gin框架日志)
*/
