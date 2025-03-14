<template>
  <div class="deployment-summary-container">
    <el-card shadow="hover">
      <template #header>
        <h2>部署摘要</h2>
      </template>
      <div class="summary-content">
        <el-alert
          v-if="!isDataComplete"
          title="请完成所有前置步骤后查看部署摘要"
          type="warning"
          :closable="false"
          show-icon
        />
        
        <div v-else>
          <el-card class="mb-4">
            <template #header>
              <h3>基础信息</h3>
            </template>
            <el-descriptions :column="1" border>
              <el-descriptions-item label="云服务提供商">{{ getCloudProviderName(cloudProvider) }}</el-descriptions-item>
              <el-descriptions-item label="区域">{{ regionName }}</el-descriptions-item>
              <el-descriptions-item label="可用区">{{ azName }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
          
          <el-card class="mb-4">
            <template #header>
              <h3>网络配置</h3>
            </template>
            <el-descriptions :column="1" border>
              <el-descriptions-item label="VPC名称">{{ vpcConfig.name }}</el-descriptions-item>
              <el-descriptions-item label="VPC CIDR">{{ vpcConfig.cidr }}</el-descriptions-item>
              <el-descriptions-item label="DNS支持">
                <el-tag :type="vpcConfig.enableDnsSupport ? 'success' : 'info'">
                  {{ vpcConfig.enableDnsSupport ? '启用' : '禁用' }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="DNS主机名">
                <el-tag :type="vpcConfig.enableDnsHostnames ? 'success' : 'info'">
                  {{ vpcConfig.enableDnsHostnames ? '启用' : '禁用' }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
            
            <el-divider content-position="left">子网配置</el-divider>
            
            <el-descriptions :column="1" border>
              <el-descriptions-item label="子网名称">{{ subnetConfig.name }}</el-descriptions-item>
              <el-descriptions-item label="子网CIDR">{{ subnetConfig.cidr }}</el-descriptions-item>
              <el-descriptions-item label="可用区">{{ azName }}</el-descriptions-item>
              <el-descriptions-item label="自动分配公网IP">
                <el-tag :type="subnetConfig.mapPublicIpOnLaunch ? 'success' : 'info'">
                  {{ subnetConfig.mapPublicIpOnLaunch ? '是' : '否' }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </el-card>
          
          <el-card class="mb-4">
            <template #header>
              <h3>云组件</h3>
            </template>
            <el-empty 
              v-if="selectedComponents.length === 0" 
              description="未选择任何云组件" 
            />
            <div v-else>
              <el-row :gutter="20">
                <el-col 
                  v-for="(component, index) in selectedComponentsDetails" 
                  :key="index"
                  :xs="24" 
                  :sm="24" 
                  :md="12" 
                  :lg="8"
                  class="mb-4"
                >
                  <el-card shadow="hover">
                    <template #header>
                      <h4>{{ component.name }}</h4>
                    </template>
                    <p>{{ component.description }}</p>
                    
                    <div v-if="component.properties && component.properties.length > 0">
                      <el-divider content-position="left">配置参数</el-divider>
                      <el-descriptions :column="1" border size="small">
                        <el-descriptions-item 
                          v-for="(prop, propIndex) in component.properties" 
                          :key="propIndex"
                          :label="prop.name"
                        >
                          {{ getComponentPropertyValue(component.value, prop.key) || '未设置' }}
                        </el-descriptions-item>
                      </el-descriptions>
                    </div>
                  </el-card>
                </el-col>
              </el-row>
            </div>
          </el-card>
          
          <el-card class="mb-4">
            <template #header>
              <h3>Terraform配置预览</h3>
            </template>
            <el-scrollbar height="400px">
              <pre class="terraform-code">{{ generateTerraformPreview() }}</pre>
            </el-scrollbar>
          </el-card>
          
          <div class="text-center mt-4">
            <el-button 
              type="success" 
              size="large" 
              @click="startDeployment"
              :disabled="!isDataComplete"
              :loading="deploying"
            >
              开始部署
            </el-button>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
// 导入API配置
import { deployInfrastructure } from '../api';

export default {
  name: 'DeploymentSummary',
  props: {
    cloudProvider: {
      type: String,
      default: ''
    },
    regionName: {
      type: String,
      default: ''
    },
    azName: {
      type: String,
      default: ''
    },
    vpcConfig: {
      type: Object,
      default: () => ({})
    },
    subnetConfig: {
      type: Object,
      default: () => ({})
    },
    selectedComponents: {
      type: Array,
      default: () => []
    },
    componentProperties: {
      type: Object,
      default: () => ({})
    },
    allComponents: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      deploying: false,
      deploymentStatus: null,
      deploymentError: null
    }
  },
  computed: {
    isDataComplete() {
      return (
        this.cloudProvider && 
        this.regionName && 
        this.azName && 
        this.vpcConfig.name && 
        this.vpcConfig.cidr && 
        this.subnetConfig.name && 
        this.subnetConfig.cidr
      );
    },
    selectedComponentsDetails() {
      return this.allComponents.filter(component => 
        this.selectedComponents.includes(component.value)
      );
    }
  },
  methods: {
    getCloudProviderName(provider) {
      const providerMap = {
        'aws': 'AWS',
        'azure': 'Azure',
        'alicloud': '阿里云',
        'baidu': '百度云',
        'huawei': '华为云',
        'tencent': '腾讯云',
        'volcengine': '火山云'
      };
      
      return providerMap[provider] || provider;
    },
    getComponentPropertyValue(componentValue, propertyKey) {
      if (
        this.componentProperties && 
        this.componentProperties[componentValue] && 
        this.componentProperties[componentValue][propertyKey] !== undefined
      ) {
        return this.componentProperties[componentValue][propertyKey];
      }
      
      // 查找默认值
      const component = this.allComponents.find(c => c.value === componentValue);
      if (component && component.properties) {
        const property = component.properties.find(p => p.key === propertyKey);
        if (property) {
          return property.defaultValue || '';
        }
      }
      
      return '';
    },
    generateTerraformPreview() {
      // 根据用户选择生成Terraform配置预览
      let preview = '';
      
      // 添加提供商配置
      switch(this.cloudProvider) {
        case 'aws':
          preview += `provider "aws" {
  region = "${this.regionName}"
}

`;
          break;
        case 'azure':
          preview += `provider "azurerm" {
  features {}
}

`;
          break;
        case 'alicloud':
          preview += `provider "alicloud" {
  region = "${this.regionName}"
}

`;
          break;
        // 其他云服务提供商...
        default:
          preview += `provider "${this.cloudProvider}" {
  region = "${this.regionName}"
}

`;
      }
      
      // 添加VPC配置
      if (this.cloudProvider === 'aws') {
        preview += `resource "aws_vpc" "${this.vpcConfig.name}" {
  cidr_block           = "${this.vpcConfig.cidr}"
  enable_dns_support   = ${this.vpcConfig.enableDnsSupport || true}
  enable_dns_hostnames = ${this.vpcConfig.enableDnsHostnames || true}
  
  tags = {
    Name = "${this.vpcConfig.name}"
  }
}

`;
      } else if (this.cloudProvider === 'azure') {
        preview += `resource "azurerm_virtual_network" "${this.vpcConfig.name}" {
  name                = "${this.vpcConfig.name}"
  address_space       = ["${this.vpcConfig.cidr}"]
  location            = "${this.regionName}"
  resource_group_name = "resource-group-name"
}

`;
      } else {
        preview += `resource "${this.cloudProvider}_vpc" "${this.vpcConfig.name}" {
  name       = "${this.vpcConfig.name}"
  cidr_block = "${this.vpcConfig.cidr}"
}

`;
      }
      
      // 添加子网配置
      if (this.cloudProvider === 'aws') {
        preview += `resource "aws_subnet" "${this.subnetConfig.name}" {
  vpc_id                  = aws_vpc.${this.vpcConfig.name}.id
  cidr_block              = "${this.subnetConfig.cidr}"
  availability_zone       = "${this.subnetConfig.az || 'us-east-1a'}"
  map_public_ip_on_launch = ${this.subnetConfig.mapPublicIpOnLaunch || false}
  
  tags = {
    Name = "${this.subnetConfig.name}"
  }
}

`;
      } else if (this.cloudProvider === 'azure') {
        preview += `resource "azurerm_subnet" "${this.subnetConfig.name}" {
  name                 = "${this.subnetConfig.name}"
  resource_group_name  = "resource-group-name"
  virtual_network_name = azurerm_virtual_network.${this.vpcConfig.name}.name
  address_prefixes     = ["${this.subnetConfig.cidr}"]
}

`;
      } else {
        preview += `resource "${this.cloudProvider}_subnet" "${this.subnetConfig.name}" {
  vpc_id     = ${this.cloudProvider}_vpc.${this.vpcConfig.name}.id
  cidr_block = "${this.subnetConfig.cidr}"
  zone_id    = "${this.subnetConfig.az || 'zone-1'}"
}

`;
      }
      
      // 添加选定的组件
      this.selectedComponents.forEach(componentValue => {
        const component = this.allComponents.find(c => c.value === componentValue);
        if (!component) return;
        
        switch(component.value) {
          case 'load-balancer':
            if (this.cloudProvider === 'aws') {
              preview += `resource "aws_lb" "load_balancer" {
  name               = "${this.vpcConfig.name}-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.lb_sg.id]
  subnets            = [aws_subnet.${this.subnetConfig.name}.id]
}

resource "aws_security_group" "lb_sg" {
  name        = "lb-sg"
  description = "Allow inbound traffic"
  vpc_id      = aws_vpc.${this.vpcConfig.name}.id

  ingress {
    from_port   = ${this.getComponentPropertyValue(component.value, 'listener_port') || '80'}
    to_port     = ${this.getComponentPropertyValue(component.value, 'listener_port') || '80'}
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

`;
            }
            break;
          
          case 'object-storage':
            if (this.cloudProvider === 'aws') {
              const bucketName = this.getComponentPropertyValue(component.value, 'bucket_name') || `${this.vpcConfig.name}-bucket`;
              preview += `resource "aws_s3_bucket" "storage_bucket" {
  bucket = "${bucketName}"
  
  tags = {
    Name = "${bucketName}"
  }
}
`;
            }
            break;
          
          case 'rds':
            if (this.cloudProvider === 'aws') {
              preview += `resource "aws_db_instance" "database" {
  allocated_storage    = ${this.getComponentPropertyValue(component.value, 'storage_size') || '20'}
  engine               = "${this.getComponentPropertyValue(component.value, 'engine') || 'mysql'}"
  engine_version       = "5.7"
  instance_class       = "db.t3.${this.getComponentPropertyValue(component.value, 'instance_type') || 'small'}"
  name                 = "mydb"
  username             = "admin"
  password             = "password"
  parameter_group_name = "default.mysql5.7"
  skip_final_snapshot  = true
  db_subnet_group_name = aws_db_subnet_group.default.name
}
resource "aws_db_subnet_group" "default" {
  name       = "main"
  subnet_ids = [aws_subnet.${this.subnetConfig.name}.id]
}
`;
            }
            break;
          
          case 'compute':
            if (this.cloudProvider === 'aws') {
              preview += `resource "aws_instance" "web" {
  count         = ${this.getComponentPropertyValue(component.value, 'instance_count') || '2'}
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.${this.getComponentPropertyValue(component.value, 'instance_type') || 'medium'}"
  subnet_id     = aws_subnet.${this.subnetConfig.name}.id
  
  tags = {
    Name = "WebServer-\${count.index + 1}"
  }
}
`;
            }
            break;
          
          // 其他组件...
        }
      });
      
      return preview;
    },
    async startDeployment() {
      this.deploying = true;
      this.deploymentStatus = null;
      this.deploymentError = null;
      
      // 准备部署数据
      const deploymentData = {
        cloudProvider: this.cloudProvider,
        region: this.regionName,
        az: this.azName,
        vpc: this.vpcConfig,
        subnet: this.subnetConfig,
        components: this.selectedComponentsDetails,
        componentProperties: this.componentProperties,
        terraformPreview: this.generateTerraformPreview()
      };
      
      try {
        // 调用后端API
        const response = await deployInfrastructure(deploymentData);
        
        // 处理成功响应
        this.deploymentStatus = 'success';
        console.log('部署请求成功:', response);
        
        // 发出部署事件
        this.$emit('start-deployment', {
          ...deploymentData,
          status: 'success',
          response
        });
        
        // 可以在这里添加导航到部署反馈页面的逻辑
        // this.$router.push('/deployment-feedback');
      } catch (error) {
        // 处理错误
        this.deploymentStatus = 'error';
        this.deploymentError = error.message || '部署请求失败';
        console.error('部署请求失败:', error);
        
        // 发出部署失败事件
        this.$emit('deployment-error', {
          error: this.deploymentError,
          data: deploymentData
        });
      } finally {
        this.deploying = false;
      }
    }
  }
}
</script>

<style scoped>
.deployment-summary-container {
  margin-bottom: 2rem;
}

.summary-content {
  margin-top: 1rem;
}

.mb-4 {
  margin-bottom: 1rem;
}

.mt-4 {
  margin-top: 1rem;
}

.text-center {
  text-align: center;
}

.terraform-code {
  margin: 0;
  color: #f8f9fa;
  font-family: monospace;
  white-space: pre-wrap;
  background-color: #2c3e50;
  padding: 1rem;
  border-radius: 4px;
}
</style>
