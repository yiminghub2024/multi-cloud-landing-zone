package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/your-org/multi-cloud-landing-zone/models"
)

// GenerateTerraformConfig 生成Terraform配置
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
	
	// 添加VPC配置
	if config.CloudProvider == "aws" {
		// 优先使用AllVpcs数组，如果存在的话
		if config.AllVpcs != nil && len(config.AllVpcs) > 0 {
			LogInfo(fmt.Sprintf("检测到多个VPC配置，共 %d 个", len(config.AllVpcs)))
			
			for i, vpc := range config.AllVpcs {
				// 使用用户提供的VPC名称作为资源名称
				terraformConfig.WriteString(fmt.Sprintf(`resource "aws_vpc" "%s" {
  cidr_block           = "%s"
  enable_dns_support   = %t
  enable_dns_hostnames = %t
  
  tags = {
    Name = "%s"
  }
}
`, vpc.Name, vpc.CIDR, vpc.EnableDnsSupport, vpc.EnableDnsHostnames, vpc.Name))
				
				LogInfo(fmt.Sprintf("已生成VPC配置 %d: 名称=%s, CIDR=%s", i+1, vpc.Name, vpc.CIDR))
			}
		} else {
			// 兼容旧版本，使用单个VPC配置
			terraformConfig.WriteString(fmt.Sprintf(`resource "aws_vpc" "%s" {
  cidr_block           = "%s"
  enable_dns_support   = %t
  enable_dns_hostnames = %t
  
  tags = {
    Name = "%s"
  }
}
`, config.VPC.Name, config.VPC.CIDR, config.VPC.EnableDnsSupport, config.VPC.EnableDnsHostnames, config.VPC.Name))
			
			LogInfo(fmt.Sprintf("已生成单个VPC配置: 名称=%s, CIDR=%s", config.VPC.Name, config.VPC.CIDR))
		}
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
		// 优先使用AllSubnets数组，如果存在的话
		if config.AllSubnets != nil && len(config.AllSubnets) > 0 {
			LogInfo(fmt.Sprintf("检测到多个子网配置，共 %d 个", len(config.AllSubnets)))
			
			for i, subnet := range config.AllSubnets {
				// 获取实际的AWS可用区ID
				actualAZ := getActualAwsAZ(config.Region, subnet.AZ)
				if actualAZ == "" {
					actualAZ = getActualAwsAZ(config.Region, config.AZ)
				}
				
				// 确定子网所属的VPC
				var vpcName string
				if config.AllVpcs != nil && len(config.AllVpcs) > 0 {
					// 如果有多个VPC，使用对应的VPC名称
					vpcIndex := subnet.VpcIndex // 使用子网中的VpcIndex字段
					if vpcIndex >= 0 && vpcIndex < len(config.AllVpcs) {
						vpcName = config.AllVpcs[vpcIndex].Name
					} else {
						// 默认使用第一个VPC
						vpcName = config.AllVpcs[0].Name
					}
				} else {
					// 使用单个VPC的名称
					vpcName = config.VPC.Name
				}
				
				terraformConfig.WriteString(fmt.Sprintf(`resource "aws_subnet" "%s" {
  vpc_id                  = aws_vpc.%s.id
  cidr_block              = "%s"
  availability_zone       = "%s"
  map_public_ip_on_launch = %t
  
  tags = {
    Name = "%s"
  }
}
`, subnet.Name, vpcName, subnet.CIDR, actualAZ, subnet.MapPublicIpOnLaunch, subnet.Name))
				
				LogInfo(fmt.Sprintf("已生成子网配置 %d: 名称=%s, CIDR=%s, VPC=%s", i+1, subnet.Name, subnet.CIDR, vpcName))
			}
		} else {
			// 兼容旧版本，使用单个子网配置
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
			
			LogInfo(fmt.Sprintf("已生成单个子网配置: 名称=%s, CIDR=%s, VPC=%s", config.Subnet.Name, config.Subnet.CIDR, config.VPC.Name))
		}
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
		terraformConfig.WriteString(fmt.Sprintf(`resource "tencentcloud_subnet" "%s" {
  name              = "%s"
  vpc_id            = tencentcloud_vpc.%s.id
  cidr_block        = "%s"
  availability_zone = "%s"
}
`, config.Subnet.Name, config.Subnet.Name, config.VPC.Name, config.Subnet.CIDR, config.AZ))
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
  name       = "%s"
  vpc_id     = %s_vpc.%s.id
  cidr_block = "%s"
  zone       = "%s"
}
`, config.CloudProvider, config.Subnet.Name, config.Subnet.Name, config.CloudProvider, config.VPC.Name, config.Subnet.CIDR, config.AZ))
	}
	
	// 添加组件配置
	for _, component := range config.Components {
		LogInfo(fmt.Sprintf("处理组件: %s", component))
		
		// 获取组件属性
		var propsMap map[string]interface{}
		if config.ComponentProps != nil {
			if props, ok := config.ComponentProps[component]; ok {
				propsMap = props
			}
		}
		if propsMap == nil {
			propsMap = make(map[string]interface{})
		}
		
		switch component {
		case "ec2":
			if config.CloudProvider == "aws" {
				// 设置默认值
				instanceType := "t2.micro"
				amiId := "ami-0c55b159cbfafe1f0" // Amazon Linux 2 AMI ID
				
				// 如果属性存在且不为nil，则使用属性值
				if instType, ok := propsMap["instance_type"]; ok && instType != nil {
					if instTypeStr, ok := instType.(string); ok && instTypeStr != "" {
						instanceType = instTypeStr
					}
				}
				
				if ami, ok := propsMap["ami_id"]; ok && ami != nil {
					if amiStr, ok := ami.(string); ok && amiStr != "" {
						amiId = amiStr
					}
				}
				
				terraformConfig.WriteString(fmt.Sprintf(`resource "aws_instance" "ec2" {
  ami           = "%s"
  instance_type = "%s"
  subnet_id     = aws_subnet.%s.id
  
  tags = {
    Name = "EC2 Instance"
  }
}
`, amiId, instanceType, config.Subnet.Name))
			}
			// 其他云提供商的EC2配置...
			
		case "rds":
			if config.CloudProvider == "aws" {
				// 设置默认值
				instanceClass := "db.t2.micro"
				engine := "mysql"
				engineVersion := "5.7"
				dbName := "mydb"
				username := "admin"
				password := "password123!"
				
				// 如果属性存在且不为nil，则使用属性值
				if instClass, ok := propsMap["instance_class"]; ok && instClass != nil {
					if instClassStr, ok := instClass.(string); ok && instClassStr != "" {
						instanceClass = instClassStr
					}
				}
				
				if eng, ok := propsMap["engine"]; ok && eng != nil {
					if engStr, ok := eng.(string); ok && engStr != "" {
						engine = engStr
					}
				}
				
				if engVer, ok := propsMap["engine_version"]; ok && engVer != nil {
					if engVerStr, ok := engVer.(string); ok && engVerStr != "" {
						engineVersion = engVerStr
					}
				}
				
				if db, ok := propsMap["db_name"]; ok && db != nil {
					if dbStr, ok := db.(string); ok && dbStr != "" {
						dbName = dbStr
					}
				}
				
				if user, ok := propsMap["username"]; ok && user != nil {
					if userStr, ok := user.(string); ok && userStr != "" {
						username = userStr
					}
				}
				
				if pass, ok := propsMap["password"]; ok && pass != nil {
					if passStr, ok := pass.(string); ok && passStr != "" {
						password = passStr
					}
				}
				
				terraformConfig.WriteString(fmt.Sprintf(`resource "aws_db_subnet_group" "default" {
  name       = "main"
  subnet_ids = [aws_subnet.%s.id]
  
  tags = {
    Name = "My DB subnet group"
  }
}

resource "aws_db_instance" "default" {
  allocated_storage    = 10
  db_name              = "%s"
  engine               = "%s"
  engine_version       = "%s"
  instance_class       = "%s"
  username             = "%s"
  password             = "%s"
  parameter_group_name = "default.mysql5.7"
  skip_final_snapshot  = true
  db_subnet_group_name = aws_db_subnet_group.default.name
}
`, config.Subnet.Name, dbName, engine, engineVersion, instanceClass, username, password))
			}
			// 其他云提供商的RDS配置...
			
		case "elb":
			if config.CloudProvider == "aws" {
				// 设置默认值
				lbType := "application"
				listenerPort := 80
				
				// 如果属性存在且不为nil，则使用属性值
				if lb, ok := propsMap["lb_type"]; ok && lb != nil {
					if lbStr, ok := lb.(string); ok && lbStr != "" {
						lbType = lbStr
					}
				}
				
				if port, ok := propsMap["listener_port"]; ok && port != nil {
					if portFloat, ok := port.(float64); ok {
						listenerPort = int(portFloat)
					}
				}
				
				terraformConfig.WriteString(fmt.Sprintf(`resource "aws_lb" "%s" {
  name               = "%s"
  internal           = false
  load_balancer_type = "%s"
  security_groups    = [aws_security_group.lb_sg.id]
  subnets            = [aws_subnet.%s.id]
  
  enable_deletion_protection = false
  
  tags = {
    Environment = "production"
  }
}

resource "aws_security_group" "lb_sg" {
  name        = "allow_http"
  description = "Allow HTTP inbound traffic"
  vpc_id      = aws_vpc.%s.id
  
  ingress {
    description = "HTTP from VPC"
    from_port   = %d
    to_port     = %d
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  tags = {
    Name = "allow_http"
  }
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.%s.arn
  port              = "%d"
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
`, component, component, lbType, config.Subnet.Name, config.VPC.Name, listenerPort, listenerPort, component, listenerPort, config.VPC.Name))
			}
			// 其他云提供商的负载均衡器配置...
			
		case "transit-gateway":
			if config.CloudProvider == "aws" {
				// 创建Transit Gateway
				// 设置默认值，避免nil值导致格式化错误
				description := "Transit Gateway for multi-cloud connectivity"
				autoAcceptSharedAttachments := "disable"
				dnsSupport := "enable"
				vpnEcmpSupport := "disable"
				tgwName := "transit-gateway"
				
				// 如果属性存在且不为nil，则使用属性值
				if desc, ok := propsMap["description"]; ok && desc != nil {
					if descStr, ok := desc.(string); ok && descStr != "" {
						description = descStr
					}
				}
				
				if autoAccept, ok := propsMap["auto_accept_shared_attachments"]; ok && autoAccept != nil {
					if autoAcceptStr, ok := autoAccept.(string); ok && autoAcceptStr != "" {
						autoAcceptSharedAttachments = autoAcceptStr
					}
				}
				
				if dns, ok := propsMap["dns_support"]; ok && dns != nil {
					if dnsStr, ok := dns.(string); ok && dnsStr != "" {
						dnsSupport = dnsStr
					}
				}
				
				if vpnEcmp, ok := propsMap["vpn_ecmp_support"]; ok && vpnEcmp != nil {
					if vpnEcmpStr, ok := vpnEcmp.(string); ok && vpnEcmpStr != "" {
						vpnEcmpSupport = vpnEcmpStr
					}
				}
				
				if name, ok := propsMap["name"]; ok && name != nil {
					if nameStr, ok := name.(string); ok && nameStr != "" {
						tgwName = nameStr
					}
				}
				
				// 使用ComponentConfig中的传输网关名称（如果存在）
				if config.ComponentConfig.TransitGatewayName != "" {
					tgwName = config.ComponentConfig.TransitGatewayName
				}
				
				terraformConfig.WriteString(fmt.Sprintf(`resource "aws_ec2_transit_gateway" "tgw" {
  description                     = "%s"
  auto_accept_shared_attachments  = "%s"
  dns_support                     = "%s"
  vpn_ecmp_support                = "%s"
  
  tags = {
    Name = "%s"
  }
}
`, description, autoAcceptSharedAttachments, dnsSupport, vpnEcmpSupport, tgwName))
				
				// 处理中转网关挂载配置
				// 检查ComponentConfig中是否有tgwAttachments数组
				var tgwAttachmentsJSON []byte
				if attachments, ok := propsMap["tgwAttachments"]; ok && attachments != nil {
					// 记录中转网关挂载配置
					attachmentsStr, ok := attachments.(string)
					if ok && attachmentsStr != "" {
						LogInfo(fmt.Sprintf("处理中转网关挂载配置: %s", attachmentsStr))
						tgwAttachmentsJSON = []byte(attachmentsStr)
					}
				}
				
				// 解析中转网关挂载配置
				if len(tgwAttachmentsJSON) > 0 {
					var attachments []map[string]interface{}
					if err := json.Unmarshal(tgwAttachmentsJSON, &attachments); err == nil {
						for i, attachment := range attachments {
							// 获取挂载配置
							vpcId, _ := attachment["vpcId"].(string)
							subnetIds, _ := attachment["subnetIds"].(string)
							attachmentName, _ := attachment["name"].(string)
							
							if attachmentName == "" {
								attachmentName = fmt.Sprintf("tgw-attachment-%d", i+1)
							}
							
							// 生成中转网关挂载配置
							terraformConfig.WriteString(fmt.Sprintf(`resource "aws_ec2_transit_gateway_vpc_attachment" "tgw_attachment_%d" {
  transit_gateway_id = aws_ec2_transit_gateway.tgw.id
  vpc_id             = %s
  subnet_ids         = [%s]
  
  tags = {
    Name = "%s"
  }
}
`, i, vpcId, subnetIds, attachmentName))
							
							LogInfo(fmt.Sprintf("已生成中转网关挂载配置 %d: 名称=%s", i+1, attachmentName))
						}
					}
				} else if config.ComponentConfig.EnableVpcAttachment {
					// 使用ComponentConfig中的配置
					tgwConfig := config.ComponentConfig.TransitGatewayConfig
					dnsSupport := "disable"
					ipv6Support := "disable"
					attachmentName := "tgw-attachment"
					
					if tgwConfig.DnsSupport {
						dnsSupport = "enable"
					}
					
					if tgwConfig.Ipv6Support {
						ipv6Support = "enable"
					}
					
					if tgwConfig.AttachmentName != "" {
						attachmentName = tgwConfig.AttachmentName
					}
					
					// 为每个VPC创建一个传输网关挂载
					if config.AllVpcs != nil && len(config.AllVpcs) > 0 {
						for i, vpc := range config.AllVpcs {
							// 确定子网ID
							var subnetIds string
							if config.AllSubnets != nil && len(config.AllSubnets) > 0 {
								// 查找属于当前VPC的子网
								var vpcSubnets []string
								for _, subnet := range config.AllSubnets {
									if subnet.VpcIndex == i {
										vpcSubnets = append(vpcSubnets, fmt.Sprintf("aws_subnet.%s.id", subnet.Name))
									}
								}
								
								if len(vpcSubnets) > 0 {
									subnetIds = strings.Join(vpcSubnets, ", ")
								} else {
									// 如果没有找到属于当前VPC的子网，使用第一个子网
									subnetIds = fmt.Sprintf("aws_subnet.%s.id", config.AllSubnets[0].Name)
								}
							} else if tgwConfig.SubnetIds != "" {
								subnetIds = tgwConfig.SubnetIds
							} else {
								subnetIds = fmt.Sprintf("aws_subnet.%s.id", config.Subnet.Name)
							}
							
							terraformConfig.WriteString(fmt.Sprintf(`resource "aws_ec2_transit_gateway_vpc_attachment" "tgw_attachment_%d" {
  transit_gateway_id = aws_ec2_transit_gateway.tgw.id
  vpc_id             = aws_vpc.%s.id
  subnet_ids         = [%s]
  
  dns_support                                     = "%s"
  ipv6_support                                    = "%s"
  transit_gateway_default_route_table_association = true
  transit_gateway_default_route_table_propagation = true
  
  tags = {
    Name = "%s-%d"
  }
}
`, i, vpc.Name, subnetIds, dnsSupport, ipv6Support, attachmentName, i+1))
							
							LogInfo(fmt.Sprintf("已生成中转网关挂载配置 %d: 名称=%s-%d, VPC=%s", i+1, attachmentName, i+1, vpc.Name))
						}
					} else {
						// 使用子网ID或默认使用当前子网
						subnetIds := tgwConfig.SubnetIds
						if subnetIds == "" {
							subnetIds = fmt.Sprintf("aws_subnet.%s.id", config.Subnet.Name)
						}
						
						terraformConfig.WriteString(fmt.Sprintf(`resource "aws_ec2_transit_gateway_vpc_attachment" "tgw_attachment" {
  transit_gateway_id = aws_ec2_transit_gateway.tgw.id
  vpc_id             = aws_vpc.%s.id
  subnet_ids         = [%s]
  
  dns_support                                     = "%s"
  ipv6_support                                    = "%s"
  transit_gateway_default_route_table_association = true
  transit_gateway_default_route_table_propagation = true
  
  tags = {
    Name = "%s"
  }
}
`, config.VPC.Name, subnetIds, dnsSupport, ipv6Support, attachmentName))
						
						LogInfo(fmt.Sprintf("已生成中转网关挂载配置: 名称=%s, VPC=%s", attachmentName, config.VPC.Name))
					}
				}
			}
			// 其他云提供商的中转网关配置...
			
		case "s3":
			if config.CloudProvider == "aws" {
				// 处理S3存储桶配置
				// 检查ComponentConfig中是否有storageBuckets数组
				var storageBucketsJSON []byte
				if buckets, ok := propsMap["storageBuckets"]; ok && buckets != nil {
					// 记录存储桶配置
					bucketsStr, ok := buckets.(string)
					if ok && bucketsStr != "" {
						LogInfo(fmt.Sprintf("处理存储桶配置: %s", bucketsStr))
						storageBucketsJSON = []byte(bucketsStr)
					}
				}
				
				// 解析存储桶配置
				if len(storageBucketsJSON) > 0 {
					var buckets []map[string]interface{}
					if err := json.Unmarshal(storageBucketsJSON, &buckets); err == nil {
						for i, bucket := range buckets {
							bucketName, _ := bucket["bucketName"].(string)
							if bucketName == "" {
								bucketName = fmt.Sprintf("my-bucket-%d", i+1)
							}
							
							// 获取存储桶策略类型
							policyType, _ := bucket["policyType"].(string)
							
							// 生成存储桶配置
							terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket" "bucket_%d" {
  bucket = "%s"
  
  tags = {
    Name        = "%s"
    Environment = "Dev"
  }
}
`, i, bucketName, bucketName))
							
							// 根据策略类型添加相应的配置
							if policyType == "public-read" {
								terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket_acl" "bucket_%d_acl" {
  bucket = aws_s3_bucket.bucket_%d.id
  acl    = "%s"
}
`, i, i, policyType))
							} else if policyType == "custom" {
								// 自定义策略
								customPolicy, _ := bucket["customPolicy"].(string)
								if customPolicy != "" {
									terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket_policy" "bucket_%d_policy" {
  bucket = aws_s3_bucket.bucket_%d.id
  policy = <<POLICY
%s
POLICY
}
`, i, i, customPolicy))
								}
							}
							
							// 添加生命周期规则
							if enableLifecycle, ok := bucket["enableLifecycleRules"].(bool); ok && enableLifecycle {
								lifecycleRule, ok := bucket["lifecycleRule"].(map[string]interface{})
								if ok {
									ruleName, _ := lifecycleRule["name"].(string)
									ruleStatus, _ := lifecycleRule["status"].(string)
									expirationDays := 0
									transitionDays := 0
									
									if expDays, ok := lifecycleRule["expirationDays"].(float64); ok {
										expirationDays = int(expDays)
									}
									
									if transDays, ok := lifecycleRule["transitionDays"].(float64); ok {
										transitionDays = int(transDays)
									}
									
									terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket_lifecycle_configuration" "bucket_%d_lifecycle" {
  bucket = aws_s3_bucket.bucket_%d.id
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
`, i, i, ruleName, ruleStatus, expirationDays, transitionDays))
								}
							}
						}
					}
				} else if config.ComponentConfig.BucketName != "" {
					// 使用ComponentConfig中的配置
					bucketName := config.ComponentConfig.BucketName
					bucketPolicyType := config.ComponentConfig.BucketPolicyType
					
					terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket" "storage" {
  bucket = "%s"
  
  tags = {
    Name        = "%s"
    Environment = "Dev"
  }
}
`, bucketName, bucketName))
					
					// 根据策略类型添加相应的配置
					if bucketPolicyType == "public-read" {
						terraformConfig.WriteString(`resource "aws_s3_bucket_acl" "storage_acl" {
  bucket = aws_s3_bucket.storage.id
  acl    = "public-read"
}
`)
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
	// 创建拓扑图数据结构
	topology := map[string]interface{}{
		"nodes": []map[string]interface{}{},
		"edges": []map[string]interface{}{},
	}
	
	// 添加VPC节点
	if config.AllVpcs != nil && len(config.AllVpcs) > 0 {
		for i, vpc := range config.AllVpcs {
			vpcNode := map[string]interface{}{
				"id":   vpc.Name,
				"type": "vpc",
				"name": vpc.Name,
				"data": map[string]interface{}{
					"cidr": vpc.CIDR,
				},
			}
			topology["nodes"] = append(topology["nodes"].([]map[string]interface{}), vpcNode)
		}
	} else {
		vpcNode := map[string]interface{}{
			"id":   config.VPC.Name,
			"type": "vpc",
			"name": config.VPC.Name,
			"data": map[string]interface{}{
				"cidr": config.VPC.CIDR,
			},
		}
		topology["nodes"] = append(topology["nodes"].([]map[string]interface{}), vpcNode)
	}
	
	// 添加子网节点
	if config.AllSubnets != nil && len(config.AllSubnets) > 0 {
		for i, subnet := range config.AllSubnets {
			subnetNode := map[string]interface{}{
				"id":   subnet.Name,
				"type": "subnet",
				"name": subnet.Name,
				"data": map[string]interface{}{
					"cidr": subnet.CIDR,
					"az":   subnet.AZ,
				},
			}
			topology["nodes"] = append(topology["nodes"].([]map[string]interface{}), subnetNode)
			
			// 添加子网与VPC的连接
			vpcIndex := subnet.VpcIndex
			var vpcName string
			
			if config.AllVpcs != nil && len(config.AllVpcs) > 0 && vpcIndex >= 0 && vpcIndex < len(config.AllVpcs) {
				vpcName = config.AllVpcs[vpcIndex].Name
			} else {
				vpcName = config.VPC.Name
			}
			
			topology["edges"] = append(topology["edges"].([]map[string]interface{}), map[string]interface{}{
				"source": subnet.Name,
				"target": vpcName,
				"label":  "belongs-to",
			})
		}
	} else {
		subnetNode := map[string]interface{}{
			"id":   config.Subnet.Name,
			"type": "subnet",
			"name": config.Subnet.Name,
			"data": map[string]interface{}{
				"cidr": config.Subnet.CIDR,
				"az":   config.AZ,
			},
		}
		topology["nodes"] = append(topology["nodes"].([]map[string]interface{}), subnetNode)
		
		// 添加子网与VPC的连接
		topology["edges"] = append(topology["edges"].([]map[string]interface{}), map[string]interface{}{
			"source": config.Subnet.Name,
			"target": config.VPC.Name,
			"label":  "belongs-to",
		})
	}
	
	// 添加组件节点
	for _, component := range config.Components {
		// 创建组件节点
		componentNode := map[string]interface{}{
			"id":   component,
			"type": component,
			"name": component,
			"data": map[string]interface{}{
				"description": "组件: " + component,
			},
		}
		
		// 添加组件节点
		topology["nodes"] = append(topology["nodes"].([]map[string]interface{}), componentNode)
		
		// 添加组件与VPC或子网的连接
		if component == "transit-gateway" {
			// 中转网关连接到所有VPC
			if config.AllVpcs != nil && len(config.AllVpcs) > 0 {
				for _, vpc := range config.AllVpcs {
					topology["edges"] = append(topology["edges"].([]map[string]interface{}), map[string]interface{}{
						"source": component,
						"target": vpc.Name,
						"label":  "connected-to",
					})
				}
			} else {
				topology["edges"] = append(topology["edges"].([]map[string]interface{}), map[string]interface{}{
					"source": component,
					"target": config.VPC.Name,
					"label":  "connected-to",
				})
			}
		} else {
			// 其他组件连接到第一个VPC
			vpcName := config.VPC.Name
			if config.AllVpcs != nil && len(config.AllVpcs) > 0 {
				vpcName = config.AllVpcs[0].Name
			}
			
			topology["edges"] = append(topology["edges"].([]map[string]interface{}), map[string]interface{}{
				"source": component,
				"target": vpcName,
				"label":  "connected-to",
			})
		}
	}
	
	return topology
}

// getActualAwsAZ 获取实际的AWS可用区ID
func getActualAwsAZ(region, az string) string {
	// 如果AZ已经包含区域前缀，则直接返回
	if strings.HasPrefix(az, region) {
		return az
	}
	
	// 否则，拼接区域和可用区
	return region + az
}

// LogInfo 记录信息日志
func LogInfo(message string) {
	// 实现日志记录逻辑
}

// LogError 记录错误日志
func LogError(message string) {
	// 实现日志记录逻辑
}
