<template>
  <div class="az-selector-container">
    <h2>选择可用区 (AZ)</h2>
    <div class="az-selection">
      <div v-if="loading" class="loading-indicator">
        <div class="spinner"></div>
        <p>正在加载可用区数据...</p>
      </div>
      <div v-else-if="availabilityZones.length === 0" class="no-data">
        没有找到可用区数据，请先选择区域
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
</template>

<script>
export default {
  name: 'AZSelector',
  props: {
    cloudProvider: {
      type: String,
      required: true
    },
    region: {
      type: String,
      required: true
    },
    initialAZs: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      selectedAZ: '',
      availabilityZones: this.initialAZs || [],
      loading: false
    }
  },
  watch: {
    region: {
      immediate: true,
      handler(newRegion) {
        if (newRegion && this.cloudProvider) {
          this.loadAvailabilityZones(this.cloudProvider, newRegion);
        } else {
          this.availabilityZones = [];
          this.selectedAZ = '';
        }
      }
    },
    initialAZs: {
      immediate: true,
      handler(newAZs) {
        if (newAZs && newAZs.length > 0) {
          this.availabilityZones = newAZs;
        }
      }
    }
  },
  methods: {
    selectAZ(az) {
      this.selectedAZ = az;
      this.$emit('az-selected', az);
    },
    loadAvailabilityZones(provider, region) {
      if (!provider || !region) return;
      
      this.loading = true;
      this.availabilityZones = [];
      this.selectedAZ = '';
      
      // 在实际应用中，这里会调用后端API从Terraform文件中提取可用区信息
      // 模拟API调用
      setTimeout(() => {
        this.availabilityZones = this.getMockAvailabilityZones(provider, region);
        this.loading = false;
      }, 800);
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
.az-selector-container {
  margin-bottom: 2rem;
}

.az-selection {
  margin-top: 1.5rem;
}

.az-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
}

.az-option {
  border: 2px solid #eee;
  border-radius: 8px;
  padding: 1rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.az-option:hover {
  background-color: #f9f9f9;
}

.az-option.selected {
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
