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
          <div class="step" :class="{ active: currentStep >= 7, completed: currentStep > 7 }">7. 部署摘要</div>
        </div>
        
        <div class="step-content">
          <!-- 步骤1: 选择云服务提供商 -->
          <div v-if="currentStep === 1" class="step-panel">
            <h2>选择云服务提供商</h2>
            <div class="cloud-selection">
              <div 
                v-for="(cloud, index) in cloudProviders" 
                :key="index" 
                class="cloud-option"
                :class="{ selected: selectedCloud === cloud.value }"
                @click="selectCloud(cloud.value)"
              >
                <div class="cloud-logo">{{ cloud.logo }}</div>
                <div class="cloud-name">{{ cloud.name }}</div>
              </div>
            </div>
          </div>
          
          <!-- 步骤2: 选择区域 -->
          <div v-if="currentStep === 2" class="step-panel">
            <h2>选择区域 (Region)</h2>
            <div class="region-selection">
              <div v-if="regions.length === 0" class="no-data">
                正在加载区域数据...
              </div>
              <div v-else class="region-list">
                <div 
                  v-for="(region, index) in regions" 
                  :key="index" 
                  class="region-option"
                  :class="{ selected: selectedRegion === region.value }"
                  @click="selectRegion(region.value)"
                >
                  {{ region.name }}
                </div>
              </div>
            </div>
          </div>
          
          <!-- 步骤3: 选择可用区 -->
          <div v-if="currentStep === 3" class="step-panel">
            <h2>选择可用区 (AZ)</h2>
            <div class="az-selection">
              <div v-if="availabilityZones.length === 0" class="no-data">
                正在加载可用区数据...
              </div>
              <div v-else class="az-list">
                <div 
                  v-for="(az, index) in availabilityZones" 
                  :key="index" 
                  class="az-option"
                  :class="{ selected: selectedAZ === az.value }"
                  @click="selectAZ(az.value)"
                >
                  {{ az.name }}
                </div>
              </div>
            </div>
          </div>
          
          <!-- 步骤4: 配置VPC -->
          <div v-if="currentStep === 4" class="step-panel">
            <h2>配置VPC</h2>
            <div class="vpc-config">
              <div class="form-group">
                <label for="vpc-name">VPC名称</label>
                <input 
                  id="vpc-name" 
                  type="text" 
                  v-model="vpcName" 
                  placeholder="请输入VPC名称"
                />
              </div>
              <div class="form-group">
                <label for="vpc-cidr">CIDR块</label>
                <input 
                  id="vpc-cidr" 
                  type="text" 
                  v-model="vpcCIDR" 
                  placeholder="例如: 10.0.0.0/16"
                />
              </div>
            </div>
          </div>
          
          <!-- 步骤5: 配置子网 -->
          <div v-if="currentStep === 5" class="step-panel">
            <h2>配置子网</h2>
            <div class="subnet-config">
              <div class="form-group">
                <label for="subnet-name">子网名称</label>
                <input 
                  id="subnet-name" 
                  type="text" 
                  v-model="subnetName" 
                  placeholder="请输入子网名称"
                />
              </div>
              <div class="form-group">
                <label for="subnet-cidr">CIDR块</label>
                <input 
                  id="subnet-cidr" 
                  type="text" 
                  v-model="subnetCIDR" 
                  placeholder="例如: 10.0.1.0/24"
                />
              </div>
              <div class="selected-vpc">
                <p>所属VPC: {{ vpcName }}</p>
              </div>
            </div>
          </div>
          
          <!-- 步骤6: 选择云组件 -->
          <div v-if="currentStep === 6" class="step-panel">
            <h2>选择云组件</h2>
            <div class="component-selection">
              <div v-if="cloudComponents.length === 0" class="no-data">
                正在加载云组件数据...
              </div>
              <div v-else class="component-list">
                <div 
                  v-for="(component, index) in cloudComponents" 
                  :key="index" 
                  class="component-option"
                  :class="{ selected: selectedComponents.includes(component.value) }"
                  @click="toggleComponent(component.value)"
                >
                  <div class="component-name">{{ component.name }}</div>
                  <div class="component-desc">{{ component.description }}</div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 步骤7: 部署摘要 -->
          <div v-if="currentStep === 7" class="step-panel">
            <h2>部署摘要</h2>
            <div class="deployment-summary">
              <div class="summary-item">
                <div class="summary-label">云服务提供商:</div>
                <div class="summary-value">{{ getCloudProviderName(selectedCloud) }}</div>
              </div>
              <div class="summary-item">
                <div class="summary-label">区域:</div>
                <div class="summary-value">{{ getRegionName(selectedRegion) }}</div>
              </div>
              <div class="summary-item">
                <div class="summary-label">可用区:</div>
                <div class="summary-value">{{ getAZName(selectedAZ) }}</div>
              </div>
              <div class="summary-item">
                <div class="summary-label">VPC:</div>
                <div class="summary-value">{{ vpcName }} ({{ vpcCIDR }})</div>
              </div>
              <div class="summary-item">
                <div class="summary-label">子网:</div>
                <div class="summary-value">{{ subnetName }} ({{ subnetCIDR }})</div>
              </div>
              <div class="summary-item">
                <div class="summary-label">选择的组件:</div>
                <div class="summary-value">
                  <ul>
                    <li v-for="(component, index) in selectedComponentsDetails" :key="index">
                      {{ component.name }}
                    </li>
                  </ul>
                </div>
              </div>
              <div class="deployment-actions">
                <button class="deploy-button" @click="startDeployment">开始部署</button>
              </div>
            </div>
          </div>
          
          <!-- 部署状态 -->
          <div v-if="isDeploying" class="deployment-status">
            <h2>正在部署</h2>
            <div class="progress-bar">
              <div class="progress" :style="{ width: deploymentProgress + '%' }"></div>
            </div>
            <div class="status-message">{{ deploymentStatus }}</div>
          </div>
          
          <!-- 部署完成 -->
          <div v-if="isDeploymentComplete" class="deployment-complete">
            <h2>部署完成</h2>
            <div class="topology-container">
              <h3>资源拓扑图</h3>
              <div class="topology-placeholder">
                <!-- 这里将显示拓扑图 -->
                拓扑图将在这里显示
              </div>
            </div>
            <div class="deployment-details">
              <h3>部署详情</h3>
              <pre>{{ deploymentDetails }}</pre>
            </div>
          </div>
        </div>
        
        <div class="step-navigation">
          <button 
            v-if="currentStep > 1" 
            class="nav-button prev" 
            @click="prevStep"
          >
            上一步
          </button>
          <button 
            v-if="currentStep < 7 && isCurrentStepValid" 
            class="nav-button next" 
            @click="nextStep"
          >
            下一步
          </button>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      currentStep: 1,
      cloudProviders: [
        { name: 'AWS', value: 'aws', logo: 'AWS' },
        { name: 'Azure', value: 'azure', logo: 'Azure' },
        { name: '阿里云', value: 'alicloud', logo: '阿里云' },
        { name: '百度云', value: 'baidu', logo: '百度云' },
        { name: '华为云', value: 'huawei', logo: '华为云' },
        { name: '腾讯云', value: 'tencent', logo: '腾讯云' },
        { name: '火山云', value: 'volcengine', logo: '火山云' }
      ],
      selectedCloud: '',
      regions: [],
      selectedRegion: '',
      availabilityZones: [],
      selectedAZ: '',
      vpcName: '',
      vpcCIDR: '',
      subnetName: '',
      subnetCIDR: '',
      cloudComponents: [],
      selectedComponents: [],
      isDeploying: false,
      deploymentProgress: 0,
      deploymentStatus: '',
      isDeploymentComplete: false,
      deploymentDetails: '',
    }
  },
  computed: {
    isCurrentStepValid() {
      switch (this.currentStep) {
        case 1:
          return this.selectedCloud !== '';
        case 2:
          return this.selectedRegion !== '';
        case 3:
          return this.selectedAZ !== '';
        case 4:
          return this.vpcName !== '' && this.vpcCIDR !== '';
        case 5:
          return this.subnetName !== '' && this.subnetCIDR !== '';
        case 6:
          return this.selectedComponents.length > 0;
        default:
          return true;
      }
    },
    selectedComponentsDetails() {
      return this.cloudComponents.filter(component => 
        this.selectedComponents.includes(component.value)
      );
    }
  },
  methods: {
    selectCloud(cloud) {
      this.selectedCloud = cloud;
      // 在实际应用中，这里会从后端获取所选云的区域数据
      this.loadRegions(cloud);
    },
    loadRegions(cloud) {
      // 模拟从后端加载区域数据
      // 在实际应用中，这里会调用API从Terraform文件中提取区域信息
      setTimeout(() => {
        this.regions = [
          { name: '华北区域', value: 'north-china' },
          { name: '华东区域', value: 'east-china' },
          { name: '华南区域', value: 'south-china' },
          { name: '美国西部', value: 'us-west' },
          { name: '美国东部', value: 'us-east' },
          { name: '欧洲', value: 'europe' }
        ];
      }, 500);
    },
    selectRegion(region) {
      this.selectedRegion = region;
      // 在实际应用中，这里会从后端获取所选区域的可用区数据
      this.loadAvailabilityZones(this.selectedCloud, region);
    },
    loadAvailabilityZones(cloud, region) {
      // 模拟从后端加载可用区数据
      // 在实际应用中，这里会调用API从Terraform文件中提取可用区信息
      setTimeout(() => {
        this.availabilityZones = [
          { name: '可用区A', value: 'zone-a' },
          { name: '可用区B', value: 'zone-b' },
          { name: '可用区C', value: 'zone-c' }
        ];
      }, 500);
    },
    selectAZ(az) {
      this.selectedAZ = az;
    },
    loadCloudComponents(cloud, region, az) {
      // 模拟从后端加载云组件数据
      // 在实际应用中，这里会调用API从Terraform文件中提取组件信息
      setTimeout(() => {
        this.cloudComponents = [
          { 
            name: '负载均衡器', 
            value: 'load-balancer',
            description: '用于分发网络流量的服务'
          },
          { 
            name: '对象存储', 
            value: 'object-storage',
            description: '用于存储和检索任意数量数据的服务'
          },
          { 
            name: '关系型数据库', 
            value: 'rds',
            description: '托管的关系型数据库服务'
          },
          { 
            name: '弹性计算实例', 
            value: 'compute',
            description: '可扩展的计算容量'
          },
          { 
            name: 'CDN', 
            value: 'cdn',
            description: '内容分发网络服务'
          }
        ];
      }, 500);
    },
    toggleComponent(component) {
      const index = this.selectedComponents.indexOf(component);
      if (index === -1) {
        this.selectedComponents.push(component);
      } else {
        this.selectedComponents.splice(index, 1);
      }
    },
    getCloudProviderName(value) {
      const provider = this.cloudProviders.find(p => p.value === value);
      return provider ? provider.name : '';
    },
    getRegionName(value) {
      const region = this.regions.find(r => r.value === value);
      return region ? region.name : '';
    },
    getAZName(value) {
      const az = this.availabilityZones.find(a => a.value === value);
      return az ? az.name : '';
    },
    nextStep() {
      if (this.currentStep < 7 && this.isCurrentStepValid) {
        this.currentStep++;
        
        // 根据步骤加载相应数据
        if (this.currentStep === 6) {
          this.loadCloudComponents(this.selectedCloud, this.selectedRegion, this.selectedAZ);
        }
      }
    },
    prevStep() {
      if (this.currentStep > 1) {
        this.currentStep--;
      }
    },
    startDeployment() {
      this.isDeploying = true;
      this.deploymentProgress = 0;
      this.deploymentStatus = '正在准备部署资源...';
      
      // 模拟部署过程
      const interval = setInterval(() => {
        this.deploymentProgress += 10;
        
        if (this.deploymentProgress < 30) {
          this.deploymentStatus = '正在创建VPC和子网...';
        } else if (this.deploymentProgress < 60) {
          this.deploymentStatus = '正在部署选定的云组件...';
        } else if (this.deploymentProgress < 90) {
          this.deploymentStatus = '正在配置网络连接...';
        } else {
          this.deploymentStatus = '正在完成部署...';
        }
        
        if (this.deploymentProgress >= 100) {
          clearInterval(interval);
          this.deploymentComplete();
        }
      }, 1000);
    },
    deploymentComplete() {
      this.isDeploying = false;
      this.isDeploymentComplete = true;
      this.deploymentDetails = JSON.stringify({
        provider: this.getCloudProviderName(this.selectedCloud),
        region: this.getRegionName(this.selectedRegion),
        az: this.getAZName(this.selectedAZ),
        vpc: {
          name: this.vpcName,
          cidr: this.vpcCIDR
        },
        subnet: {
          name: this.subnetName,
          cidr: this.subnetCIDR
        },
        components: this.selectedComponentsDetails.map(c => c.name)
      }, null, 2);
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin: 0;
  padding: 0;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: #2c3e50;
  color: white;
  padding: 1rem;
  text-align: center;
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

.step<response clipped><NOTE>To save on context only part of this file has been shown to you. You should retry this tool after you have searched inside the file with `grep -n` in order to find the line numbers of what you are looking for.</NOTE>