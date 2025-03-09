<template>
  <div class="vpc-config-container">
    <h2>配置VPC</h2>
    <div class="vpc-form">
      <div class="form-group">
        <label for="vpc-name">VPC名称</label>
        <input 
          id="vpc-name" 
          type="text" 
          v-model="vpcName" 
          placeholder="请输入VPC名称"
          @input="validateForm"
        />
        <div v-if="errors.vpcName" class="error-message">{{ errors.vpcName }}</div>
      </div>
      <div class="form-group">
        <label for="vpc-cidr">CIDR块</label>
        <input 
          id="vpc-cidr" 
          type="text" 
          v-model="vpcCIDR" 
          placeholder="例如: 10.0.0.0/16"
          @input="validateForm"
        />
        <div v-if="errors.vpcCIDR" class="error-message">{{ errors.vpcCIDR }}</div>
        <div class="helper-text">推荐使用私有IP地址范围，如10.0.0.0/16、172.16.0.0/12或192.168.0.0/16</div>
      </div>
      <div class="vpc-options">
        <h3>VPC选项</h3>
        <div class="option-group">
          <label class="checkbox-label">
            <input type="checkbox" v-model="enableDnsSupport">
            启用DNS支持
          </label>
          <div class="option-description">允许VPC中的实例通过VPC提供的DNS服务器解析DNS</div>
        </div>
        <div class="option-group">
          <label class="checkbox-label">
            <input type="checkbox" v-model="enableDnsHostnames">
            启用DNS主机名
          </label>
          <div class="option-description">为VPC中的实例分配DNS主机名</div>
        </div>
      </div>
    </div>
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
      isValid: false
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
  margin-top: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.error-message {
  color: #e74c3c;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}

.helper-text {
  color: #7f8c8d;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}

.vpc-options {
  margin-top: 2rem;
  padding: 1.5rem;
  background-color: #f9f9f9;
  border-radius: 8px;
}

.vpc-options h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  font-size: 1.1rem;
}

.option-group {
  margin-bottom: 1rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  font-weight: bold;
  cursor: pointer;
}

.checkbox-label input {
  margin-right: 0.5rem;
}

.option-description {
  margin-top: 0.25rem;
  margin-left: 1.5rem;
  font-size: 0.9rem;
  color: #7f8c8d;
}
</style>
