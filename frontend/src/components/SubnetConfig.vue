<template>
  <div class="subnet-config-container">
    <el-card shadow="hover">
      <template #header>
        <h2>配置子网</h2>
      </template>
      <div class="subnet-form">
        <el-alert
          type="info"
          :closable="false"
          title="所属VPC信息"
        >
          <el-descriptions :column="1" border>
            <el-descriptions-item label="VPC名称">{{ vpcName }}</el-descriptions-item>
            <el-descriptions-item label="VPC CIDR">{{ vpcCIDR }}</el-descriptions-item>
          </el-descriptions>
        </el-alert>
        
        <el-form :model="subnetForm" label-position="top" :rules="rules" ref="subnetFormRef" class="mt-4">
          <el-form-item label="子网名称" prop="subnetName">
            <el-input 
              v-model="subnetName" 
              placeholder="请输入子网名称"
              @input="validateForm"
            />
          </el-form-item>
          
          <el-form-item label="CIDR块" prop="subnetCIDR">
            <el-input 
              v-model="subnetCIDR" 
              placeholder="例如: 10.0.1.0/24"
              @input="validateForm"
            />
            <el-text type="info" size="small">子网CIDR必须是VPC CIDR的子集</el-text>
          </el-form-item>
          
          <el-divider content-position="left">子网选项</el-divider>
          
          <el-form-item>
            <el-switch
              v-model="mapPublicIpOnLaunch"
              active-text="自动分配公网IP"
              @change="validateForm"
            />
            <div class="mt-2">
              <el-text type="info" size="small">在此子网中启动的实例将自动分配公网IP地址</el-text>
            </div>
          </el-form-item>
          
          <el-form-item label="可用区">
            <el-tag size="large" type="success">{{ azName }}</el-tag>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
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
      isValid: false,
      subnetForm: {},
      rules: {
        subnetName: [
          { required: true, message: '子网名称不能为空', trigger: 'blur' },
          { min: 2, message: '子网名称至少需要2个字符', trigger: 'blur' },
          { max: 64, message: '子网名称不能超过64个字符', trigger: 'blur' }
        ],
        subnetCIDR: [
          { required: true, message: 'CIDR块不能为空', trigger: 'blur' },
          { 
            pattern: /^([0-9]{1,3}\.){3}[0-9]{1,3}\/([0-9]|[1-2][0-9]|3[0-2])$/, 
            message: '请输入有效的CIDR格式，例如: 10.0.1.0/24', 
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
  margin-top: 1rem;
}

.mt-2 {
  margin-top: 0.5rem;
}

.mt-4 {
  margin-top: 1rem;
}
</style>
