<template>
  <div class="region-selector-container">
    <el-card shadow="hover">
      <template #header>
        <h2>选择区域 (Region)</h2>
      </template>
      <div class="region-selection">
        <el-skeleton v-if="loading" :rows="3" animated />
        
        <el-empty v-else-if="regions.length === 0" description="没有找到区域数据，请先选择云服务提供商" />
        
        <el-row v-else :gutter="20">
          <el-col 
            v-for="(region, index) in regions" 
            :key="index" 
            :xs="24" 
            :sm="12" 
            :md="8" 
            :lg="6"
            class="region-col"
          >
            <el-card 
              class="region-option" 
              :class="{ 'is-selected': selectedRegion === region.value }"
              shadow="hover"
              @click="selectRegion(region.value)"
            >
              {{ region.name }}
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>
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
          // AWS 全球区域
          { name: '美国东部 (弗吉尼亚北部)', value: 'us-east-1' },
          { name: '美国东部 (俄亥俄)', value: 'us-east-2' },
          { name: '美国西部 (加利福尼亚北部)', value: 'us-west-1' },
          { name: '美国西部 (俄勒冈)', value: 'us-west-2' },
          { name: '加拿大 (蒙特利尔)', value: 'ca-central-1' },
          { name: '加拿大 (卡尔加里)', value: 'ca-west-1' },
          { name: '欧洲 (斯德哥尔摩)', value: 'eu-north-1' },
          { name: '欧洲 (巴黎)', value: 'eu-west-3' },
          { name: '欧洲 (伦敦)', value: 'eu-west-2' },
          { name: '欧洲 (爱尔兰)', value: 'eu-west-1' },
          { name: '欧洲 (法兰克福)', value: 'eu-central-1' },
          { name: '欧洲 (米兰)', value: 'eu-south-1' },
          { name: '欧洲 (西班牙)', value: 'eu-south-2' },
          { name: '欧洲 (苏黎世)', value: 'eu-central-2' },
          { name: '亚太地区 (孟买)', value: 'ap-south-1' },
          { name: '亚太地区 (海得拉巴)', value: 'ap-south-2' },
          { name: '亚太地区 (东京)', value: 'ap-northeast-1' },
          { name: '亚太地区 (首尔)', value: 'ap-northeast-2' },
          { name: '亚太地区 (大阪)', value: 'ap-northeast-3' },
          { name: '亚太地区 (新加坡)', value: 'ap-southeast-1' },
          { name: '亚太地区 (悉尼)', value: 'ap-southeast-2' },
          { name: '亚太地区 (雅加达)', value: 'ap-southeast-3' },
          { name: '亚太地区 (墨尔本)', value: 'ap-southeast-4' },
          { name: '亚太地区 (马来西亚)', value: 'ap-southeast-5' },
          { name: '亚太地区 (泰国)', value: 'ap-southeast-7' },
          { name: '亚太地区 (香港)', value: 'ap-east-1' },
          { name: '南美洲 (圣保罗)', value: 'sa-east-1' },
          { name: '中东 (巴林)', value: 'me-south-1' },
          { name: '中东 (阿联酋)', value: 'me-central-1' },
          { name: '非洲 (开普敦)', value: 'af-south-1' },
          { name: '以色列 (特拉维夫)', value: 'il-central-1' },
          { name: '墨西哥 (中部)', value: 'mx-central-1' },
          { name: 'AWS GovCloud (美国东部)', value: 'us-gov-east-1' },
          { name: 'AWS GovCloud (美国西部)', value: 'us-gov-west-1' },
          
          // AWS 中国区域
          { name: '中国 (北京)', value: 'cn-north-1' },
          { name: '中国 (宁夏)', value: 'cn-northwest-1' }
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
  margin-top: 1rem;
}

.region-col {
  margin-bottom: 20px;
}

.region-option {
  cursor: pointer;
  transition: all 0.3s;
  height: 100%;
  text-align: center;
}

.region-option:hover {
  transform: translateY(-5px);
}

.region-option.is-selected {
  border-color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
}
</style>
