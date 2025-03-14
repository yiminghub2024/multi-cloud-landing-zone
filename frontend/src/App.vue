<template>
  <div id="app">
    <header class="header">
      <h1>多云着陆区部署平台</h1>
    </header>
    <main class="main-content">
      <div class="step-container">
        <div class="step-indicator">
          <div class="step" :class="{ active: currentStep >= 1, completed: currentStep > 1 }">1. 选择云服务提供商</div>
          <div class="step" :class="{ active: currentStep >= 2, completed: currentStep > 2 }">2. 选择区域(Region)</div>
          <div class="step" :class="{ active: currentStep >= 3, completed: currentStep > 3 }">3. 选择可用区(AZ)</div>
          <div class="step" :class="{ active: currentStep >= 4, completed: currentStep > 4 }">4. 配置VPC</div>
          <div class="step" :class="{ active: currentStep >= 5, completed: currentStep > 5 }">5. 配置子网</div>
          <div class="step" :class="{ active: currentStep >= 6, completed: currentStep > 6 }">6. 选择云组件</div>
          <div class="step" :class="{ active: currentStep >= 7, completed: currentStep > 7 }">7. 云组件配置</div>
          <div class="step" :class="{ active: currentStep >= 8, completed: currentStep > 8 }">8. 部署摘要</div>
        </div>
        
        <div class="step-content">
          <cloud-provider 
            v-if="currentStep === 1"
            @cloud-selected="handleCloudSelected"
            @regions-loaded="handleRegionsLoaded"
          />
          
          <region-selector 
            v-if="currentStep === 2"
            :cloud-provider="selectedCloud"
            :initial-regions="regions"
            @region-selected="handleRegionSelected"
            @azs-loaded="handleAZsLoaded"
          />
          
          <a-z-selector 
            v-if="currentStep === 3"
            :cloud-provider="selectedCloud"
            :region="selectedRegion"
            :initial-a-zs="availabilityZones"
            @az-selected="handleAZSelected"
          />
          
          <v-p-c-config 
            v-if="currentStep === 4"
            :cloud-provider="selectedCloud"
            :region="selectedRegion"
            @vpc-updated="handleVPCUpdated"
          />
          
          <subnet-config 
            v-if="currentStep === 5"
            :cloud-provider="selectedCloud"
            :region="selectedRegion"
            :az="selectedAZ"
            :az-name="getAZName(selectedAZ)"
            :vpc-name="vpcConfig.name"
            :vpc-c-i-d-r="vpcConfig.cidr"
            @subnet-updated="handleSubnetUpdated"
          />
          
          <component-selector 
            v-if="currentStep === 6"
            :cloud-provider="selectedCloud"
            :region="selectedRegion"
            :az="selectedAZ"
            :initial-selected-components="selectedComponents"
            @components-updated="handleComponentsUpdated"
          />
          
          <cloud-component-config
            v-if="currentStep === 7"
            :cloud-provider="selectedCloud"
            :region="selectedRegion"
            :selected-components="selectedComponents"
            :component-properties="componentProperties"
            :initial-config="componentConfig"
            @config-updated="handleComponentConfigUpdated"
          />
          
          <deployment-summary 
            v-if="currentStep === 8 && !isDeploying && !isDeploymentComplete"
            :cloud-provider="selectedCloud"
            :region-name="getRegionName(selectedRegion)"
            :az-name="getAZName(selectedAZ)"
            :vpc-config="vpcConfig"
            :subnet-config="subnetConfig"
            :selected-components="selectedComponents"
            :component-properties="componentProperties"
            :component-config="componentConfig"
            :all-components="cloudComponents"
            @start-deployment="startDeployment"
          />
          
          <deployment-feedback 
            v-if="isDeploying || isDeploymentComplete"
            :status="deploymentStatus"
            :progress="deploymentProgress"
            :message="deploymentMessage"
            :logs="deploymentLogs"
            :result="deploymentResult"
            :topology="deploymentTopology"
            @refresh-status="refreshDeploymentStatus"
            @retry-deployment="retryDeployment"
            @back-to-summary="backToSummary"
          />
        </div>
        
        <div class="step-navigation">
          <el-button 
            v-if="currentStep > 1 && !isDeploying && !isDeploymentComplete" 
            @click="prevStep"
            icon="el-icon-arrow-left"
          >
            上一步
          </el-button>
          <el-button 
            v-if="currentStep < 7 && !isDeploying && !isDeploymentComplete && isCurrentStepValid" 
            type="primary"
            @click="nextStep"
            icon="el-icon-arrow-right"
          >
            下一步
          </el-button>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import CloudProvider from './components/CloudProvider.vue';
import RegionSelector from './components/RegionSelector.vue';
import AZSelector from './components/AZSelector.vue';
import VPCConfig from './components/VPCConfig.vue';
import SubnetConfig from './components/SubnetConfig.vue';
import ComponentSelector from './components/ComponentSelector.vue';
import CloudComponentConfig from './components/CloudComponentConfig.vue';
import DeploymentSummary from './components/DeploymentSummary.vue';
import DeploymentFeedback from './components/DeploymentFeedback.vue';

export default {
  name: 'App',
  components: {
    CloudProvider,
    RegionSelector,
    AZSelector,
    VPCConfig,
    SubnetConfig,
    ComponentSelector,
    CloudComponentConfig,
    DeploymentSummary,
    DeploymentFeedback
  },
  data() {
    return {
      currentStep: 1,
      selectedCloud: '',
      regions: [],
      selectedRegion: '',
      availabilityZones: [],
      selectedAZ: '',
      vpcConfig: {
        name: '',
        cidr: '',
        enableDnsSupport: true,
        enableDnsHostnames: true
      },
      subnetConfig: {
        name: '',
        cidr: '',
        az: '',
        mapPublicIpOnLaunch: true
      },
      cloudComponents: [],
      selectedComponents: [],
      componentProperties: {},
      componentConfig: {
        bucketPolicyType: 'private',
        customBucketPolicy: '',
        enableLifecycleRules: false,
        lifecycleRule: {
          name: 'default-lifecycle-rule',
          status: 'Enabled',
          expirationDays: 365,
          transitionDays: 30
        },
        enableRouteTables: false,
        enableVpcAttachment: false,
        transitGatewayConfig: {
          routeTableName: 'main-route-table',
          defaultRouteTable: true,
          subnetIds: '',
          dnsSupport: true,
          ipv6Support: false
        }
      },
      
      // 部署相关
      isDeploying: false,
      isDeploymentComplete: false,
      deploymentStatus: 'idle', // idle, preparing, deploying, completed, failed
      deploymentProgress: 0,
      deploymentMessage: '',
      deploymentLogs: [],
      deploymentResult: {},
      deploymentTopology: {
        nodes: [],
        edges: []
      }
    }
  },
  computed: {
    isCurrentStepValid() {
      switch(this.currentStep) {
        case 1:
          return !!this.selectedCloud;
        case 2:
          return !!this.selectedRegion;
        case 3:
          return !!this.selectedAZ;
        case 4:
          return !!this.vpcConfig.name && !!this.vpcConfig.cidr;
        case 5:
          return !!this.subnetConfig.name && !!this.subnetConfig.cidr;
        case 6:
          return true; // 组件选择是可选的
        case 7:
          return true;
        default:
          return false;
      }
    },
    selectedComponentsDetails() {
      return this.cloudComponents.filter(component => 
        this.selectedComponents.includes(component.value)
      );
    }
  },
  methods: {
    // 导航方法
    nextStep() {
      if (this.currentStep < 7 && this.isCurrentStepValid) {
        this.currentStep++;
      }
    },
    prevStep() {
      if (this.currentStep > 1) {
        this.currentStep--;
      }
    },
    
    // 事件处理方法
    handleCloudSelected(cloud) {
      this.selectedCloud = cloud;
      this.selectedRegion = '';
      this.selectedAZ = '';
      this.regions = [];
      this.availabilityZones = [];
    },
    handleRegionsLoaded(regions) {
      this.regions = regions;
    },
    handleRegionSelected(region) {
      this.selectedRegion = region;
      this.selectedAZ = '';
      this.availabilityZones = [];
    },
    handleAZsLoaded(azs) {
      this.availabilityZones = azs;
    },
    handleAZSelected(az) {
      this.selectedAZ = az;
    },
    handleVPCUpdated(vpc) {
      this.vpcConfig = vpc;
    },
    handleSubnetUpdated(subnet) {
      this.subnetConfig = subnet;
    },
    handleComponentsUpdated(data) {
      this.selectedComponents = data.selectedComponents;
      this.componentProperties = data.componentProperties;
    },
    handleComponentConfigUpdated(config) {
      this.componentConfig = config;
    },
    
    // 辅助方法
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
    getRegionName(regionValue) {
      const region = this.regions.find(r => r.value === regionValue);
      return region ? region.name : regionValue;
    },
    getAZName(azValue) {
      const az = this.availabilityZones.find(a => a.value === azValue);
      return az ? az.name : azValue;
    },
    
    // 部署方法
    startDeployment(deploymentData) {
      this.isDeploying = true;
      this.deploymentStatus = 'preparing';
      this.deploymentProgress = 0;
      this.deploymentMessage = '正在准备部署资源...';
      this.deploymentLogs = ['[INFO] 开始部署准备...'];
      
      // 模拟部署过程
      this.simulateDeployment(deploymentData);
    },
    simulateDeployment(deploymentData) {
      // 在实际应用中，这里会调用后端API进行真实部署
      // 这里仅作为演示，模拟部署过程
      
      let progress = 0;
      const interval = setInterval(() => {
        progress += 5;
        this.deploymentProgress = progress;
        
        if (progress === 10) {
          this.deploymentStatus = 'deploying';
          this.deploymentMessage = '正在部署VPC...';
          this.deploymentLogs.push(`[INFO] 开始创建VPC: ${deploymentData.vpc.name}`);
        } else if (progress === 30) {
          this.deploymentMessage = '正在部署子网...';
          this.deploymentLogs.push(`[INFO] VPC创建成功: ${deploymentData.vpc.name}`);
          this.deploymentLogs.push(`[INFO] 开始创建子网: ${deploymentData.subnet.name}`);
        } else if (progress === 50) {
          this.deploymentMessage = '正在部署云组件...';
          this.deploymentLogs.push(`[INFO] 子网创建成功: ${deploymentData.subnet.name}`);
          if (deploymentData.components.length > 0) {
            this.deploymentLogs.push(`[INFO] 开始创建云组件...`);
            deploymentData.components.forEach(component => {
              this.deploymentLogs.push(`[INFO] 创建组件: ${component.name}`);
            });
          }
        } else if (progress === 80) {
          this.deploymentMessage = '正在配置网络连接...';
          this.deploymentLogs.push(`[INFO] 云组件创建成功`);
          this.deploymentLogs.push(`[INFO] 配置网络连接...`);
        } else if (progress >= 100) {
          clearInterval(interval);
          this.deploymentStatus = 'completed';
          this.deploymentProgress = 100;
          this.deploymentMessage = '部署成功完成！';
          this.deploymentLogs.push(`[INFO] 部署成功完成`);
          this.isDeploymentComplete = true;
          
          // 设置部署结果
          this.deploymentResult = {
            deploymentId: 'dep-' + Math.random().toString(36).substr(2, 9),
            cloudProvider: deploymentData.cloudProvider,
            region: deploymentData.region,
            az: deploymentData.az,
            vpc: deploymentData.vpc,
            subnet: deploymentData.subnet,
            components: deploymentData.components
          };
          
          // 生成拓扑图数据
          this.generateTopologyData(deploymentData);
        }
      }, 500);
    },
    generateTopologyData(deploymentData) {
      // 生成拓扑图数据
      const nodes = [];
      const edges = [];
      
      // 添加云节点
      nodes.push({
        id: 'cloud',
        name: this.getCloudProviderName(deploymentData.cloudProvider),
        type: 'cloud'
      });
      
      // 添加VPC节点
      nodes.push({
        id: 'vpc',
        name: deploymentData.vpc.name,
        type: 'vpc'
      });
      
      // 添加VPC与云的连接
      edges.push({
        source: 'cloud',
        target: 'vpc'
      });
      
      // 添加子网节点
      nodes.push({
        id: 'subnet',
        name: deploymentData.subnet.name,
        type: 'subnet'
      });
      
      // 添加子网与VPC的连接
      edges.push({
        source: 'vpc',
        target: 'subnet'
      });
      
      // 添加组件节点
      deploymentData.components.forEach(component => {
        const nodeId = component.value;
        nodes.push({
          id: nodeId,
          name: component.name,
          type: component.value
        });
        
        // 添加组件与子网的连接
        edges.push({
          source: 'subnet',
          target: nodeId
        });
      });
      
      this.deploymentTopology = {
        nodes,
        edges
      };
    },
    refreshDeploymentStatus() {
      // 在实际应用中，这里会调用后端API获取最新的部署状态
      console.log('Refreshing deployment status...');
    },
    retryDeployment() {
      // 重试部署
      this.isDeploying = true;
      this.isDeploymentComplete = false;
      this.deploymentStatus = 'preparing';
      this.deploymentProgress = 0;
      this.deploymentMessage = '正在准备重新部署...';
      this.deploymentLogs = ['[INFO] 开始重新部署...'];
      
      // 模拟重新部署
      setTimeout(() => {
        this.simulateDeployment({
          cloudProvider: this.selectedCloud,
          region: this.selectedRegion,
          az: this.selectedAZ,
          vpc: this.vpcConfig,
          subnet: this.subnetConfig,
          components: this.selectedComponentsDetails
        });
      }, 1000);
    },
    backToSummary() {
      // 返回摘要页面
      this.isDeploying = false;
      this.isDeploymentComplete = false;
      this.deploymentStatus = 'idle';
      this.deploymentProgress = 0;
      this.deploymentMessage = '';
      this.deploymentLogs = [];
    }
  }
}
</script>

<style>
#app {
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.header {
  background-color: #3498db;
  color: white;
  padding: 1rem 2rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
.header h1 {
  margin: 0;
  font-size: 1.5rem;
}
.main-content {
  flex: 1;
  padding: 2rem;
}
.step-container {
  max-width: 1000px;
  margin: 0 auto;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}
.step-indicator {
  display: flex;
  background-color: #f5f5f5;
  padding: 1rem;
  overflow-x: auto;
}
.step {
  padding: 0.5rem 1rem;
  margin-right: 0.5rem;
  border-radius: 4px;
  font-size: 0.9rem;
  color: #666;
  white-space: nowrap;
}
.step.active {
  background-color: #3498db;
  color: white;
}
.step.completed {
  background-color: #2ecc71;
  color: white;
}
.step-content {
  padding: 2rem;
  min-height: 400px;
}
.step-panel {
  animation: fadeIn 0.3s ease-in-out;
}
.step-navigation {
  display: flex;
  justify-content: space-between;
  padding: 1rem 2rem 2rem;
  border-top: 1px solid #eee;
}
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
