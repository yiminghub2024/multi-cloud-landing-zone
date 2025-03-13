package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/multi-cloud-landing-zone/backend/models"
)

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
		terraformConfig.WriteString(fmt.Sprintf(`resource "aws_subnet" "%s" {
  vpc_id                  = aws_vpc.%s.id
  cidr_block              = "%s"
  availability_zone       = "%s"
  map_public_ip_on_launch = %t
  
  tags = {
    Name = "%s"
  }
}

`, config.Subnet.Name, config.VPC.Name, config.Subnet.CIDR, config.AZ, config.Subnet.MapPublicIpOnLaunch, config.Subnet.Name))
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
			} else if config.CloudProvider == "azure" {
				terraformConfig.WriteString(fmt.Sprintf(`resource "azurerm_public_ip" "lb_ip" {
  name                = "lb-ip"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "load_balancer" {
  name                = "%s-%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.lb_ip.id
  }
}

resource "azurerm_lb_backend_address_pool" "backend_pool" {
  loadbalancer_id = azurerm_lb.load_balancer.id
  name            = "BackEndAddressPool"
}

resource "azurerm_lb_rule" "lb_rule" {
  loadbalancer_id                = azurerm_lb.load_balancer.id
  name                           = "LBRule"
  protocol                       = "Tcp"
  frontend_port                  = %s
  backend_port                   = %s
  frontend_ip_configuration_name = "PublicIPAddress"
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.backend_pool.id]
}

`, component.Name, config.VPC.Name, props["listener_port"], props["listener_port"]))
			}
		case "object-storage":
			if config.CloudProvider == "aws" {
				terraformConfig.WriteString(fmt.Sprintf(`resource "aws_s3_bucket" "storage" {
  bucket = "%s"
}

resource "aws_s3_bucket_acl" "storage_acl" {
  bucket = aws_s3_bucket.storage.id
  acl    = "private"
}

`, props["bucket_name"]))
			} else if config.CloudProvider == "azure" {
				terraformConfig.WriteString(fmt.Sprintf(`resource "azurerm_storage_account" "storage" {
  name                     = "%s"
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "container" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.storage.name
  container_access_type = "private"
}

`, props["storage_account_name"]))
			}
		case "database":
			if config.CloudProvider == "aws" {
				terraformConfig.WriteString(fmt.Sprintf(`resource "aws_db_instance" "database" {
  allocated_storage    = %s
  engine               = "%s"
  engine_version       = "%s"
  instance_class       = "%s"
  db_name              = "%s"
  username             = "%s"
  password             = "%s"
  parameter_group_name = "default.mysql5.7"
  skip_final_snapshot  = true
  vpc_security_group_ids = []
  db_subnet_group_name = aws_db_subnet_group.default.name
}

resource "aws_db_subnet_group" "default" {
  name       = "main"
  subnet_ids = [aws_subnet.%s.id]
}

`, props["allocated_storage"], props["engine"], props["engine_version"], props["instance_class"], props["db_name"], props["username"], props["password"], config.Subnet.Name))
			} else if config.CloudProvider == "azure" {
				terraformConfig.WriteString(fmt.Sprintf(`resource "azurerm_mysql_server" "database" {
  name                = "%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  administrator_login          = "%s"
  administrator_login_password = "%s"

  sku_name   = "%s"
  storage_mb = %s
  version    = "%s"

  auto_grow_enabled                 = true
  backup_retention_days             = 7
  geo_redundant_backup_enabled      = false
  infrastructure_encryption_enabled = false
  public_network_access_enabled     = true
  ssl_enforcement_enabled           = true
}

resource "azurerm_mysql_database" "database" {
  name                = "%s"
  resource_group_name = azurerm_resource_group.rg.name
  server_name         = azurerm_mysql_server.database.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}

`, props["server_name"], props["username"], props["password"], props["sku_name"], props["storage_mb"], props["version"], props["db_name"]))
			}
		}
	}

	LogInfo("Terraform配置生成完成")
	return terraformConfig.String()
}

// SaveTerraformConfig 保存Terraform配置到文件
func SaveTerraformConfig(config string, filePath string) error {
	LogInfo(fmt.Sprintf("正在保存Terraform配置到文件: %s", filePath))
	
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		LogError(fmt.Sprintf("创建目录失败: %s, 错误: %v", dir, err))
		return fmt.Errorf("创建目录失败: %w", err)
	}
	
	// 写入文件
	if err := os.WriteFile(filePath, []byte(config), 0644); err != nil {
		LogError(fmt.Sprintf("写入Terraform配置文件失败: %s, 错误: %v", filePath, err))
		return fmt.Errorf("写入Terraform配置文件失败: %w", err)
	}
	
	LogInfo(fmt.Sprintf("Terraform配置已成功保存到: %s", filePath))
	
	// 记录文件内容
	LogDebug(fmt.Sprintf("Terraform配置文件内容:\n%s", config))
	
	return nil
}

// GenerateTopology 生成资源拓扑图
func GenerateTopology(config models.DeploymentConfig) models.Topology {
	LogInfo("开始生成资源拓扑图")
	
	var topology models.Topology
	
	// 添加VPC节点
	vpcNode := models.TopologyNode{
		ID:   "vpc-" + config.VPC.Name,
		Type: "vpc",
		Name: config.VPC.Name,
		Data: map[string]interface{}{
			"cidr": config.VPC.CIDR,
		},
	}
	topology.Nodes = append(topology.Nodes, vpcNode)
	
	// 添加子网节点
	subnetNode := models.TopologyNode{
		ID:   "subnet-" + config.Subnet.Name,
		Type: "subnet",
		Name: config.Subnet.Name,
		Data: map[string]interface{}{
			"cidr": config.Subnet.CIDR,
			"az":   config.AZ,
		},
	}
	topology.Nodes = append(topology.Nodes, subnetNode)
	
	// 添加VPC到子网的连接
	topology.Edges = append(topology.Edges, models.TopologyEdge{
		Source: "vpc-" + config.VPC.Name,
		Target: "subnet-" + config.Subnet.Name,
		Label:  "contains",
	})
	
	// 添加组件节点
	for _, component := range config.Components {
		componentNode := models.TopologyNode{
			ID:   component.Value + "-" + component.Name,
			Type: component.Value,
			Name: component.Name,
			Data: map[string]interface{}{
				"properties": config.ComponentProperties[component.Value],
			},
		}
		topology.Nodes = append(topology.Nodes, componentNode)
		
		// 添加子网到组件的连接
		topology.Edges = append(topology.Edges, models.TopologyEdge{
			Source: "subnet-" + config.Subnet.Name,
			Target: component.Value + "-" + component.Name,
			Label:  "hosts",
		})
	}
	
	LogInfo(fmt.Sprintf("资源拓扑图生成完成，包含 %d 个节点和 %d 个连接", len(topology.Nodes), len(topology.Edges)))
	return topology
}
