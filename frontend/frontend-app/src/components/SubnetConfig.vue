<template>
  <div class="subnet-config-container">
    <h2>配置子网</h2>
    <div class="subnet-form">
      <div class="vpc-info">
        <h3>所属VPC信息</h3>
        <div class="vpc-details">
          <div class="detail-item">
            <span class="detail-label">VPC名称:</span>
            <span class="detail-value">{{ vpcName }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">VPC CIDR:</span>
            <span class="detail-value">{{ vpcCIDR }}</span>
          </div>
        </div>
      </div>
      
      <div class="form-group">
        <label for="subnet-name">子网名称</label>
        <input 
          id="subnet-name" 
          type="text" 
          v-model="subnetName" 
          placeholder="请输入子网名称"
          @input="validateForm"
        />
        <div v-if="errors.subnetName" class="error-message">{{ errors.subnetName }}</div>
      </div>
      
      <div class="form-group">
        <label for="subnet-cidr">CIDR块</label>
        <input 
          id="subnet-cidr" 
          type="text" 
          v-model="subnetCIDR" 
          placeholder="例如: 10.0.1.0/24"
          @input="validateForm"
        />
        <div v-if="errors.subnetCIDR" class="error-message">{{ errors.subnetCIDR }}</div>
        <div class="helper-text">子网CIDR必须是VPC CIDR的子集</div>
      </div>
      
      <div class="subnet-options">
        <h3>子网选项</h3>
        <div class="option-group">
          <label class="checkbox-label">
            <input type="checkbox" v-model="mapPublicIpOnLaunch">
            自动分配公网IP
          </label>
          <div class="option-description">在此子网中启动的实例将自动分配公网IP地址</div>
        </div>
        
        <div class="az-selection">
          <label>可用区</label>
          <div class="selected-az">
            <span class="az-label">已选择可用区:</span>
            <span class="az-value">{{ azName }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SubnetConfig',
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
    azName: {
      type: String,
      required: true
    },
    vpcName: {
      type: String,
      required: true
    },
    vpcCIDR: {
      type: String,
      required: true
    },
    initialSubnetName: {
      type: String,
      default: ''
    },
    initialSubnetCIDR: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      subnetName: this.initialSubnetName || '',
      subnetCIDR: this.initialSubnetCIDR || '',
      mapPublicIpOnLaunch: true,
      errors: {
        subnetName: '',
        subnetCIDR: ''
      },
      isValid: false
    }
  },
  watch: {
    cloudProvider() {
      // 根据云服务提供商设置默认值
      this.setDefaultsByProvider();
    },
    vpcCIDR() {
      // 当VPC CIDR变化时，重新生成默认子网CIDR
      if (!this.initialSubnetCIDR) {
        this.generateDefaultSubnetCIDR();
      }
      this.validateForm();
    },
    az() {
      // 当可用区变化时，更新子网名称
      this.setDefaultsByProvider();
    },
    initialSubnetName(newVal) {
      this.subnetName = newVal;
      this.validateForm();
    },
    initialSubnetCIDR(newVal) {
      this.subnetCIDR = newVal;
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
      const azSuffix = this.az.split('-').pop();
      
      switch(this.cloudProvider) {
        case 'aws':
          this.subnetName = this.subnetName || `AWS-Subnet-${this.region}-${azSuffix}`;
          break;
        case 'azure':
          this.subnetName = this.subnetName || `Azure-Subnet-${this.region}-${azSuffix}`;
          break;
        case 'alicloud':
          this.subnetName = this.subnetName || `Alicloud-Subnet-${this.region}-${azSuffix}`;
          break;
        case 'baidu':
          this.subnetName = this.subnetName || `Baidu-Subnet-${this.region}-${azSuffix}`;
          break;
        case 'huawei':
          this.subnetName = this.subnetName || `Huawei-Subnet-${this.region}-${azSuffix}`;
          break;
        case 'tencent':
          this.subnetName = this.subnetName || `Tencent-Subnet-${this.region}-${azSuffix}`;
          break;
        case 'volcengine':
          this.subnetName = this.subnetName || `Volcengine-Subnet-${this.region}-${azSuffix}`;
          break;
        default:
          this.subnetName = this.subnetName || `Subnet-${this.region}-${azSuffix}`;
      }
      
      if (!this.subnetCIDR) {
        this.generateDefaultSubnetCIDR();
      }
      
      this.validateForm();
    },
    generateDefaultSubnetCIDR() {
      // 根据VPC CIDR生成默认子网CIDR
      if (!this.vpcCIDR) return;
      
      try {
        // 简单的CIDR计算，将VPC CIDR的掩码增加8位
        const parts = this.vpcCIDR.split('/');
        const ipParts = parts[0].split('.');
        const mask = parseInt(parts[1]);
        
        if (mask <= 24) {
          // 如果VPC掩码小于等于24，则子网掩码为VPC掩码+8
          this.subnetCIDR = `${ipParts[0]}.${ipParts[1]}.${ipParts[2]}.0/${mask + 8}`;
        } else {
          // 如果VPC掩码大于24，则子网掩码与VPC掩码相同
          this.subnetCIDR = this.vpcCIDR;
        }
      } catch (e) {
        console.error('Error generating default subnet CIDR:', e);
        this.subnetCIDR = '';
      }
    },
    validateForm() {
      // 验证子网名称
      if (!this.subnetName) {
        this.errors.subnetName = '子网名称不能为空';
      } else if (this.subnetName.length < 2) {
        this.errors.subnetName = '子网名称至少需要2个字符';
      } else if (this.subnetName.length > 64) {
        this.errors.subnetName = '子网名称不能超过64个字符';
      } else {
        this.errors.subnetName = '';
      }
      
      // 验证CIDR块
      const cidrPattern = /^([0-9]{1,3}\.){3}[0-9]{1,3}\/([0-9]|[1-2][0-9]|3[0-2])$/;
      if (!this.subnetCIDR) {
        this.errors.subnetCIDR = 'CIDR块不能为空';
      } else if (!cidrPattern.test(this.subnetCIDR)) {
        this.errors.subnetCIDR = '请输入有效的CIDR格式，例如: 10.0.1.0/24';
      } else {
        // 验证IP地址部分是否有效
        const ipParts = this.subnetCIDR.split('/')[0].split('.');
        const isValidIp = ipParts.every(part => parseInt(part) >= 0 && parseInt(part) <= 255);
        
        if (!isValidIp) {
          this.errors.subnetCIDR = 'IP地址部分无效，每个部分应在0-255之间';
        } else if (!this.isSubnetOfVpc()) {
          this.errors.subnetCIDR = '子网CIDR必须是VPC CIDR的子集';
        } else {
          this.errors.subnetCIDR = '';
        }
      }
      
      // 检查表单是否有效
      this.isValid = !this.errors.subnetName && !this.errors.subnetCIDR;
      
      // 如果表单有效，发出事件
      if (this.isValid) {
        this.$emit('subnet-updated', {
          name: this.subnetName,
          cidr: this.subnetCIDR,
          az: this.az,
          mapPublicIpOnLaunch: this.mapPublicIpOnLaunch
        });
      }
    },
    isSubnetOfVpc() {
      // 简单检查子网CIDR是否是VPC CIDR的子集
      // 注意：这是一个简化的检查，实际应用中可能需要更复杂的IP地址计算
      try {
        const vpcParts = this.vpcCIDR.split('/');
        const subnetParts = this.subnetCIDR.split('/');
        
        const vpcIp = vpcParts[0].split('.').map(Number);
        const subnetIp = subnetParts[0].split('.').map(Number);
        
        const vpcMask = parseInt(vpcParts[1]);
        const subnetMask = parseInt(subnetParts[1]);
        
        // 子网掩码必须大于等于VPC掩码
        if (subnetMask < vpcMask) {
          return false;
        }
        
        // 检查子网IP是否在VPC IP范围内
        // 这是一个简化的检查，只比较前缀
        const prefixBytes = Math.floor(vpcMask / 8);
        
        for (let i = 0; i < prefixBytes; i++) {
          if (vpcIp[i] !== subnetIp[i]) {
            return false;
          }
        }
        
        // 检查部分字节
        if (vpcMask % 8 !== 0) {
          const bitMask = 256 - (1 << (8 - (vpcMask % 8)));
          const vpcPartialByte = vpcIp[prefixBytes] & bitMask;
          const subnetPartialByte = subnetIp[prefixBytes] & bitMask;
          
          if (vpcPartialByte !== subnetPartialByte) {
            return false;
          }
        }
        
        return true;
      } catch (e) {
        console.error('Error checking if subnet is part of VPC:', e);
        return false;
      }
    }
  }
}
</script>

<style scoped>
.subnet-config-container {
  margin-bottom: 2rem;
}

.subnet-form {
  margin-top: 1.5rem;
}

.vpc-info {
  margin-bottom: 2rem;
  padding: 1.5rem;
  background-color: #f5f7fa;
  border-radius: 8px;
  border-left: 4px solid #3498db;
}

.vpc-info h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  font-size: 1.1rem;
  color: #2c3e50;
}

.vpc-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.detail-item {
  display: flex;
}

.detail-label {
  width: 100px;
  font-weight: bold;
  color: #7f8c8d;
}

.detail-value {
  flex: 1;
  font-family: monospace;
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

.subnet-options {
  margin-top: 2rem;
  padding: 1.5rem;
  background-color: #f9f9f9;
  border-radius: 8px;
}

.subnet-options h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  font-size: 1.1rem;
}

.option-group {
  margin-bottom: 1.5rem;
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

.az-selection {
  margin-top: 1rem;
}

.az-selection label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

.selected-az {
  padding: 0.75rem;
  background-color: #edf2f7;
  border-radius: 4px;
}

.az-label {
  font-weight: bold;
  margin-right: 0.5rem;
}

.az-value {
  font-family: monospace;
}
</style>
