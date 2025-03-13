package main

import (
        "fmt"
        "io"
        "os"
        "path/filepath"
        "time"

        "github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
        "github.com/joho/godotenv"
        "github.com/multi-cloud-landing-zone/backend/routes"
        "github.com/multi-cloud-landing-zone/backend/utils"
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
                
                // 同时使用新的日志系统记录
                utils.LogError(fmt.Sprintf("Terraform执行失败: %v, 命令: %s %v", err, command, args))
                utils.LogDebug(fmt.Sprintf("Terraform执行详细输出:\n%s", output))
        } else {
                fields["success"] = true
                Logger.WithFields(fields).Info("Terraform执行成功")
                // 记录详细输出
                Logger.WithFields(logrus.Fields{
                        "command": command,
                        "output":  output,
                }).Debug("Terraform执行详细输出")
                
                // 同时使用新的日志系统记录
                utils.LogInfo(fmt.Sprintf("Terraform执行成功, 命令: %s %v", command, args))
                utils.LogDebug(fmt.Sprintf("Terraform执行详细输出:\n%s", output))
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
                                
                                // 使用新的日志系统记录
                                utils.LogError(fmt.Sprintf("服务器内部错误: %v, 路径: %s, 方法: %s", 
                                        err, c.Request.URL.Path, c.Request.Method))
                                
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
        
        // 初始化新的日志系统
        if err := utils.InitLogger(); err != nil {
                Logger.Errorf("初始化新日志系统失败: %v", err)
        } else {
                Logger.Info("新日志系统初始化成功")
                utils.LogInfo("新日志系统初始化成功，与原有日志系统集成")
        }

        // 加载环境变量
        if err := godotenv.Load(); err != nil {
                Logger.Warn("未找到.env文件，使用默认环境变量")
                utils.LogWarn("未找到.env文件，使用默认环境变量")
        } else {
                Logger.Info("环境变量加载成功")
                utils.LogInfo("环境变量加载成功")
                
                // 记录关键环境变量（不包含敏感信息）
                Logger.WithFields(logrus.Fields{
                        "GIN_MODE": os.Getenv("GIN_MODE"),
                        "PORT":     os.Getenv("PORT"),
                        "LOG_LEVEL": os.Getenv("LOG_LEVEL"),
                }).Debug("环境变量配置")
                
                utils.LogDebug(fmt.Sprintf("环境变量配置: GIN_MODE=%s, PORT=%s, LOG_LEVEL=%s",
                        os.Getenv("GIN_MODE"), os.Getenv("PORT"), os.Getenv("LOG_LEVEL")))
        }

        // 设置运行模式
        gin.SetMode(gin.ReleaseMode)
        if os.Getenv("GIN_MODE") == "debug" {
                gin.SetMode(gin.DebugMode)
                Logger.Info("Gin运行在Debug模式")
                utils.LogInfo("Gin运行在Debug模式")
        } else {
                Logger.Info("Gin运行在Release模式")
                utils.LogInfo("Gin运行在Release模式")
        }

        // 配置Gin的日志输出
        currentTime := time.Now().Format("2006-01-02")
        ginLogPath := filepath.Join(LogDir, fmt.Sprintf("gin-%s.log", currentTime))
        ginLogFile, err := os.OpenFile(ginLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
                Logger.Errorf("打开Gin日志文件失败: %v", err)
                utils.LogError(fmt.Sprintf("打开Gin日志文件失败: %v", err))
        } else {
                gin.DefaultWriter = io.MultiWriter(ginLogFile, os.Stdout)
                Logger.Info("Gin日志文件路径: " + ginLogPath)
                utils.LogInfo("Gin日志文件路径: " + ginLogPath)
        }

        // 初始化Gin路由，使用自定义中间件
        router := gin.New() // 不使用默认中间件
        router.Use(gin.Logger())
        router.Use(CustomRecovery())
        router.Use(RequestLoggerMiddleware()) // 添加请求日志中间件
        	router.Use(cors.New(cors.Config{
					AllowOrigins:     []string{"*"}, // 允许所有来源
							AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
									AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
											ExposeHeaders:    []string{"Content-Length"},
													AllowCredentials: true,
															MaxAge:           12 * time.Hour,
																}))
        Logger.Info("Gin路由初始化完成")
        utils.LogInfo("Gin路由初始化完成")

        // 配置CORS - 使用更宽松的配置并添加调试日志
        router.Use(func(c *gin.Context) {
            origin := c.Request.Header.Get("Origin")
            
            // 添加调试日志
            Logger.WithFields(logrus.Fields{
                "request_method": c.Request.Method,
                "request_path": c.Request.URL.Path,
                "origin": origin,
                "request_headers": c.Request.Header,
            }).Info("收到跨域请求")
            
            // 设置CORS头部
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
            c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
            c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID")
            c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24小时内不再发送预检请求
            
            // 处理OPTIONS预检请求
            if c.Request.Method == "OPTIONS" {
                Logger.WithFields(logrus.Fields{
                    "cors_headers": c.Writer.Header(),
                }).Info("响应OPTIONS预检请求")
                
                c.AbortWithStatus(204)
                return
            }
            
            c.Next()
            
            // 记录响应状态
            Logger.WithFields(logrus.Fields{
                "status": c.Writer.Status(),
                "response_headers": c.Writer.Header(),
            }).Info("CORS请求处理完成")
        })
        
        Logger.Info("CORS配置完成 - 使用增强的自定义中间件，允许所有来源并支持凭证")
        utils.LogInfo("CORS配置完成 - 使用增强的自定义中间件，允许所有来源并支持凭证")

        // 注册API路由
        routes.SetupRoutes(router)
        Logger.Info("API路由注册完成")
        utils.LogInfo("API路由注册完成")

        // 获取端口
        port := os.Getenv("PORT")
        if port == "" {
                port = "3000" // 默认端口
                Logger.Info("使用默认端口: 3000")
                utils.LogInfo("使用默认端口: 3000")
        } else {
                Logger.Infof("使用配置端口: %s", port)
                utils.LogInfo(fmt.Sprintf("使用配置端口: %s", port))
        }

        // 启动服务器
        Logger.Infof("服务器启动在 http://0.0.0.0:%s (可通过 http://10.168.0.5:%s 访问)", port, port)
        utils.LogInfo(fmt.Sprintf("服务器启动在 http://0.0.0.0:%s (可通过 http://10.168.0.5:%s 访问)", port, port))
        if err := router.Run(":" + port); err != nil {
                Logger.Fatalf("启动服务器失败: %v", err)
                utils.LogError(fmt.Sprintf("启动服务器失败: %v", err))
        }
}

/*
使用说明:

1. 在routes包中，您需要修改处理函数以记录详细日志，例如:

func HandleCreateDeployment(c *gin.Context) {
    // 导入全局日志实例
    logger := main.Logger // 您需要导入main包
    
    // 或者使用新的日志系统
    utils.LogInfo("收到部署请求")
    
    // 解析请求
    var deploymentRequest DeploymentRequest
    if err := c.ShouldBindJSON(&deploymentRequest); err != nil {
        logger.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("解析部署请求失败")
        utils.LogError(fmt.Sprintf("解析部署请求失败: %v", err))
        c.JSON(400, gin.H{"error": "无效的请求数据"})
        return
    }
    
    // 记录请求详情
    logger.WithFields(logrus.Fields{
        "deployment_name": deploymentRequest.Name,
        "cloud_provider": deploymentRequest.Provider,
        "region": deploymentRequest.Region,
    }).Info("收到部署请求")
    
    utils.LogInfo(fmt.Sprintf("收到部署请求: 名称=%s, 提供商=%s, 区域=%s",
        deploymentRequest.Name, deploymentRequest.Provider, deploymentRequest.Region))
    
    // 执行Terraform命令
    cmd := exec.Command("terraform", "apply", "-auto-approve")
    output, err := cmd.CombinedOutput()
    
    // 记录Terraform执行情况
    if err != nil {
        logger.WithFields(logrus.Fields{
            "command": "terraform apply",
            "error": err.Error(),
            "output": string(output),
        }).Error("Terraform执行失败")
        
        utils.LogError(fmt.Sprintf("Terraform执行失败: %v", err))
        utils.LogDebug(fmt.Sprintf("Terraform执行详细输出:\n%s", string(output)))
        
        c.JSON(500, gin.H{"error": "部署失败"})
        return
    }
    
    logger.WithFields(logrus.Fields{
        "command": "terraform apply",
        "output": string(output),
    }).Info("Terraform执行成功")
    
    utils.LogInfo("Terraform执行成功")
    utils.LogDebug(fmt.Sprintf("Terraform执行详细输出:\n%s", string(output)))
    
    // 记录成功信息
    logger.WithFields(logrus.Fields{
        "deployment_id": "新生成的ID",
        "status": "success",
    }).Info("部署成功完成")
    
    utils.LogInfo(fmt.Sprintf("部署成功完成, 部署ID: %s", "新生成的ID"))
    
    c.JSON(200, gin.H{"message": "部署成功", "deployment_id": "新生成的ID"})
}

2. 设置环境变量LOG_LEVEL=debug可以获取更详细的日志

3. 日志文件位于logs目录下，包括:
   - multi-cloud-landing-zone-YYYY-MM-DD.log (应用日志)
   - gin-YYYY-MM-DD.log (Gin框架日志)
*/

