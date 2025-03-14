<template>
  <div class="cloud-component-config-container">
    <el-card shadow="hover">
      <template #header>
        <h2>云组件配置</h2>
      </template>
      <div class="component-config">
        <el-skeleton v-if="loading" :rows="5" animated />
        
        <el-empty v-else-if="selectedComponents.length === 0" description="没有选择任何云组件，请返回上一步选择云组件" />
        
        <div v-else class="component-config-list">
          <el-tabs v-model="activeTab" type="card">
            <!-- 对象存储桶策略配置 -->
            <el-tab-pane v-if="hasObjectStorage" label="对象存储桶策略" name="object-storage">
              <el-card shadow="hover" class="config-card">
                <el-form label-position="top">
                  <el-form-item label="存储桶策略类型">
                    <el-select v-model="bucketPolicyType" placeholder="请选择策略类型">
                      <el-option label="公共读取" value="public-read" />
                      <el-option label="公共读写" value="public-read-write" />
                      <el-option label="私有" value="private" />
                      <el-option label="自定义" value="custom" />
                    </el-select>
                  </el-form-item>
                  
                  <el-form-item v-if="bucketPolicyType === 'custom'" label="自定义策略 (JSON格式)">
                    <el-input
                      type="textarea"
                      v-model="customBucketPolicy"
                      :rows="10"
                      placeholder="请输入JSON格式的存储桶策略"
                    />
                  </el-form-item>
                  
                  <el-form-item label="存储桶生命周期规则">
                    <el-switch v-model="enableLifecycleRules" />
                  </el-form-item>
                  
                  <template v-if="enableLifecycleRules">
                    <el-divider content-position="left">生命周期规则配置</el-divider>
                    
                    <el-form-item label="规则名称">
                      <el-input v-model="lifecycleRule.name" placeholder="请输入规则名称" />
                    </el-form-item>
                    
                    <el-form-item label="状态">
                      <el-radio-group v-model="lifecycleRule.status">
                        <el-radio label="Enabled">启用</el-radio>
                        <el-radio label="Disabled">禁用</el-radio>
                      </el-radio-group>
                    </el-form-item>
                    
                    <el-form-item label="过期天数">
                      <el-input-number v-model="lifecycleRule.expirationDays" :min="1" :max="3650" />
                      <el-text type="info" size="small">对象创建后多少天过期</el-text>
                    </el-form-item>
                    
                    <el-form-item label="转换为低频存储天数">
                      <el-input-number v-model="lifecycleRule.transitionDays" :min="1" :max="3650" />
                      <el-text type="info" size="small">对象创建后多少天转换为低频存储</el-text>
                    </el-form-item>
                  </template>
                </el-form>
              </el-card>
            </el-tab-pane>
            
            <!-- Transit Gateway配置 -->
            <el-tab-pane v-if="hasTransitGateway" label="Transit Gateway配置" name="transit-gateway">
              <el-card shadow="hover" class="config-card">
                <el-form label-position="top">
                  <el-form-item label="路由表配置">
                    <el-switch v-model="enableRouteTables" />
                  </el-form-item>
                  
                  <template v-if="enableRouteTables">
                    <el-divider content-position="left">路由表配置</el-divider>
                    
                    <el-form-item label="路由表名称">
                      <el-input v-model="transitGatewayConfig.routeTableName" placeholder="请输入路由表名称" />
                    </el-form-item>
                    
                    <el-form-item label="默认路由表">
                      <el-switch v-model="transitGatewayConfig.defaultRouteTable" />
                    </el-form-item>
                  </template>
                  
                  <el-form-item label="VPC附件">
                    <el-switch v-model="enableVpcAttachment" />
                  </el-form-item>
                  
                  <template v-if="enableVpcAttachment">
                    <el-divider content-position="left">VPC附件配置</el-divider>
                    
                    <el-form-item label="子网IDs">
                      <el-input v-model="transitGatewayConfig.subnetIds" placeholder="请输入子网IDs，多个ID用逗号分隔" />
                    </el-form-item>
                    
                    <el-form-item label="DNS支持">
                      <el-switch v-model="transitGatewayConfig.dnsSupport" />
                    </el-form-item>
                    
                    <el-form-item label="IPv6支持">
                      <el-switch v-model="transitGatewayConfig.ipv6Support" />
                    </el-form-item>
                  </template>
                </el-form>
              </el-card>
            </el-tab-pane>
            
            <!-- 其他组件配置标签页可以在这里添加 -->
            <el-tab-pane v-if="hasOtherComponents" label="其他组件配置" name="other">
              <el-card shadow="hover" class="config-card">
                <el-empty description="其他组件的高级配置将在后续版本中提供" />
              </el-card>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'CloudComponentConfig',
  props: {
    cloudProvider: {
      type: String,
      required: true
    },
    region: {
      type: String,
      required: true
    },
    selectedComponents: {
      type: Array,
      default: () => []
    },
    componentProperties: {
      type: Object,
      default: () => ({})
    },
    initialConfig: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      loading: false,
      activeTab: 'object-storage',
      
      // 对象存储桶策略配置
      bucketPolicyType: this.initialConfig.bucketPolicyType || 'private',
      customBucketPolicy: this.initialConfig.customBucketPolicy || '',
      enableLifecycleRules: this.initialConfig.enableLifecycleRules || false,
      lifecycleRule: this.initialConfig.lifecycleRule || {
        name: 'default-lifecycle-rule',
        status: 'Enabled',
        expirationDays: 365,
        transitionDays: 30
      },
      
      // Transit Gateway配置
      enableRouteTables: this.initialConfig.enableRouteTables || false,
      enableVpcAttachment: this.initialConfig.enableVpcAttachment || false,
      transitGatewayConfig: this.initialConfig.transitGatewayConfig || {
        routeTableName: 'main-route-table',
        defaultRouteTable: true,
        subnetIds: '',
        dnsSupport: true,
        ipv6Support: false
      }
    }
  },
  computed: {
    hasObjectStorage() {
      return this.selectedComponents.includes('object-storage');
    },
    hasTransitGateway() {
      return this.selectedComponents.includes('transit-gateway');
    },
    hasOtherComponents() {
      return this.selectedComponents.length > 0 && 
             !this.selectedComponents.every(comp => 
               ['object-storage', 'transit-gateway'].includes(comp)
             );
    },
    componentConfig() {
      return {
        bucketPolicyType: this.bucketPolicyType,
        customBucketPolicy: this.customBucketPolicy,
        enableLifecycleRules: this.enableLifecycleRules,
        lifecycleRule: this.lifecycleRule,
        enableRouteTables: this.enableRouteTables,
        enableVpcAttachment: this.enableVpcAttachment,
        transitGatewayConfig: this.transitGatewayConfig
      };
    }
  },
  watch: {
    selectedComponents() {
      this.updateActiveTab();
    },
    componentConfig: {
      deep: true,
      handler() {
        this.emitConfigUpdate();
      }
    },
    initialConfig: {
      deep: true,
      handler(newVal) {
        if (newVal) {
          this.bucketPolicyType = newVal.bucketPolicyType || 'private';
          this.customBucketPolicy = newVal.customBucketPolicy || '';
          this.enableLifecycleRules = newVal.enableLifecycleRules || false;
          this.lifecycleRule = newVal.lifecycleRule || {
            name: 'default-lifecycle-rule',
            status: 'Enabled',
            expirationDays: 365,
            transitionDays: 30
          };
          this.enableRouteTables = newVal.enableRouteTables || false;
          this.enableVpcAttachment = newVal.enableVpcAttachment || false;
          this.transitGatewayConfig = newVal.transitGatewayConfig || {
            routeTableName: 'main-route-table',
            defaultRouteTable: true,
            subnetIds: '',
            dnsSupport: true,
            ipv6Support: false
          };
        }
      }
    }
  },
  mounted() {
    this.updateActiveTab();
  },
  methods: {
    updateActiveTab() {
      if (this.hasObjectStorage) {
        this.activeTab = 'object-storage';
      } else if (this.hasTransitGateway) {
        this.activeTab = 'transit-gateway';
      } else if (this.hasOtherComponents) {
        this.activeTab = 'other';
      }
    },
    emitConfigUpdate() {
      this.$emit('config-updated', this.componentConfig);
    }
  }
}
</script>

<style scoped>
.cloud-component-config-container {
  margin-bottom: 2rem;
}

.component-config {
  margin-top: 1rem;
}

.config-card {
  margin-top: 1rem;
}

.el-divider {
  margin: 16px 0;
}
</style>
