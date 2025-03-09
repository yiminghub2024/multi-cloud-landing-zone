<template>
  <div class="component-selector-container">
    <h2>选择云组件</h2>
    <div class="component-selection">
      <div v-if="loading" class="loading-indicator">
        <div class="spinner"></div>
        <p>正在加载云组件数据...</p>
      </div>
      <div v-else-if="cloudComponents.length === 0" class="no-data">
        没有找到云组件数据，请检查Terraform配置文件
      </div>
      <div v-else class="component-list">
        <div 
          v-for="(component, index) in cloudComponents" 
          :key="index" 
          class="component-option"
          :class="{ selected: selectedComponents.includes(component.value) }"
          @click="toggleComponent(component.value)"
        >
          <div class="component-header">
            <div class="component-name">{{ component.name }}</div>
            <div class="component-checkbox">
              <input 
                type="checkbox" 
                :checked="selectedComponents.includes(component.value)"
                @click.stop
              >
            </div>
          </div>
          <div class="component-desc">{{ component.description }}</div>
          <div v-if="component.properties && component.properties.length > 0" class="component-properties">
            <div 
              v-for="(prop, propIndex) in component.properties" 
              :key="propIndex"
              class="property-item"
            >
              <label :for="`prop-${component.value}-${propIndex}`">{{ prop.name }}</label>
              <input 
                :id="`prop-${component.value}-${propIndex}`"
                :type="prop.type === 'number' ? 'number' : 'text'"
                v-model="componentProperties[component.value][prop.key]"
                :placeholder="prop.placeholder || ''"
                :disabled="!selectedComponents.includes(component.value)"
                @input="updateComponentProperty(component.value, prop.key)"
              >
              <div class="property-desc">{{ prop.description }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ComponentSelector',
  props: {
    cloudProvider: {
      type: String,
      required: true
    },
    region: {
      type: String,
      required: true
    },
    az: {
      type: String,
      required: true
    },
    initialSelectedComponents: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      selectedComponents: this.initialSelectedComponents || [],
      cloudComponents: [],
      componentProperties: {},
      loading: false
    }
  },
  watch: {
    cloudProvider() {
      this.loadComponents();
    },
    region() {
      this.loadComponents();
    },
    initialSelectedComponents(newVal) {
      this.selectedComponents = newVal;
    }
  },
  mounted() {
    this.loadComponents();
  },
  methods: {
    loadComponents() {
      if (!this.cloudProvider || !this.region) return;
      
      this.loading = true;
      this.cloudComponents = [];
      this.selectedComponents = [];
      this.componentProperties = {};
      
      // 在实际应用中，这里会调用后端API从Terraform文件中提取组件信息
      // 模拟API调用
      setTimeout(() => {
        this.cloudComponents = this.getMockComponents(this.cloudProvider, this.region);
        
        // 初始化组件属性
        this.cloudComponents.forEach(component => {
          if (!this.componentProperties[component.value]) {
            this.componentProperties[component.value] = {};
          }
          
          if (component.properties) {
            component.properties.forEach(prop => {
              this.componentProperties[component.value][prop.key] = prop.defaultValue || '';
            });
          }
        });
        
        this.loading = false;
      }, 800);
    },
    toggleComponent(componentValue) {
      const index = this.selectedComponents.indexOf(componentValue);
      if (index === -1) {
        this.selectedComponents.push(componentValue);
      } else {
        this.selectedComponents.splice(index, 1);
      }
      
      this.$emit('components-updated', {
        selectedComponents: this.selectedComponents,
        componentProperties: this.componentProperties
      });
    },
    updateComponentProperty(componentValue, propertyKey) {
      this.$emit('components-updated', {
        selectedComponents: this.selectedComponents,
        componentProperties: this.componentProperties
      });
    },
    getMockComponents(provider, region) {
      // 模拟不同云服务提供商的组件数据
      const commonComponents = [
        {
          name: '负载均衡器',
          value: 'load-balancer',
          description: '用于分发网络流量的服务，提高应用程序的可用性和容错能力',
          properties: [
            {
              name: '实例数量',
              key: 'instance_count',
              type: 'number',
              defaultValue: '1',
              placeholder: '请输入实例数量',
              description: '负载均衡器实例的数量'
            },
            {
              name: '监听端口',
              key: 'listener_port',
              type: 'number',
              defaultValue: '80',
              placeholder: '请输入监听端口',
              description: '负载均衡器监听的端口'
            }
          ]
        },
        {
          name: '对象存储',
          value: 'object-storage',
          description: '用于存储和检索任意数量数据的服务，适用于静态网站、备份和归档等场景',
          properties: [
            {
              name: '存储桶名称',
              key: 'bucket_name',
              type: 'text',
              defaultValue: '',
              placeholder: '请输入全局唯一的存储桶名称',
              description: '存储桶名称必须全局唯一'
            },
            {
              name: '存储类型',
              key: 'storage_class',
              type: 'text',
              defaultValue: 'Standard',
              placeholder: '例如: Standard, IA, Archive',
              description: '存储类型决定数据的访问频率和成本'
            }
          ]
        },
        {
          name: '关系型数据库',
          value: 'rds',
          description: '托管的关系型数据库服务，支持多种数据库引擎，自动备份和高可用性',
          properties: [
            {
              name: '数据库引擎',
              key: 'engine',
              type: 'text',
              defaultValue: 'MySQL',
              placeholder: '例如: MySQL, PostgreSQL',
              description: '数据库引擎类型'
            },
            {
              name: '实例类型',
              key: 'instance_type',
              type: 'text',
              defaultValue: 'small',
              placeholder: '例如: small, medium, large',
              description: '数据库实例的规格'
            },
            {
              name: '存储容量(GB)',
              key: 'storage_size',
              type: 'number',
              defaultValue: '20',
              placeholder: '请输入存储容量',
              description: '数据库存储容量，单位为GB'
            }
          ]
        },
        {
          name: '弹性计算实例',
          value: 'compute',
          description: '可扩展的计算容量，适用于各种应用场景，支持多种操作系统和配置',
          properties: [
            {
              name: '实例数量',
              key: 'instance_count',
              type: 'number',
              defaultValue: '2',
              placeholder: '请输入实例数量',
              description: '计算实例的数量'
            },
            {
              name: '实例类型',
              key: 'instance_type',
              type: 'text',
              defaultValue: 'medium',
              placeholder: '例如: small, medium, large',
              description: '计算实例的规格'
            },
            {
              name: '操作系统',
              key: 'os',
              type: 'text',
              defaultValue: 'Linux',
              placeholder: '例如: Linux, Windows',
              description: '实例的操作系统类型'
            }
          ]
        },
        {
          name: 'CDN',
          value: 'cdn',
          description: '内容分发网络服务，加速静态内容分发，提高用户访问速度和体验',
          properties: [
            {
              name: '域名',
              key: 'domain',
              type: 'text',
              defaultValue: '',
              placeholder: '请输入加速域名',
              description: '需要加速的域名'
            },
            {
              name: '源站类型',
              key: 'origin_type',
              type: 'text',
              defaultValue: 'OSS',
              placeholder: '例如: OSS, ECS, Custom',
              description: '内容源站的类型'
            }
          ]
        }
      ];
      
      // 根据不同的云服务提供商添加特定组件
      let providerSpecificComponents = [];
      
      switch(provider) {
        case 'aws':
          providerSpecificComponents = [
            {
              name: 'Lambda函数',
              value: 'lambda',
              description: 'AWS Lambda是一项无服务器计算服务，无需预置或管理服务器即可运行代码',
              properties: [
                {
                  name: '运行时',
                  key: 'runtime',
                  type: 'text',
                  defaultValue: 'nodejs14.x',
                  placeholder: '例如: nodejs14.x, python3.9',
                  description: 'Lambda函数的运行时环境'
                },
                {
                  name: '内存大小(MB)',
                  key: 'memory_size',
                  type: 'number',
                  defaultValue: '128',
                  placeholder: '请输入内存大小',
                  description: 'Lambda函数的内存大小，单位为MB'
                }
              ]
            }
          ];
          break;
        case 'azure':
          providerSpecificComponents = [
            {
              name: 'Azure Functions',
              value: 'azure-functions',
              description: 'Azure Functions是一项无服务器计算服务，可以运行事件驱动的代码而无需管理基础设施',
              properties: [
                {
                  name: '运行时',
                  key: 'runtime',
                  type: 'text',
                  defaultValue: 'node',
                  placeholder: '例如: node, dotnet, java',
                  description: 'Azure Functions的运行时环境'
                }
              ]
            }
          ];
          break;
        case 'alicloud':
          providerSpecificComponents = [
            {
              name: '函数计算',
              value: 'fc',
              description: '阿里云函数计算是一个事件驱动的全托管计算服务，无需管理服务器等基础设施',
              properties: [
                {
                  name: '运行时',
                  key: 'runtime',
                  type: 'text',
                  defaultValue: 'nodejs10',
                  placeholder: '例如: nodejs10, python3',
                  description: '函数计算的运行时环境'
                }
              ]
            }
          ];
          break;
        // 其他云服务提供商的特定组件...
      }
      
      return [...commonComponents, ...providerSpecificComponents];
    }
  }
}
</script>

<style scoped>
.component-selector-container {
  margin-bottom: 2rem;
}

.component-selection {
  margin-top: 1.5rem;
}

.component-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.component-option {
  border: 2px solid #eee;
  border-radius: 8px;
  padding: 1.5rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.component-option:hover {
  background-color: #f9f9f9;
}

.component-option.selected {
  border-color: #3498db;
  background-color: rgba(52, 152, 219, 0.05);
}

.component-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.component-name {
  font-weight: bold;
  font-size: 1.1rem;
}

.component-checkbox {
  display: flex;
  align-items: center;
}

.component-checkbox input {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.component-desc {
  color: #666;
  margin-bottom: 1rem;
}

.component-properties {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px dashed #ddd;
}

.property-item {
  margin-bottom: 1rem;
}

.property-item label {
  display: block;
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.property-item input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.property-item input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.property-item input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.property-desc {
  font-size: 0.9rem;
  color: #7f8c8d;
  margin-top: 0.5rem;
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
