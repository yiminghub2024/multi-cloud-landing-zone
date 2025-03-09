<template>
  <div class="region-selector-container">
    <h2>选择区域 (Region)</h2>
    <div class="region-selection">
      <div v-if="loading" class="loading-indicator">
        <div class="spinner"></div>
        <p>正在加载区域数据...</p>
      </div>
      <div v-else-if="regions.length === 0" class="no-data">
        没有找到区域数据，请先选择云服务提供商
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
</template>

<script>
export default {
  name: 'RegionSelector',
  props: {
    cloudProvider: {
      type: String,
      required: true
    },
    initialRegions: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      selectedRegion: '',
      regions: this.initialRegions || [],
      loading: false
    }
  },
  watch: {
    cloudProvider: {
      immediate: true,
      handler(newProvider) {
        if (newProvider) {
          this.loadRegions(newProvider);
        } else {
          this.regions = [];
          this.selectedRegion = '';
        }
      }
    },
    initialRegions: {
      immediate: true,
      handler(newRegions) {
        if (newRegions && newRegions.length > 0) {
          this.regions = newRegions;
        }
      }
    }
  },
  methods: {
    selectRegion(region) {
      this.selectedRegion = region;
      this.$emit('region-selected', region);
      
      // 在实际应用中，这里会调用后端API获取所选区域的可用区数据
      this.fetchAvailabilityZones(this.cloudProvider, region);
    },
    loadRegions(provider) {
      if (!provider) return;
      
      this.loading = true;
      this.regions = [];
      this.selectedRegion = '';
      
      // 在实际应用中，这里会调用后端API从Terraform文件中提取区域信息
      // 模拟API调用
      setTimeout(() => {
        this.regions = this.getMockRegions(provider);
        this.loading = false;
      }, 800);
    },
    fetchAvailabilityZones(provider, region) {
      // 这里将来会实现与后端的API调用，从Terraform文件中提取可用区信息
      console.log(`Fetching availability zones for ${provider} in region ${region}...`);
      
      // 模拟API调用
      setTimeout(() => {
        // 假设这是从后端获取的数据
        const azs = this.getMockAvailabilityZones(provider, region);
        
        // 将获取的可用区数据发送给父组件
        this.$emit('azs-loaded', azs);
      }, 500);
    },
    getMockRegions(provider) {
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
      
      return mockRegions[provider] || [];
    },
    getMockAvailabilityZones(provider, region) {
      // 模拟不同区域的可用区数据
      const mockAZs = {
        'us-east-1': [
          { name: 'us-east-1a', value: 'us-east-1a' },
          { name: 'us-east-1b', value: 'us-east-1b' },
          { name: 'us-east-1c', value: 'us-east-1c' }
        ],
        'us-west-2': [
          { name: 'us-west-2a', value: 'us-west-2a' },
          { name: 'us-west-2b', value: 'us-west-2b' },
          { name: 'us-west-2c', value: 'us-west-2c' }
        ],
        'eastus': [
          { name: 'eastus-1', value: 'eastus-1' },
          { name: 'eastus-2', value: 'eastus-2' },
          { name: 'eastus-3', value: 'eastus-3' }
        ],
        'cn-beijing': [
          { name: '可用区A', value: 'cn-beijing-a' },
          { name: '可用区B', value: 'cn-beijing-b' },
          { name: '可用区C', value: 'cn-beijing-c' }
        ],
        'cn-shanghai': [
          { name: '可用区A', value: 'cn-shanghai-a' },
          { name: '可用区B', value: 'cn-shanghai-b' },
          { name: '可用区C', value: 'cn-shanghai-c' }
        ]
      };
      
      // 如果没有特定区域的数据，返回通用的可用区
      if (!mockAZs[region]) {
        return [
          { name: '可用区A', value: `${region}-a` },
          { name: '可用区B', value: `${region}-b` },
          { name: '可用区C', value: `${region}-c` }
        ];
      }
      
      return mockAZs[region];
    }
  }
}
</script>

<style scoped>
.region-selector-container {
  margin-bottom: 2rem;
}

.region-selection {
  margin-top: 1.5rem;
}

.region-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
}

.region-option {
  border: 2px solid #eee;
  border-radius: 8px;
  padding: 1rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.region-option:hover {
  background-color: #f9f9f9;
}

.region-option.selected {
  border-color: #3498db;
  background-color: rgba(52, 152, 219, 0.1);
}

.loading-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

.no-data {
  padding: 2rem;
  text-align: center;
  color: #999;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
