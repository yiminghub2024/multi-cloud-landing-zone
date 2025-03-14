<template>
  <div class="vpc-config-container">
    <el-card shadow="hover">
      <template #header>
        <h2>配置VPC</h2>
      </template>
      <div class="vpc-form">
        <el-form :model="vpcForm" label-position="top" :rules="rules" ref="vpcFormRef">
          <el-form-item label="VPC名称" prop="vpcName">
            <el-input 
              v-model="vpcName" 
              placeholder="请输入VPC名称"
              @input="validateForm"
            />
          </el-form-item>
          
          <el-form-item label="CIDR块" prop="vpcCIDR">
            <el-input 
              v-model="vpcCIDR" 
              placeholder="例如: 10.0.0.0/16"
              @input="validateForm"
            />
            <el-text type="info" size="small">推荐使用私有IP地址范围，如10.0.0.0/16、172.16.0.0/12或192.168.0.0/16</el-text>
          </el-form-item>
          
          <el-divider content-position="left">VPC选项</el-divider>
          
          <el-form-item>
            <el-space direction="vertical" alignment="start" :size="20" fill>
              <el-switch
                v-model="enableDnsSupport"
                active-text="启用DNS支持"
                @change="validateForm"
              />
              <el-text type="info" size="small">允许VPC中的实例通过VPC提供的DNS服务器解析DNS</el-text>
              
              <el-switch
                v-model="enableDnsHostnames"
                active-text="启用DNS主机名"
                @change="validateForm"
              />
              <el-text type="info" size="small">为VPC中的实例分配DNS主机名</el-text>
            </el-space>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'VPCConfig',
  props: {
    cloudProvider: {
      type: String,
      required: true
    },
    region: {
      type: String,
      required: true
    },
    initialVpcName: {
      type: String,
      default: ''
    },
    initialVpcCIDR: {
      type: String,
      default: '10.0.0.0/16'
    }
  },
  data() {
    return {
      vpcName: this.initialVpcName || '',
      vpcCIDR: this.initialVpcCIDR || '10.0.0.0/16',
      enableDnsSupport: true,
      enableDnsHostnames: true,
      errors: {
        vpcName: '',
        vpcCIDR: ''
      },
      isValid: false,
      vpcForm: {},
      rules: {
        vpcName: [
          { required: true, message: 'VPC名称不能为空', trigger: 'blur' },
          { min: 2, message: 'VPC名称至少需要2个字符', trigger: 'blur' },
          { max: 64, message: 'VPC名称不能超过64个字符', trigger: 'blur' }
        ],
        vpcCIDR: [
          { required: true, message: 'CIDR块不能为空', trigger: 'blur' },
          { 
            pattern: /^([0-9]{1,3}\.){3}[0-9]{1,3}\/([0-9]|[1-2][0-9]|3[0-2])$/, 
            message: '请输入有效的CIDR格式，例如: 10.0.0.0/16', 
            trigger: 'blur' 
          }
        ]
      }
    }
  },
  watch: {
    cloudProvider() {
      // 根据云服务提供商设置默认值
      this.setDefaultsByProvider();
    },
    initialVpcName(newVal) {
      this.vpcName = newVal;
      this.validateForm();
    },
    initialVpcCIDR(newVal) {
      this.vpcCIDR = newVal;
      this.validateForm();
    }
  },
  mounted() {
    this.setDefaultsByProvider();
    this.validateForm();
  },
  methods: {
    setDefaultsByProvider() {
      // 根据不同的云服务提供商设置默认值
      switch(this.cloudProvider) {
        case 'aws':
          this.vpcName = this.vpcName || `AWS-VPC-${this.region}`;
          break;
        case 'azure':
          this.vpcName = this.vpcName || `Azure-VNet-${this.region}`;
          break;
        case 'alicloud':
          this.vpcName = this.vpcName || `Alicloud-VPC-${this.region}`;
          break;
        case 'baidu':
          this.vpcName = this.vpcName || `Baidu-VPC-${this.region}`;
          break;
        case 'huawei':
          this.vpcName = this.vpcName || `Huawei-VPC-${this.region}`;
          break;
        case 'tencent':
          this.vpcName = this.vpcName || `Tencent-VPC-${this.region}`;
          break;
        case 'volcengine':
          this.vpcName = this.vpcName || `Volcengine-VPC-${this.region}`;
          break;
        default:
          this.vpcName = this.vpcName || `VPC-${this.region}`;
      }
      
      this.validateForm();
    },
    validateForm() {
      // 验证VPC名称
      if (!this.vpcName) {
        this.errors.vpcName = 'VPC名称不能为空';
      } else if (this.vpcName.length < 2) {
        this.errors.vpcName = 'VPC名称至少需要2个字符';
      } else if (this.vpcName.length > 64) {
        this.errors.vpcName = 'VPC名称不能超过64个字符';
      } else {
        this.errors.vpcName = '';
      }
      
      // 验证CIDR块
      const cidrPattern = /^([0-9]{1,3}\.){3}[0-9]{1,3}\/([0-9]|[1-2][0-9]|3[0-2])$/;
      if (!this.vpcCIDR) {
        this.errors.vpcCIDR = 'CIDR块不能为空';
      } else if (!cidrPattern.test(this.vpcCIDR)) {
        this.errors.vpcCIDR = '请输入有效的CIDR格式，例如: 10.0.0.0/16';
      } else {
        // 验证IP地址部分是否有效
        const ipParts = this.vpcCIDR.split('/')[0].split('.');
        const isValidIp = ipParts.every(part => parseInt(part) >= 0 && parseInt(part) <= 255);
        
        if (!isValidIp) {
          this.errors.vpcCIDR = 'IP地址部分无效，每个部分应在0-255之间';
        } else {
          this.errors.vpcCIDR = '';
        }
      }
      
      // 检查表单是否有效
      this.isValid = !this.errors.vpcName && !this.errors.vpcCIDR;
      
      // 如果表单有效，发出事件
      if (this.isValid) {
        this.$emit('vpc-updated', {
          name: this.vpcName,
          cidr: this.vpcCIDR,
          enableDnsSupport: this.enableDnsSupport,
          enableDnsHostnames: this.enableDnsHostnames
        });
      }
    }
  }
}
</script>

<style scoped>
.vpc-config-container {
  margin-bottom: 2rem;
}

.vpc-form {
  margin-top: 1rem;
}
</style>
