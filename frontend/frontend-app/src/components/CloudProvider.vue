<template>
  <div class="cloud-provider-container">
    <h2>选择云服务提供商</h2>
    <div class="cloud-selection">
      <div 
        v-for="(cloud, index) in cloudProviders" 
        :key="index" 
        class="cloud-option"
        :class="{ selected: selectedCloud === cloud.value }"
        @click="selectCloud(cloud.value)"
      >
        <div class="cloud-logo">
          <img v-if="cloud.logoUrl" :src="cloud.logoUrl" :alt="cloud.name + ' logo'" class="provider-logo">
          <div v-else class="logo-placeholder">{{ cloud.logo }}</div>
        </div>
        <div class="cloud-name">{{ cloud.name }}</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CloudProvider',
  data() {
    return {
      selectedCloud: '',
      cloudProviders: [
        { 
          name: 'AWS', 
          value: 'aws', 
          logo: 'AWS',
          logoUrl: null,
          description: 'Amazon Web Services'
        },
        { 
          name: 'Azure', 
          value: 'azure', 
          logo: 'Azure',
          logoUrl: null,
          description: 'Microsoft Azure'
        },
        { 
          name: '阿里云', 
          value: 'alicloud', 
          logo: '阿里云',
          logoUrl: null,
          description: 'Alibaba Cloud'
        },
        { 
          name: '百度云', 
          value: 'baidu', 
          logo: '百度云',
          logoUrl: null,
          description: 'Baidu Cloud'
        },
        { 
          name: '华为云', 
          value: 'huawei', 
          logo: '华为云',
          logoUrl: null,
          description: 'Huawei Cloud'
        },
        { 
          name: '腾讯云', 
          value: 'tencent', 
          logo: '腾讯云',
          logoUrl: null,
          description: 'Tencent Cloud'
        },
        { 
          name: '火山云', 
          value: 'volcengine', 
          logo: '火山云',
          logoUrl: null,
          description: 'Volcengine Cloud'
        }
      ]
    }
  },
  methods: {
    selectCloud(cloud) {
      this.selectedCloud = cloud;
      this.$emit('cloud-selected', cloud);
      
      // 在实际应用中，这里会调用后端API获取所选云的区域数据
      this.fetchCloudRegions(cloud);
    },
    fetchCloudRegions(cloud) {
      // 这里将来会实现与后端的API调用，从Terraform文件中提取区域信息
      console.log(`Fetching regions for ${cloud}...`);
      
      // 模拟API调用
      setTimeout(() => {
        // 假设这是从后端获取的数据
        const regions = this.getMockRegions(cloud);
        
        // 将获取的区域数据发送给父组件
        this.$emit('regions-loaded', regions);
      }, 500);
    },
    getMockRegions(cloud) {
      // 模拟不同云服务提供商的区域数据
      const mockRegions = {
        aws: [
          { name: '美国东部 (弗吉尼亚北部)', value: 'us-east-1' },
          { name: '美国东部 (俄亥俄)', value: 'us-east-2' },
          { name: '美国西部 (加利福尼亚北部)', value: 'us-west-1' },
          { name: '美国西部 (俄勒冈)', value: 'us-west-2' },
          { name: '亚太地区 (香港)', value: 'ap-east-1' },
          { name: '亚太地区 (东京)', value: 'ap-northeast-1' }
        ],
        azure: [
          { name: '美国东部', value: 'eastus' },
          { name: '美国东部2', value: 'eastus2' },
          { name: '美国西部', value: 'westus' },
          { name: '美国西部2', value: 'westus2' },
          { name: '东亚', value: 'eastasia' },
          { name: '东南亚', value: 'southeastasia' }
        ],
        alicloud: [
          { name: '华北 1 (青岛)', value: 'cn-qingdao' },
          { name: '华北 2 (北京)', value: 'cn-beijing' },
          { name: '华北 3 (张家口)', value: 'cn-zhangjiakou' },
          { name: '华东 1 (杭州)', value: 'cn-hangzhou' },
          { name: '华东 2 (上海)', value: 'cn-shanghai' },
          { name: '华南 1 (深圳)', value: 'cn-shenzhen' }
        ],
        baidu: [
          { name: '华北-北京', value: 'bj' },
          { name: '华南-广州', value: 'gz' },
          { name: '华东-苏州', value: 'su' }
        ],
        huawei: [
          { name: '华北-北京一', value: 'cn-north-1' },
          { name: '华北-北京四', value: 'cn-north-4' },
          { name: '华东-上海一', value: 'cn-east-3' },
          { name: '华南-广州', value: 'cn-south-1' },
          { name: '亚太-香港', value: 'ap-southeast-1' }
        ],
        tencent: [
          { name: '华南地区(广州)', value: 'ap-guangzhou' },
          { name: '华东地区(上海)', value: 'ap-shanghai' },
          { name: '华北地区(北京)', value: 'ap-beijing' },
          { name: '西南地区(成都)', value: 'ap-chengdu' },
          { name: '西南地区(重庆)', value: 'ap-chongqing' },
          { name: '港澳台地区(中国香港)', value: 'ap-hongkong' }
        ],
        volcengine: [
          { name: '华北-北京', value: 'cn-beijing' },
          { name: '华东-上海', value: 'cn-shanghai' },
          { name: '华南-广州', value: 'cn-guangzhou' }
        ]
      };
      
      return mockRegions[cloud] || [];
    }
  }
}
</script>

<style scoped>
.cloud-provider-container {
  margin-bottom: 2rem;
}

.cloud-selection {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1rem;
  margin-top: 1.5rem;
}

.cloud-option {
  border: 2px solid #eee;
  border-radius: 8px;
  padding: 1.5rem 1rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cloud-option:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.cloud-option.selected {
  border-color: #3498db;
  background-color: rgba(52, 152, 219, 0.1);
}

.cloud-logo {
  font-size: 2rem;
  margin-bottom: 1rem;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.provider-logo {
  max-width: 100%;
  max-height: 60px;
}

.logo-placeholder {
  width: 60px;
  height: 60px;
  background-color: #f5f5f5;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}

.cloud-name {
  font-weight: bold;
}
</style>
