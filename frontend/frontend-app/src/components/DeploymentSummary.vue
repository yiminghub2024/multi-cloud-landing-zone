<template>
  <div class="deployment-summary-container">
    <h2>部署摘要</h2>
    <div class="summary-content">
      <div v-if="!isDataComplete" class="incomplete-data">
        <p>请完成所有前置步骤后查看部署摘要</p>
      </div>
      <div v-else>
        <div class="summary-section">
          <h3>基础信息</h3>
          <div class="summary-item">
            <div class="summary-label">云服务提供商:</div>
            <div class="summary-value">{{ getCloudProviderName(cloudProvider) }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">区域:</div>
            <div class="summary-value">{{ regionName }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">可用区:</div>
            <div class="summary-value">{{ azName }}</div>
          </div>
        </div>
        
        <div class="summary-section">
          <h3>网络配置</h3>
          <div class="summary-item">
            <div class="summary-label">VPC名称:</div>
            <div class="summary-value">{{ vpcConfig.name }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">VPC CIDR:</div>
            <div class="summary-value">{{ vpcConfig.cidr }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">DNS支持:</div>
            <div class="summary-value">{{ vpcConfig.enableDnsSupport ? '启用' : '禁用' }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">DNS主机名:</div>
            <div class="summary-value">{{ vpcConfig.enableDnsHostnames ? '启用' : '禁用' }}</div>
          </div>
          
          <div class="summary-subsection">
            <h4>子网配置</h4>
            <div class="summary-item">
              <div class="summary-label">子网名称:</div>
              <div class="summary-value">{{ subnetConfig.name }}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">子网CIDR:</div>
              <div class="summary-value">{{ subnetConfig.cidr }}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">可用区:</div>
              <div class="summary-value">{{ azName }}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">自动分配公网IP:</div>
              <div class="summary-value">{{ subnetConfig.mapPublicIpOnLaunch ? '是' : '否' }}</div>
            </div>
          </div>
        </div>
        
        <div class="summary-section">
          <h3>云组件</h3>
          <div v-if="selectedComponents.length === 0" class="no-components">
            <p>未选择任何云组件</p>
          </div>
          <div v-else class="component-summary">
            <div 
              v-for="(component, index) in selectedComponentsDetails" 
              :key="index"
              class="component-summary-item"
            >
              <div class="component-summary-header">
                <h4>{{ component.name }}</h4>
              </div>
              <div class="component-summary-desc">{{ component.description }}</div>
              
              <div v-if="component.properties && component.properties.length > 0" class="component-properties-summary">
                <h5>配置参数</h5>
                <div 
                  v-for="(prop, propIndex) in component.properties" 
                  :key="propIndex"
                  class="property-summary-item"
                >
                  <div class="property-summary-label">{{ prop.name }}:</div>
                  <div class="property-summary-value">
                    {{ getComponentPropertyValue(component.value, prop.key) || '未设置' }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="terraform-preview">
          <h3>Terraform配置预览</h3>
          <div class="terraform-code">
            <pre>{{ generateTerraformPreview() }}</pre>
          </div>
        </div>
        
        <div class="deployment-actions">
          <button 
            class="deploy-button" 
            @click="startDeployment"
            :disabled="!isDataComplete"
          >
            开始部署
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
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
  region = "${this.vpcConfig.region || 'us-east-1'}"
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
  region = "${this.vpcConfig.region || 'cn-hangzhou'}"
}

`;
          break;
        // 其他云服务提供商...
        default:
          preview += `provider "${this.cloudProvider}" {
  region = "${this.vpcConfig.region || 'default-region'}"
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
  location            = "${this.vpcConfig.region || 'eastus'}"
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
      this.selectedComponentsDetails.forEach(component => {
        switch(component.value) {
          case 'load-balancer':
            if (this.cloudProvider === 'aws') {
              preview += `resource "aws_lb" "load_balancer" {
  name               = "${component.name}-${this.vpcConfig.name}"
  internal           = false
  load_balancer_type = "application"
  subnets            = [aws_subnet.${this.subnetConfig.name}.id]
  
  enable_deletion_protection = false
}

resource "aws_lb_listener" "front_end" {
  load_balancer_arn = aws_lb.load_balancer.arn
  port              = "${this.getComponentPropertyValue(component.value, 'listener_port') || '80'}"
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
  vpc_id   = aws_vpc.${this.vpcConfig.name}.id
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
    startDeployment() {
      // 发出部署事件
      this.$emit('start-deployment', {
        cloudProvider: this.cloudProvider,
        region: this.regionName,
        az: this.azName,
        vpc: this.vpcConfig,
        subnet: this.subnetConfig,
        components: this.selectedComponentsDetails,
        componentProperties: this.componentProperties,
        terraformPreview: this.generateTerraformPreview()
      });
    }
  }
}
</script>

<style scoped>
.deployment-summary-container {
  margin-bottom: 2rem;
}

.summary-content {
  margin-top: 1.5rem;
}

.incomplete-data {
  padding: 2rem;
  text-align: center;
  color: #e74c3c;
  background-color: #fdf2f2;
  border-radius: 8px;
}

.summary-section {
  margin-bottom: 2rem;
  padding: 1.5rem;
  background-color: #f8f9fa;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.summary-section h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #e9ecef;
  color: #2c3e50;
}

.summary-subsection {
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px dashed #e9ecef;
}

.summary-subsection h4 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #2c3e50;
}

.summary-item {
  display: flex;
  margin-bottom: 0.75rem;
}

.summary-label {
  width: 180px;
  font-weight: bold;
  color: #6c757d;
}

.summary-value {
  flex: 1;
  font-family: monospace;
}

.no-components {
  padding: 1.5rem;
  text-align: center;
  color: #6c757d;
  background-color: #f1f3f5;
  border-radius: 4px;
}

.component-summary {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.component-summary-item {
  padding: 1.5rem;
  background-color: #fff;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.component-summary-header h4 {
  margin-top: 0;
  margin-bottom: 0.75rem;
  color: #3498db;
}

.component-summary-desc {
  margin-bottom: 1.5rem;
  color: #6c757d;
  font-size: 0.9rem;
}

.component-properties-summary h5 {
  margin-top: 0;
  margin-bottom: 1rem;
  padding-top: 1rem;
  border-top: 1px dashed #e9ecef;
  color: #2c3e50;
  font-size: 0.95rem;
}

.property-summary-item {
  display: flex;
  margin-bottom: 0.5rem;
}

.property-summary-label {
  width: 120px;
  font-weight: bold;
  color: #6c757d;
  font-size: 0.9rem;
}

.property-summary-value {
  flex: 1;
  font-family: monospace;
  font-size: 0.9rem;
}

.terraform-preview {
  margin-top: 2rem;
  margin-bottom: 2rem;
}

.terraform-preview h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #2c3e50;
}

.terraform-code {
  padding: 1.5rem;
  background-color: #2c3e50;
  border-radius: 8px;
  overflow: auto;
  max-height: 400px;
}

.terraform-code pre {
  margin: 0;
  color: #f8f9fa;
  font-family: monospace;
  white-space: pre-wrap;
}

.deployment-actions {
  margin-top: 2rem;
  text-align: center;
}

.deploy-button {
  padding: 1rem 2rem;
  background-color: #2ecc71;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1.1rem;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.deploy-button:hover {
  background-color: #27ae60;
}

.deploy-button:disabled {
  background-color: #95a5a6;
  cursor: not-allowed;
}
</style>
