package utils

import (
        "encoding/json"
        "fmt"
        "os"
        "path/filepath"
        "strings"

        "github.com/multi-cloud-landing-zone/backend/models"
)

// 可用区映射，将前端显示名称映射到实际的AWS可用区ID
var awsAvailabilityZoneMap = map[string]map[string]string{
        "us-east-1": {
                "可用区A": "us-east-1a",
                "可用区B": "us-east-1b",
                "可用区C": "us-east-1c",
        },
        "us-west-2": {
                "可用区A": "us-west-2a",
                "可用区B": "us-west-2b",
                "可用区C": "us-west-2c",
        },
        "cn-north-1": {
                "可用区A": "cn-north-1a",
                "可用区B": "cn-north-1b",
                "可用区C": "cn-north-1d",
        },
        "cn-northwest-1": {
                "可用区A": "cn-northwest-1a",
                "可用区B": "cn-northwest-1b",
                "可用区C": "cn-northwest-1c",
        },
}

// 获取实际的AWS可用区ID
func getActualAwsAZ(region string, displayName string) string {
        // 如果已经是有效的可用区ID格式（包含区域名称），则直接返回
        if strings.HasPrefix(displayName, region) {
                return displayName
        }
        
        // 检查是否有该区域的映射
        if regionMap, ok := awsAvailabilityZoneMap[region]; ok {
                // 检查是否有该显示名称的映射
                if actualAZ, ok := regionMap[displayName]; ok {
                        LogInfo(fmt.Sprintf("将可用区显示名称 '%s' 映射到实际AWS可用区ID '%s'", displayName, actualAZ))
                        return actualAZ
                }
        }
        
        // 如果没有找到映射，使用默认格式（区域名称+字母a）
        defaultAZ := region + "a"
        LogWarn(fmt.Sprintf("未找到可用区 '%s' 在区域 '%s' 的映射，使用默认可用区 '%s'", displayName, region, defaultAZ))
        return defaultAZ
}

// GenerateTerraformConfig 根据部署配置生成Terraform配置
func GenerateTerraformConfig(config models.DeploymentConfig) string {
        // 记录开始生成Terraform配置
        LogInfo(fmt.Sprintf("开始为云提供商 %s 生成Terraform配置", config.CloudProvider))
        
        // 记录详细的部署配置参数
        configJSON, _ := json.MarshalIndent(config, "", "  ")
        LogInfo(fmt.Sprintf("部署配置详情:\n%s", string(configJSON)))
        
        var terraformConfig strings.Builder

        // 添加提供商配置
        switch config.CloudProvider {
        case "aws":
                terraformConfig.WriteString(fmt.Sprintf(`provider "aws" {
  region = "%s"
}

`, config.Region))
        case "azure":
                terraformConfig.WriteString(`provider "azurerm" {
  features {}
}

`)
        case "alicloud":
                terraformConfig.WriteString(fmt.Sprintf(`provider "alicloud" {
  region = "%s"
}

`, config.Region))
        case "baidu":
                terraformConfig.WriteString(fmt.Sprintf(`provider "baiducloud" {
  region = "%s"
}

`, config.Region))
        case "huawei":
                terraformConfig.WriteString(fmt.Sprintf(`provider "huaweicloud" {
  region = "%s"
}

`, config.Region))
        case "tencent":
                terraformConfig.WriteString(fmt.Sprintf(`provider "tencentcloud" {
  region = "%s"
}

`, config.Region))
        case "volcengine":
                terraformConfig.WriteString(fmt.Sprintf(`provider "volcengine" {
  region = "%s"
}

`, config.Region))
        default:
                terraformConfig.WriteString(fmt.Sprintf(`provider "%s" {
  region = "%s"
}

`, config.CloudProvider, config.Region))
        }

        LogInfo(fmt.Sprintf("已生成云提供商配置: %s, 区域: %s", config.CloudProvider, config.Region))

        // 添加VPC配置
        if config.CloudProvider == "aws" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "aws_vpc" "%s" {
  cidr_block           = "%s"
  enable_dns_support   = %t
  enable_dns_hostnames = %t
  
  tags = {
    Name = "%s"
  }
}

`, config.VPC.Name, config.VPC.CIDR, config.VPC.EnableDnsSupport, config.VPC.EnableDnsHostnames, config.VPC.Name))
        } else if config.CloudProvider == "azure" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "azurerm_resource_group" "rg" {
  name     = "rg-%s"
  location = "%s"
}

resource "azurerm_virtual_network" "%s" {
  name                = "%s"
  address_space       = ["%s"]
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
}

`, config.VPC.Name, config.Region, config.VPC.Name, config.VPC.Name, config.VPC.CIDR))
        } else if config.CloudProvider == "alicloud" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "alicloud_vpc" "%s" {
  vpc_name   = "%s"
  cidr_block = "%s"
}

`, config.VPC.Name, config.VPC.Name, config.VPC.CIDR))
        } else if config.CloudProvider == "baidu" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "baiducloud_vpc" "%s" {
  name       = "%s"
  cidr_block = "%s"
}

`, config.VPC.Name, config.VPC.Name, config.VPC.CIDR))
        } else if config.CloudProvider == "huawei" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "huaweicloud_vpc" "%s" {
  name        = "%s"
  cidr        = "%s"
  description = "VPC created by multi-cloud landing zone platform"
}

`, config.VPC.Name, config.VPC.Name, config.VPC.CIDR))
        } else if config.CloudProvider == "tencent" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "tencentcloud_vpc" "%s" {
  name       = "%s"
  cidr_block = "%s"
}

`, config.VPC.Name, config.VPC.Name, config.VPC.CIDR))
        } else if config.CloudProvider == "volcengine" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "volcengine_vpc" "%s" {
  vpc_name   = "%s"
  cidr_block = "%s"
}

`, config.VPC.Name, config.VPC.Name, config.VPC.CIDR))
        } else {
                terraformConfig.WriteString(fmt.Sprintf(`resource "%s_vpc" "%s" {
  name       = "%s"
  cidr_block = "%s"
}

`, config.CloudProvider, config.VPC.Name, config.VPC.Name, config.VPC.CIDR))
        }

        LogInfo(fmt.Sprintf("已生成VPC配置: 名称=%s, CIDR=%s", config.VPC.Name, config.VPC.CIDR))

        // 添加子网配置
        if config.CloudProvider == "aws" {
                // 获取实际的AWS可用区ID
                actualAZ := getActualAwsAZ(config.Region, config.AZ)
                
                terraformConfig.WriteString(fmt.Sprintf(`resource "aws_subnet" "%s" {
  vpc_id                  = aws_vpc.%s.id
  cidr_block              = "%s"
  availability_zone       = "%s"
  map_public_ip_on_launch = %t
  
  tags = {
    Name = "%s"
  }
}

`, config.Subnet.Name, config.VPC.Name, config.Subnet.CIDR, actualAZ, config.Subnet.MapPublicIpOnLaunch, config.Subnet.Name))
        } else if config.CloudProvider == "azure" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "azurerm_subnet" "%s" {
  name                 = "%s"
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.%s.name
  address_prefixes     = ["%s"]
}

`, config.Subnet.Name, config.Subnet.Name, config.VPC.Name, config.Subnet.CIDR))
        } else if config.CloudProvider == "alicloud" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "alicloud_vswitch" "%s" {
  vpc_id     = alicloud_vpc.%s.id
  cidr_block = "%s"
  zone_id    = "%s"
  name       = "%s"
}

`, config.Subnet.Name, config.VPC.Name, config.Subnet.CIDR, config.AZ, config.Subnet.Name))
        } else if config.CloudProvider == "baidu" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "baiducloud_subnet" "%s" {
  name        = "%s"
  zone_name   = "%s"
  cidr        = "%s"
  vpc_id      = baiducloud_vpc.%s.id
  description = "Subnet created by multi-cloud landing zone platform"
}

`, config.Subnet.Name, config.Subnet.Name, config.AZ, config.Subnet.CIDR, config.VPC.Name))
        } else if config.CloudProvider == "huawei" {
                // 从CIDR中提取网关IP
                cidrParts := strings.Split(config.Subnet.CIDR, "/")
                ipParts := strings.Split(cidrParts[0], ".")
                gatewayIP := fmt.Sprintf("%s.%s.%s.1", ipParts[0], ipParts[1], ipParts[2])

                terraformConfig.WriteString(fmt.Sprintf(`resource "huaweicloud_vpc_subnet" "%s" {
  name       = "%s"
  cidr       = "%s"
  gateway_ip = "%s"
  vpc_id     = huaweicloud_vpc.%s.id
}

`, config.Subnet.Name, config.Subnet.Name, config.Subnet.CIDR, gatewayIP, config.VPC.Name))
        } else if config.CloudProvider == "tencent" {
                // 对于腾讯云，也需要处理可用区映射
                var actualAZ string
                if config.CloudProvider == "tencent" && strings.HasPrefix(config.Region, "ap-") {
                    // 简单处理腾讯云可用区
                    if strings.Contains(config.AZ, "可用区") {
                        zoneChar := ""
                        if strings.Contains(config.AZ, "A") || strings.Contains(config.AZ, "a") {
                            zoneChar = "1"
                        } else if strings.Contains(config.AZ, "B") || strings.Contains(config.AZ, "b") {
                            zoneChar = "2"
                        } else if strings.Contains(config.AZ, "C") || strings.Contains(config.AZ, "c") {
                            zoneChar = "3"
                        }
                        actualAZ = config.Region + "-" + zoneChar
                    } else {
                        actualAZ = config.AZ
                    }
                } else {
                    actualAZ = config.AZ
                }
                
                terraformConfig.WriteString(fmt.Sprintf(`resource "tencentcloud_subnet" "%s" {
  name              = "%s"
  vpc_id            = tencentcloud_vpc.%s.id
  cidr_block        = "%s"
  availability_zone = "%s"
}

`, config.Subnet.Name, config.Subnet.Name, config.VPC.Name, config.Subnet.CIDR, actualAZ))
        } else if config.CloudProvider == "volcengine" {
                terraformConfig.WriteString(fmt.Sprintf(`resource "volcengine_subnet" "%s" {
  subnet_name = "%s"
  cidr_block  = "%s"
  zone_id     = "%s"
  vpc_id      = volcengine_vpc.%s.id
}

`, config.Subnet.Name, config.Subnet.Name, config.Subnet.CIDR, config.AZ, config.VPC.Name))
        } else {
                terraformConfig.WriteString(fmt.Sprintf(`resource "%s_subnet" "%s" {
  vpc_id     = %s_vpc.%s.id
  cidr_block = "%s"
  zone_id    = "%s"
}

`, config.CloudProvider, config.Subnet.Name, config.CloudProvider, config.VPC.Name, config.Subnet.CIDR, config.AZ))
        }

        LogInfo(fmt.Sprintf("已生成子网配置: 名称=%s, CIDR=%s, 可用区=%s", config.Subnet.Name, config.Subnet.CIDR, config.AZ))

        // 添加选定的组件
        for _, component := range config.Components {
                props := config.ComponentProperties[component.Value]
                LogInfo(fmt.Sprintf("正在生成组件配置: %s (类型: %s)", component.Name, component.Value))
                
                // 记录组件属性
                propsJSON, _ := json.MarshalIndent(props, "", "  ")
                LogInfo(fmt.Sprintf("组件属性:\n%s", string(propsJSON)))

                switch component.Value {
                case "load-balancer":
                        if config.CloudProvider == "aws" {
                                terraformConfig.WriteString(fmt.Sprintf(`resource "aws_lb" "load_balancer" {
  name               = "%s-%s"
  internal           = false
  load_balancer_type = "application"
  subnets            = [aws_subnet.%s.id]
  
  enable_deletion_protection = false
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.load_balancer.arn
  port              = "%s"
  protocol          = "HTTP"
  
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.front_end.arn
  }
}

resource "aws_lb_target_group" "front_end" {
  name     = "tf-lb-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.%s.id
}

`, component.Name, config.VPC.Name, config.Subnet.Name, props["listener_port"], config.VPC.Name))
                        }
                        // 其他云提供商的负载均衡器配置...
                        
                case "transit-gateway":
                        if config.CloudProvider == "aws" {
                                // 创建Transit Gateway
                                terraformConfig.WriteString(fmt.Sprintf(`resource "aws_ec2_transit_gateway" "tgw" {
  description                     = "%s"
  auto_accept_shared_attachments  = "%s"
  dns_support                     = "%s"
  vpn_ecmp_support                = "%s"
  
  tags = {
    Name = "transit-gateway"
  }
}

`, props["description"], props["auto_accept_shared_attachments"], props["dns_support"], props["vpn_ecmp_support"]))

                                // 如果启用了路由表配置
                                if config.ComponentConfig.EnableRouteTables {
                                        tgwConfig := config.ComponentConfig.TransitGatewayConfig
                                        defaultAssociation := "false"
                                        if tgwConfig.DefaultRouteTable {
                                                defaultAssociation = "true"
                                        }
                                        
                                        terraformConfig.WriteString(fmt.Sprintf(`resource "aws_ec2_transit_gateway_route_table" "tgw_rt" {
  transit_gateway_id = aws_ec2_transit_gateway.tgw.id
  
  tags = {
    Name = "%s"
  }
}

resource "aws_ec2_transit_gateway_route_table_association" "tgw_rt_assoc" {
  transit_gateway_attachment_id  = aws_ec2_transit_gateway_vpc_attachment.tgw_attachment[0].id
  transit_gateway_route_table_id = aws_ec2_transit_gateway_route_table.tgw_rt.id
}

`, tgwConfig.RouteTableName))
                                }

                                // 如果启用了VPC附件
                                if config.ComponentConfig.EnableVpcAttachment {
                                        tgwConfig := config.ComponentConfig.TransitGatewayConfig
                                        dnsSupport := "disable"
                                        ipv6Support := "disable"
                                        
                                        if tgwConfig.DnsSupport {
                                                dnsSupport = "enable"
                                        }
                                        
                                        if tgwConfig.Ipv6Support {
                                                ipv6Support = "enable"
                                        }
                                        
                                        // 使用子网ID或默认使用当前子网
                                        subnetIds := tgwConfig.SubnetIds
                                        if subnetIds == "" {
                                                subnetIds = fmt.Sprintf("aws_subnet.%s.id", config.Subnet.Name)
                                        }
                                        
                                        terraformConfig.WriteString(fmt.Sprintf(`resource "aws_ec2_transit_gateway_vpc_attachment" "tgw_attachment" {
  transit_gateway_id = aws_ec2_transit_gateway.tgw.id
  vpc_id             = aws_vpc.%s.id
  subnet_ids         = [%s]
  
  dns_support        = "%s"
  ipv6_support       = "%s"
  
  tags = {
    Name = "tgw-attachment"
  }
}

`, config.VPC.Name, subnetIds, dnsSupport, ipv6Support))
                                }
                        }
                        // 其他云提供商的Transit Gateway配置...
                        
                case "object-storage":
                        if config.CloudProvider == "aws" {
                                bucketName := props["bucket_name"]
                                terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket" "storage" {
  bucket = "%s"
  
  tags = {
    Name = "%s"
  }
}

`, bucketName, bucketName))

                                // 根据组件配置添加存储桶策略
                                bucketPolicyType := config.ComponentConfig.BucketPolicyType
                                if bucketPolicyType == "public-read" {
                                        terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket_acl" "storage_acl" {
  bucket = aws_s3_bucket.storage.id
  acl    = "public-read"
}

`))
                                } else if bucketPolicyType == "public-read-write" {
                                        terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket_acl" "storage_acl" {
  bucket = aws_s3_bucket.storage.id
  acl    = "public-read-write"
}

`))
                                } else if bucketPolicyType == "private" {
                                        terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket_acl" "storage_acl" {
  bucket = aws_s3_bucket.storage.id
  acl    = "private"
}

`))
                                } else if bucketPolicyType == "custom" && config.ComponentConfig.CustomBucketPolicy != "" {
                                        // 自定义策略
                                        terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket_policy" "storage_policy" {
  bucket = aws_s3_bucket.storage.id
  policy = <<POLICY
%s
POLICY
}

`, config.ComponentConfig.CustomBucketPolicy))
                                }

                                // 添加生命周期规则
                                if config.ComponentConfig.EnableLifecycleRules {
                                        lifecycleRule := config.ComponentConfig.LifecycleRule
                                        terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket_lifecycle_configuration" "storage_lifecycle" {
  bucket = aws_s3_bucket.storage.id

  rule {
    id     = "%s"
    status = "%s"

    expiration {
      days = %d
    }

    transition {
      days          = %d
      storage_class = "STANDARD_IA"
    }
  }
}

`, lifecycleRule.Name, lifecycleRule.Status, lifecycleRule.ExpirationDays, lifecycleRule.TransitionDays))
                                }
                        }
                        // 其他云提供商的对象存储配置...
                }
        }

        return terraformConfig.String()
}

// SaveTerraformConfig 保存Terraform配置到文件
func SaveTerraformConfig(config string, filePath string) error {
        // 创建目录（如果不存在）
        dir := filepath.Dir(filePath)
        if err := os.MkdirAll(dir, 0755); err != nil {
                LogError(fmt.Sprintf("创建目录失败: %v", err))
                return err
        }

        // 写入文件
        if err := os.WriteFile(filePath, []byte(config), 0644); err != nil {
                LogError(fmt.Sprintf("写入Terraform配置文件失败: %v", err))
                return err
        }

        LogInfo(fmt.Sprintf("Terraform配置已保存到: %s", filePath))
        return nil
}

// GenerateTopology 生成资源拓扑图
func GenerateTopology(config models.DeploymentConfig) map[string]interface{} {
        // 简单的拓扑图生成
        topology := map[string]interface{}{
                "nodes": []map[string]interface{}{
                        {
                                "id":    "vpc",
                                "type":  "vpc",
                                "name":  config.VPC.Name,
                                "cidr":  config.VPC.CIDR,
                                "cloud": config.CloudProvider,
                        },
                        {
                                "id":    "subnet",
                                "type":  "subnet",
                                "name":  config.Subnet.Name,
                                "cidr":  config.Subnet.CIDR,
                                "az":    config.AZ,
                                "cloud": config.CloudProvider,
                        },
                },
                "edges": []map[string]interface{}{
                        {
                                "source": "subnet",
                                "target": "vpc",
                                "type":   "contains",
                        },
                },
        }

        // 添加组件节点
        for i, component := range config.Components {
                nodeID := fmt.Sprintf("component-%d", i)
                props := config.ComponentProperties[component.Value]

                // 创建组件节点
                componentNode := map[string]interface{}{
                        "id":    nodeID,
                        "type":  component.Value,
                        "name":  component.Name,
                        "cloud": config.CloudProvider,
                }

                // 添加组件特定属性
                for k, v := range props {
                        componentNode[k] = v
                }

                // 添加节点
                topology["nodes"] = append(topology["nodes"].([]map[string]interface{}), componentNode)

                // 添加边（连接到子网）
                topology["edges"] = append(topology["edges"].([]map[string]interface{}), map[string]interface{}{
                        "source": nodeID,
                        "target": "subnet",
                        "type":   "deployed_in",
                })
        }

        return topology
}
