package utils

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "time"
)

// CORS配置工具函数

// GetAllowedOrigins 返回允许的源列表
func GetAllowedOrigins() []string {
	// 从环境变量获取允许的源，如果没有设置，使用默认值
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsEnv != "" {
		// 解析环境变量中的源列表
		return strings.Split(allowedOriginsEnv, ",")
	}
	
	// 默认允许的源
	return []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		"http://10.168.0.5",
		"http://10.168.0.5:8080",
		"http://10.168.0.5:3000",
	}
}

// LogCORSRequest 记录CORS请求详情
func LogCORSRequest(method, path, origin string, headers map[string][]string) {
	LogInfo(fmt.Sprintf("CORS请求: 方法=%s, 路径=%s, 来源=%s", method, path, origin))
	
	// 检查是否是预检请求
	if method == "OPTIONS" {
		LogInfo("收到CORS预检请求")
		
		// 记录预检请求的特殊头部
		acrMethod := getHeaderValue(headers, "Access-Control-Request-Method")
		acrHeaders := getHeaderValue(headers, "Access-Control-Request-Headers")
		
		if acrMethod != "" {
			LogInfo(fmt.Sprintf("预检请求方法: %s", acrMethod))
		}
		
		if acrHeaders != "" {
			LogInfo(fmt.Sprintf("预检请求头部: %s", acrHeaders))
		}
	}
	
	// 检查源是否在允许列表中
	allowedOrigins := GetAllowedOrigins()
	originAllowed := false
	
	for _, allowed := range allowedOrigins {
		if allowed == origin || allowed == "*" {
			originAllowed = true
			break
		}
	}
	
	if !originAllowed && origin != "" {
		LogWarn(fmt.Sprintf("请求来源不在允许列表中: %s", origin))
		LogInfo(fmt.Sprintf("允许的来源列表: %v", allowedOrigins))
	}
}

// LogCORSResponse 记录CORS响应详情
func LogCORSResponse(status int, headers map[string][]string) {
	LogInfo(fmt.Sprintf("CORS响应: 状态码=%d", status))
	
	// 记录关键CORS响应头
	acao := getHeaderValue(headers, "Access-Control-Allow-Origin")
	acac := getHeaderValue(headers, "Access-Control-Allow-Credentials")
	acam := getHeaderValue(headers, "Access-Control-Allow-Methods")
	acah := getHeaderValue(headers, "Access-Control-Allow-Headers")
	
	if acao != "" {
		LogInfo(fmt.Sprintf("响应头 Access-Control-Allow-Origin: %s", acao))
	} else {
		LogWarn("响应缺少 Access-Control-Allow-Origin 头部")
	}
	
	if acac != "" {
		LogInfo(fmt.Sprintf("响应头 Access-Control-Allow-Credentials: %s", acac))
	}
	
	if acam != "" {
		LogInfo(fmt.Sprintf("响应头 Access-Control-Allow-Methods: %s", acam))
	}
	
	if acah != "" {
		LogInfo(fmt.Sprintf("响应头 Access-Control-Allow-Headers: %s", acah))
	}
}

// 从头部映射中获取指定头部的值
func getHeaderValue(headers map[string][]string, key string) string {
	if values, exists := headers[key]; exists && len(values) > 0 {
		return values[0]
	}
	return ""
}

// IsCORSError 检查错误是否与CORS相关
func IsCORSError(err error) bool {
	if err == nil {
		return false
	}
	
	errMsg := err.Error()
	return strings.Contains(errMsg, "CORS") || 
		   strings.Contains(errMsg, "跨源") || 
		   strings.Contains(errMsg, "cross-origin") ||
		   strings.Contains(errMsg, "Access-Control-Allow-Origin")
}

// LogCORSError 记录CORS错误详情
func LogCORSError(err error, method, path, origin string) {
	if !IsCORSError(err) {
		return
	}
	
	LogError(fmt.Sprintf("CORS错误: %v", err))
	LogError(fmt.Sprintf("CORS错误详情: 方法=%s, 路径=%s, 来源=%s", method, path, origin))
	
	// 提供排查建议
	LogInfo("CORS错误排查建议:")
	LogInfo("1. 确认后端CORS配置中包含请求的源")
	LogInfo("2. 检查预检请求(OPTIONS)是否正确处理")
	LogInfo("3. 如果使用credentials: 'include'，确保不使用通配符'*'作为Access-Control-Allow-Origin")
	LogInfo("4. 检查Nginx配置是否正确转发CORS头部")
	LogInfo(fmt.Sprintf("5. 当前允许的源列表: %v", GetAllowedOrigins()))
}
