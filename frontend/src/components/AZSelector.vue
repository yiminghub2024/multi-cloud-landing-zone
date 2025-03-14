<template>
  <div class="az-selector-container">
    <el-card shadow="hover">
      <template #header>
        <h2>选择可用区 (AZ)</h2>
      </template>
      <div class="az-selection">
        <el-skeleton v-if="loading" :rows="3" animated />
        
        <el-empty v-else-if="availabilityZones.length === 0" description="没有找到可用区数据，请先选择区域" />
        
        <el-row v-else :gutter="20">
          <el-col 
            v-for="(az, index) in availabilityZones" 
            :key="index" 
            :xs="24" 
            :sm="12" 
            :md="8" 
            :lg="6"
            class="az-col"
          >
            <el-card 
              class="az-option" 
              :class="{ 'is-selected': selectedAZ === az.value }"
              shadow="hover"
              @click="selectAZ(az.value)"
            >
              {{ az.name }}
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>
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
  margin-top: 1rem;
}

.az-col {
  margin-bottom: 20px;
}

.az-option {
  cursor: pointer;
  transition: all 0.3s;
  height: 100%;
  text-align: center;
}

.az-option:hover {
  transform: translateY(-5px);
}

.az-option.is-selected {
  border-color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
}
</style>
