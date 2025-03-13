<template>
  <div class="component-selector-container">
    <el-card shadow="hover">
      <template #header>
        <h2>选择云组件</h2>
      </template>
      <div class="component-selection">
        <el-skeleton v-if="loading" :rows="5" animated />
        
        <el-empty v-else-if="cloudComponents.length === 0" description="没有找到云组件数据，请检查Terraform配置文件" />
        
        <div v-else class="component-list">
          <el-collapse accordion>
            <el-collapse-item 
              v-for="(component, index) in cloudComponents" 
              :key="index"
              :title="component.name"
              :name="component.value"
            >
              <template #extra>
                <el-checkbox 
                  v-model="selectedComponents" 
                  :label="component.value"
                  @change="() => toggleComponent(component.value)"
                  @click.stop
                />
              </template>
              
              <el-card shadow="hover" class="component-card">
                <el-descriptions :column="1" border>
                  <el-descriptions-item label="描述">{{ component.description }}</el-descriptions-item>
                </el-descriptions>
                
                <div v-if="component.properties && component.properties.length > 0" class="component-properties">
                  <el-divider content-position="left">组件属性</el-divider>
                  
                  <el-form label-position="top">
                    <el-form-item 
                      v-for="(prop, propIndex) in component.properties" 
                      :key="propIndex"
                      :label="prop.name"
                    >
                      <el-input 
                        v-if="prop.type !== 'number'"
                        v-model="componentProperties[component.value][prop.key]"
                        :placeholder="prop.placeholder || ''"
                        :disabled="!selectedComponents.includes(component.value)"
                        @input="updateComponentProperty(component.value, prop.key)"
                      />
                      <el-input-number 
                        v-else
                        v-model="componentProperties[component.value][prop.key]"
                        :placeholder="prop.placeholder || ''"
                        :disabled="!selectedComponents.includes(component.value)"
                        @change="updateComponentProperty(component.value, prop.key)"
                        :min="0"
                        controls-position="right"
                      />
                      <el-text type="info" size="small">{{ prop.description }}</el-text>
                    </el-form-item>
                  </el-form>
                </div>
              </el-card>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
    </el-card>
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
        this.cloudComponents = this.getCloudComponents(this.cloudProvider, this.region);
        
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
    // eslint-disable-next-line no-unused-vars
    updateComponentProperty(componentValue, propertyKey) {
      // 这里可以根据需要使用componentValue和propertyKey参数
      // 例如：console.log(`更新组件 ${componentValue} 的属性 ${propertyKey}`);
      
      this.$emit('components-updated', {
        selectedComponents: this.selectedComponents,
        componentProperties: this.componentProperties
      });
    },
    // eslint-disable-next-line no-unused-vars
    getCloudComponents(provider, region) {
      // 这里可以根据需要使用region参数
      // 例如：console.log(`获取 ${provider} 在 ${region} 区域的组件`);
      
      // 所有云服务提供商通用的组件
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
              type: 'string',
              defaultValue: '',
              placeholder: '请输入存储桶名称',
              description: '全局唯一的存储桶名称'
            },
            {
              name: '访问权限',
              key: 'access_control',
              type: 'string',
              defaultValue: 'private',
              placeholder: '请输入访问权限',
              description: '存储桶的访问权限，如private、public-read等'
            }
          ]
        },
        {
          name: '数据库',
          value: 'database',
          description: '托管的关系型数据库服务，提供高可用性和可扩展性',
          properties: [
            {
              name: '数据库引擎',
              key: 'engine',
              type: 'string',
              defaultValue: 'mysql',
              placeholder: '请输入数据库引擎',
              description: '数据库引擎类型，如mysql、postgresql等'
            },
            {
              name: '数据库版本',
              key: 'engine_version',
              type: 'string',
              defaultValue: '5.7',
              placeholder: '请输入数据库版本',
              description: '数据库引擎版本'
            },
            {
              name: '实例类型',
              key: 'instance_class',
              type: 'string',
              defaultValue: 'db.t3.micro',
              placeholder: '请输入实例类型',
              description: '数据库实例的计算和内存容量'
            },
            {
              name: '存储容量(GB)',
              key: 'allocated_storage',
              type: 'number',
              defaultValue: '20',
              placeholder: '请输入存储容量',
              description: '数据库的存储容量，单位为GB'
            },
            {
              name: '数据库名称',
              key: 'db_name',
              type: 'string',
              defaultValue: 'mydb',
              placeholder: '请输入数据库名称',
              description: '初始数据库的名称'
            },
            {
              name: '用户名',
              key: 'username',
              type: 'string',
              defaultValue: 'admin',
              placeholder: '请输入用户名',
              description: '数据库管理员用户名'
            },
            {
              name: '密码',
              key: 'password',
              type: 'string',
              defaultValue: '',
              placeholder: '请输入密码',
              description: '数据库管理员密码'
            }
          ]
        }
      ];
      
      // 根据云服务提供商返回特定的组件
      switch (provider) {
        case 'aws':
          return [
            ...commonComponents,
            {
              name: 'EC2实例',
              value: 'ec2',
              description: 'AWS的虚拟服务器，提供可调整的计算能力',
              properties: [
                {
                  name: '实例类型',
                  key: 'instance_type',
                  type: 'string',
                  defaultValue: 't2.micro',
                  placeholder: '请输入实例类型',
                  description: 'EC2实例的类型，决定CPU、内存、存储和网络容量'
                },
                {
                  name: 'AMI ID',
                  key: 'ami_id',
                  type: 'string',
                  defaultValue: 'ami-0c55b159cbfafe1f0',
                  placeholder: '请输入AMI ID',
                  description: '用于启动实例的Amazon机器映像ID'
                }
              ]
            }
          ];
        case 'azure':
          return [
            ...commonComponents,
            {
              name: '虚拟机',
              value: 'vm',
              description: 'Azure的虚拟机，提供可扩展的计算资源',
              properties: [
                {
                  name: '虚拟机大小',
                  key: 'vm_size',
                  type: 'string',
                  defaultValue: 'Standard_B1s',
                  placeholder: '请输入虚拟机大小',
                  description: '虚拟机的大小，决定CPU、内存和存储容量'
                },
                {
                  name: '操作系统',
                  key: 'os_type',
                  type: 'string',
                  defaultValue: 'Linux',
                  placeholder: '请输入操作系统类型',
                  description: '虚拟机的操作系统类型'
                }
              ]
            }
          ];
        case 'alicloud':
          return [
            ...commonComponents,
            {
              name: 'ECS实例',
              value: 'ecs',
              description: '阿里云的弹性计算服务，提供可伸缩的计算能力',
              properties: [
                {
                  name: '实例规格',
                  key: 'instance_type',
                  type: 'string',
                  defaultValue: 'ecs.t5-lc1m1.small',
                  placeholder: '请输入实例规格',
                  description: 'ECS实例的规格，决定CPU、内存、存储和网络容量'
                },
                {
                  name: '系统盘大小(GB)',
                  key: 'system_disk_size',
                  type: 'number',
                  defaultValue: '40',
                  placeholder: '请输入系统盘大小',
                  description: '系统盘的大小，单位为GB'
                }
              ]
            }
          ];
        default:
          return commonComponents;
      }
    }
  }
}
</script>

<style scoped>
.component-selector-container {
  margin-bottom: 2rem;
}

.component-selection {
  margin-top: 1rem;
}

.component-card {
  margin-top: 1rem;
}

.component-properties {
  margin-top: 1.5rem;
}
</style>
