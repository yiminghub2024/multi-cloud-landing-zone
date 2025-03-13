package controllers

import (
        "github.com/gin-gonic/gin"
        "github.com/multi-cloud-landing-zone/backend/models"
)

// GetProviders 返回支持的云服务提供商列表
func GetProviders(c *gin.Context) {
        providers := []models.Provider{
                {Name: "AWS", Value: "aws", Logo: "AWS"},
                {Name: "Azure", Value: "azure", Logo: "Azure"},
                {Name: "阿里云", Value: "alicloud", Logo: "阿里云"},
                {Name: "百度云", Value: "baidu", Logo: "百度云"},
                {Name: "华为云", Value: "huawei", Logo: "华为云"},
                {Name: "腾讯云", Value: "tencent", Logo: "腾讯云"},
                {Name: "火山云", Value: "volcengine", Logo: "火山云"},
        }

        c.JSON(200, gin.H{
                "success": true,
                "data":    providers,
        })
}

// GetRegions 返回指定云服务提供商的区域列表
func GetRegions(c *gin.Context) {
        provider := c.Param("provider")
        var regions []models.Region

        switch provider {
        case "aws":
                regions = []models.Region{
                        {Name: "美国东部 (弗吉尼亚北部)", Value: "us-east-1"},
                        {Name: "美国东部 (俄亥俄)", Value: "us-east-2"},
                        {Name: "美国西部 (加利福尼亚北部)", Value: "us-west-1"},
                        {Name: "美国西部 (俄勒冈)", Value: "us-west-2"},
                        {Name: "亚太地区 (香港)", Value: "ap-east-1"},
                        {Name: "亚太地区 (东京)", Value: "ap-northeast-1"},
                }
        case "azure":
                regions = []models.Region{
                        {Name: "美国东部", Value: "eastus"},
                        {Name: "美国东部2", Value: "eastus2"},
                        {Name: "美国西部", Value: "westus"},
                        {Name: "美国西部2", Value: "westus2"},
                        {Name: "东亚", Value: "eastasia"},
                        {Name: "东南亚", Value: "southeastasia"},
                }
        case "alicloud":
                regions = []models.Region{
                        {Name: "华北 1 (青岛)", Value: "cn-qingdao"},
                        {Name: "华北 2 (北京)", Value: "cn-beijing"},
                        {Name: "华北 3 (张家口)", Value: "cn-zhangjiakou"},
                        {Name: "华东 1 (杭州)", Value: "cn-hangzhou"},
                        {Name: "华东 2 (上海)", Value: "cn-shanghai"},
                        {Name: "华南 1 (深圳)", Value: "cn-shenzhen"},
                }
        case "baidu":
                regions = []models.Region{
                        {Name: "华北-北京", Value: "bj"},
                        {Name: "华南-广州", Value: "gz"},
                        {Name: "华东-苏州", Value: "su"},
                }
        case "huawei":
                regions = []models.Region{
                        {Name: "华北-北京一", Value: "cn-north-1"},
                        {Name: "华北-北京四", Value: "cn-north-4"},
                        {Name: "华东-上海一", Value: "cn-east-3"},
                        {Name: "华南-广州", Value: "cn-south-1"},
                        {Name: "亚太-香港", Value: "ap-southeast-1"},
                }
        case "tencent":
                regions = []models.Region{
                        {Name: "华南地区(广州)", Value: "ap-guangzhou"},
                        {Name: "华东地区(上海)", Value: "ap-shanghai"},
                        {Name: "华北地区(北京)", Value: "ap-beijing"},
                        {Name: "西南地区(成都)", Value: "ap-chengdu"},
                        {Name: "西南地区(重庆)", Value: "ap-chongqing"},
                        {Name: "港澳台地区(中国香港)", Value: "ap-hongkong"},
                }
        case "volcengine":
                regions = []models.Region{
                        {Name: "华北-北京", Value: "cn-beijing"},
                        {Name: "华东-上海", Value: "cn-shanghai"},
                        {Name: "华南-广州", Value: "cn-guangzhou"},
                }
        }

        c.JSON(200, gin.H{
                "success": true,
                "data":    regions,
        })
}

// GetAvailabilityZones 返回指定云服务提供商和区域的可用区列表
func GetAvailabilityZones(c *gin.Context) {
        _ = c.Param("provider") // 使用空白标识符解决未使用变量的警告
        region := c.Param("region")
        var azs []models.AvailabilityZone

        // 模拟不同区域的可用区数据
        mockAZs := map[string][]models.AvailabilityZone{
                "us-east-1": {
                        {Name: "us-east-1a", Value: "us-east-1a"},
                        {Name: "us-east-1b", Value: "us-east-1b"},
                        {Name: "us-east-1c", Value: "us-east-1c"},
                },
                "us-west-2": {
                        {Name: "us-west-2a", Value: "us-west-2a"},
                        {Name: "us-west-2b", Value: "us-west-2b"},
                        {Name: "us-west-2c", Value: "us-west-2c"},
                },
                "eastus": {
                        {Name: "eastus-1", Value: "eastus-1"},
                        {Name: "eastus-2", Value: "eastus-2"},
                        {Name: "eastus-3", Value: "eastus-3"},
                },
                "cn-beijing": {
                        {Name: "可用区A", Value: "cn-beijing-a"},
                        {Name: "可用区B", Value: "cn-beijing-b"},
                        {Name: "可用区C", Value: "cn-beijing-c"},
                },
                "cn-shanghai": {
                        {Name: "可用区A", Value: "cn-shanghai-a"},
                        {Name: "可用区B", Value: "cn-shanghai-b"},
                        {Name: "可用区C", Value: "cn-shanghai-c"},
                },
        }

        // 如果有特定区域的数据，使用它，否则生成通用的可用区
        if val, ok := mockAZs[region]; ok {
                azs = val
        } else {
                azs = []models.AvailabilityZone{
                        {Name: "可用区A", Value: region + "-a"},
                        {Name: "可用区B", Value: region + "-b"},
                        {Name: "可用区C", Value: region + "-c"},
                }
        }

        c.JSON(200, gin.H{
                "success": true,
                "data":    azs,
        })
}

// GetComponents 返回指定云服务提供商和区域可用的云组件列表
func GetComponents(c *gin.Context) {
        providerParam := c.Param("provider")
        _ = c.Param("region") // 使用空白标识符解决未使用变量的警告

        // 通用组件
        commonComponents := []models.Component{
                {
                        Name:        "负载均衡器",
                        Value:       "load-balancer",
                        Description: "用于分发网络流量的服务，提高应用程序的可用性和容错能力",
                        Properties: []models.ComponentProperty{
                                {
                                        Name:         "实例数量",
                                        Key:          "instance_count",
                                        Type:         "number",
                                        DefaultValue: "1",
                                        Placeholder:  "请输入实例数量",
                                        Description:  "负载均衡器实例的数量",
                                },
                                {
                                        Name:         "监听端口",
                                        Key:          "listener_port",
                                        Type:         "number",
                                        DefaultValue: "80",
                                        Placeholder:  "请输入监听端口",
                                        Description:  "负载均衡器监听的端口",
                                },
                        },
                },
                {
                        Name:        "对象存储",
                        Value:       "object-storage",
                        Description: "用于存储和检索任意数量数据的服务，适用于静态网站、备份和归档等场景",
                        Properties: []models.ComponentProperty{
                                {
                                        Name:         "存储桶名称",
                                        Key:          "bucket_name",
                                        Type:         "text",
                                        DefaultValue: "",
                                        Placeholder:  "请输入全局唯一的存储桶名称",
                                        Description:  "存储桶名称必须全局唯一",
                                },
                                {
                                        Name:         "存储类型",
                                        Key:          "storage_class",
                                        Type:         "text",
                                        DefaultValue: "Standard",
                                        Placeholder:  "例如: Standard, IA, Archive",
                                        Description:  "存储类型决定数据的访问频率和成本",
                                },
                        },
                },
        }

        // 根据不同的云服务提供商添加特定组件
        var components []models.Component
        components = append(components, commonComponents...)

        // 这里可以根据provider和region添加特定的组件
        if providerParam == "aws" {
                components = append(components, models.Component{
                        Name:        "Lambda函数",
                        Value:       "lambda",
                        Description: "AWS Lambda是一项无服务器计算服务，可运行代码而无需预置或管理服务器",
                        Properties: []models.ComponentProperty{
                                {
                                        Name:         "函数名称",
                                        Key:          "function_name",
                                        Type:         "text",
                                        DefaultValue: "",
                                        Placeholder:  "请输入函数名称",
                                        Description:  "Lambda函数的名称",
                                },
                                {
                                        Name:         "运行时",
                                        Key:          "runtime",
                                        Type:         "text",
                                        DefaultValue: "nodejs14.x",
                                        Placeholder:  "例如: nodejs14.x, python3.9",
                                        Description:  "Lambda函数的运行时环境",
                                },
                                {
                                        Name:         "内存大小",
                                        Key:          "memory_size",
                                        Type:         "number",
                                        DefaultValue: "128",
                                        Placeholder:  "请输入内存大小(MB)",
                                        Description:  "Lambda函数的内存大小，单位为MB",
                                },
                        },
                })
        } else if providerParam == "azure" {
                components = append(components, models.Component{
                        Name:        "Azure Functions",
                        Value:       "azure-functions",
                        Description: "Azure Functions是一项无服务器计算服务，可运行代码而无需预置或管理服务器",
                        Properties: []models.ComponentProperty{
                                {
                                        Name:         "函数名称",
                                        Key:          "function_name",
                                        Type:         "text",
                                        DefaultValue: "",
                                        Placeholder:  "请输入函数名称",
                                        Description:  "Azure Functions的名称",
                                },
                                {
                                        Name:         "运行时",
                                        Key:          "runtime",
                                        Type:         "text",
                                        DefaultValue: "node",
                                        Placeholder:  "例如: node, python, dotnet",
                                        Description:  "Azure Functions的运行时环境",
                                },
                        },
                })
        }

        c.JSON(200, gin.H{
                "success": true,
                "data":    components,
        })
}
