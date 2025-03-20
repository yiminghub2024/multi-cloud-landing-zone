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
        
        // 配置CORS - 修复CORS配置，不再使用AllowAllOrigins和AllowCredentials的冲突配置
        router.Use(cors.New(cors.Config{
                AllowOrigins:     []string{"http://localhost", "http://localhost:8080", "http://127.0.0.1", "http://127.0.0.1:8080"},
                AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
                AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
                ExposeHeaders:    []string{"Content-Length"},
                AllowCredentials: true,
                MaxAge:           12 * time.Hour,
        }))
        
        Logger.Info("Gin路由初始化完成")
        utils.LogInfo("Gin路由初始化完成")

        // 配置CORS - 使用自定义中间件处理动态Origin
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
            
            // 对于OPTIONS请求，立即返回200
            if c.Request.Method == "OPTIONS" {
                Logger.WithFields(logrus.Fields{
                    "request_method": c.Request.Method,
                    "request_path": c.Request.URL.Path,
                    "origin": origin,
                    "cors_headers": c.Writer.Header(),
                }).Info("处理OPTIONS预检请求")
                
                c.AbortWithStatus(200)
                return
            }
            
            // 继续处理非OPTIONS请求
            c.Next()
            
            // 记录CORS响应头
            Logger.WithFields(logrus.Fields{
                "response_status": c.Writer.Status(),
                "response_headers": c.Writer.Header(),
                "request_path": c.Request.URL.Path,
            }).Info("CORS请求处理完成")
        })
        
        Logger.Info("CORS配置完成 - 使用增强的自定义中间件，支持动态Origin和凭证")
        utils.LogInfo("CORS配置完成 - 使用增强的自定义中间件，支持动态Origin和凭证")

        // 设置API路由
        routes.SetupRoutes(router)
        Logger.Info("API路由设置完成")
        utils.LogInfo("API路由设置完成")

        // 获取端口配置
        port := os.Getenv("PORT")
        if port == "" {
                port = "3000" // 默认端口
        }

        // 启动服务器
        serverAddr := fmt.Sprintf("0.0.0.0:%s", port)
        Logger.Info("服务器启动，监听地址: " + serverAddr)
        utils.LogInfo("服务器启动，监听地址: " + serverAddr)
        
        if err := router.Run(serverAddr); err != nil {
                Logger.Fatalf("服务器启动失败: %v", err)
	utils.LogError(fmt.Sprintf("服务器启动失败: %v", err))        
}
}
