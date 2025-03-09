const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');
const fs = require('fs');
const path = require('path');
const { exec } = require('child_process');
const util = require('util');

const app = express();
const port = 3000;

// 中间件
app.use(cors());
app.use(bodyParser.json());

// 存储部署状态的对象
const deploymentStatus = {
  status: 'idle', // idle, preparing, deploying, completed, failed
  progress: 0,
  message: '',
  logs: [],
  result: null,
  topology: null
};

// 工具函数 - 执行shell命令
const execPromise = util.promisify(exec);
async function runCommand(command, workingDir) {
  try {
    const { stdout, stderr } = await execPromise(command, { cwd: workingDir });
    return { success: true, stdout, stderr };
  } catch (error) {
    return { success: false, error: error.message };
  }
}

// 工具函数 - 写入文件
function writeToFile(filePath, content) {
  return new Promise((resolve, reject) => {
    fs.writeFile(filePath, content, 'utf8', (err) => {
      if (err) {
        reject(err);
      } else {
        resolve();
      }
    });
  });
}

// 工具函数 - 读取文件
function readFile(filePath) {
  return new Promise((resolve, reject) => {
    fs.readFile(filePath, 'utf8', (err, data) => {
      if (err) {
        reject(err);
      } else {
        resolve(data);
      }
    });
  });
}

// 工具函数 - 检查目录是否存在，不存在则创建
function ensureDirectoryExists(directory) {
  if (!fs.existsSync(directory)) {
    fs.mkdirSync(directory, { recursive: true });
  }
}

// 路由 - 获取云服务提供商列表
app.get('/api/providers', (req, res) => {
  const providers = [
    { name: 'AWS', value: 'aws', logo: 'AWS' },
    { name: 'Azure', value: 'azure', logo: 'Azure' },
    { name: '阿里云', value: 'alicloud', logo: '阿里云' },
    { name: '百度云', value: 'baidu', logo: '百度云' },
    { name: '华为云', value: 'huawei', logo: '华为云' },
    { name: '腾讯云', value: 'tencent', logo: '腾讯云' },
    { name: '火山云', value: 'volcengine', logo: '火山云' }
  ];
  
  res.json({ success: true, data: providers });
});

// 路由 - 获取区域列表
app.get('/api/regions/:provider', (req, res) => {
  const { provider } = req.params;
  
  // 从Terraform文件中提取区域信息
  // 这里使用模拟数据，实际应用中应该解析Terraform文件
  let regions = [];
  
  switch(provider) {
    case 'aws':
      regions = [
        { name: '美国东部 (弗吉尼亚北部)', value: 'us-east-1' },
        { name: '美国东部 (俄亥俄)', value: 'us-east-2' },
        { name: '美国西部 (加利福尼亚北部)', value: 'us-west-1' },
        { name: '美国西部 (俄勒冈)', value: 'us-west-2' },
        { name: '亚太地区 (香港)', value: 'ap-east-1' },
        { name: '亚太地区 (东京)', value: 'ap-northeast-1' }
      ];
      break;
    case 'azure':
      regions = [
        { name: '美国东部', value: 'eastus' },
        { name: '美国东部2', value: 'eastus2' },
        { name: '美国西部', value: 'westus' },
        { name: '美国西部2', value: 'westus2' },
        { name: '东亚', value: 'eastasia' },
        { name: '东南亚', value: 'southeastasia' }
      ];
      break;
    case 'alicloud':
      regions = [
        { name: '华北 1 (青岛)', value: 'cn-qingdao' },
        { name: '华北 2 (北京)', value: 'cn-beijing' },
        { name: '华北 3 (张家口)', value: 'cn-zhangjiakou' },
        { name: '华东 1 (杭州)', value: 'cn-hangzhou' },
        { name: '华东 2 (上海)', value: 'cn-shanghai' },
        { name: '华南 1 (深圳)', value: 'cn-shenzhen' }
      ];
      break;
    case 'baidu':
      regions = [
        { name: '华北-北京', value: 'bj' },
        { name: '华南-广州', value: 'gz' },
        { name: '华东-苏州', value: 'su' }
      ];
      break;
    case 'huawei':
      regions = [
        { name: '华北-北京一', value: 'cn-north-1' },
        { name: '华北-北京四', value: 'cn-north-4' },
        { name: '华东-上海一', value: 'cn-east-3' },
        { name: '华南-广州', value: 'cn-south-1' },
        { name: '亚太-香港', value: 'ap-southeast-1' }
      ];
      break;
    case 'tencent':
      regions = [
        { name: '华南地区(广州)', value: 'ap-guangzhou' },
        { name: '华东地区(上海)', value: 'ap-shanghai' },
        { name: '华北地区(北京)', value: 'ap-beijing' },
        { name: '西南地区(成都)', value: 'ap-chengdu' },
        { name: '西南地区(重庆)', value: 'ap-chongqing' },
        { name: '港澳台地区(中国香港)', value: 'ap-hongkong' }
      ];
      break;
    case 'volcengine':
      regions = [
        { name: '华北-北京', value: 'cn-beijing' },
        { name: '华东-上海', value: 'cn-shanghai' },
        { name: '华南-广州', value: 'cn-guangzhou' }
      ];
      break;
    default:
      regions = [];
  }
  
  res.json({ success: true, data: regions });
});

// 路由 - 获取可用区列表
app.get('/api/azs/:provider/:region', (req, res) => {
  const { provider, region } = req.params;
  
  // 从Terraform文件中提取可用区信息
  // 这里使用模拟数据，实际应用中应该解析Terraform文件
  let azs = [];
  
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
  
  // 如果有特定区域的数据，使用它，否则生成通用的可用区
  if (mockAZs[region]) {
    azs = mockAZs[region];
  } else {
    azs = [
      { name: '可用区A', value: `${region}-a` },
      { name: '可用区B', value: `${region}-b` },
      { name: '可用区C', value: `${region}-c` }
    ];
  }
  
  res.json({ success: true, data: azs });
});

// 路由 - 获取云组件列表
app.get('/api/components/:provider/:region', (req, res) => {
  const { provider, region } = req.params;
  
  // 从Terraform文件中提取组件信息
  // 这里使用模拟数据，实际应用中应该解析Terraform文件
  
  // 通用组件
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
  
  const components = [...commonComponents, ...providerSpecificComponents];
  
  res.json({ success: true, data: components });
});

// 路由 - 部署
app.post('/api/deploy', async (req, res) => {
  const deploymentConfig = req.body;
  
  // 更新部署状态
  deploymentStatus.status = 'preparing';
  deploymentStatus.progress = 0;
  deploymentStatus.message = '正在准备部署资源...';
  deploymentStatus.logs = ['开始部署过程...'];
  deploymentStatus.result = null;
  deploymentStatus.topology = null;
  
  // 异步处理部署
  processDeploy(deploymentConfig);
  
  // 立即返回响应，不等待部署完成
  res.json({ 
    success: true, 
    message: '部署已开始',
    deploymentId: new Date().getTime().toString()
  });
});

// 路由 - 获取部署状态
app.get('/api/deployment/status', (req, res) => {
  res.json({ success: true, data: deploymentStatus });
});

// 异步处理部署
async function processDeploy(config) {
  try {
    // 提取配置信息
    const { cloudProvider, region, az, vpc, subnet, components, componentProperties } = config;
    
    // 创建部署工作目录
    const deploymentId = new Date().getTime().toString();
    const workDir = path.join(__dirname, 'terraform', 'deployments', deploymentId);
    ensureDirectoryExists(workDir);
    
    // 更新状态
    deploymentStatus.logs.push(`创建部署工作目录: ${workDir}`);
    deploymentStatus.progress = 10;
    deploymentStatus.message = '正在生成Terraform配置...';
    
    // 生成Terraform配置文件
    const terraformConfig = generateTerraformConfig(config);
    await writeToFile(path.join(workDir, 'main.tf'), terraformConfig);
    
    deploymentStatus.logs.push('生成Terraform配置文件完成');
    deploymentStatus.progress = 20;
    deploymentStatus.message = '正在初始化Terraform...';
    
    // 初始化Terraform
    const initResult = await runCommand('terraform init', workDir);
    if (!initResult.success) {
      throw new Error(`Terraform初始化失败: ${initResult.error}`);
    }
    
    deploymentStatus.logs.push('Terraform初始化完成');
    deploymentStatus.logs.push(initResult.stdout);
    deploymentStatus.progress = 30;
    deploymentStatus.message = '正在验证Terraform配置...';
    
    // 验证Terraform配置
    const validateResult = await runCommand('terraform validate', workDir);
    if (!validateResult.success) {
      throw new Error(`Terraform配置验证失败: ${validateResult.error}`);
    }
    
    deploymentStatus.logs.push('Terraform配置验证通过');
    deploymentStatus.logs.push(validateResult.stdout);
    deploymentStatus.progress = 40;
    deploymentStatus.message = '正在生成Terraform执行计划...';
    
    // 生成执行计划
    const planResult = await runCommand('terraform plan -out=tfplan', workDir);
    if (!planResult.success) {
      throw new Error(`Terraform计划生成失败: ${planResult.error}`);
    }
    
    deploymentStatus.logs.push('Terraform执行计划生成完成');
    deploymentStatus.logs.push(planResult.stdout);
    deploymentStatus.progress = 60;
    deploymentStatus.message = '正在执行Terraform部署...';
    
    // 执行部署
    const applyResult = await runCommand('terraform apply -auto-approve tfplan', workDir);
    if (!applyResult.success) {
      throw new Error(`Terraform部署失败: ${applyResult.error}`);
    }
    
    deploymentStatus.logs.push('Terraform部署执行完成');
    deploymentStatus.logs.push(applyResult.stdout);
    deploymentStatus.progress = 90;
    deploymentStatus.message = '正在生成资源拓扑图...';
    
    // 生成拓扑图
    const topology = generateTopology(config);
    
    // 完成部署
    deploymentStatus.status = 'completed';
    deploymentStatus.progress = 100;
    deploymentStatus.message = '部署完成';
    deploymentStatus.result = {
      deploymentId,
      cloudProvider,
      region,
      az,
      vpc,
      subnet,
      components
    };
    deploymentStatus.topology = topology;
    
  } catch (error) {
    // 处理错误
    deploymentStatus.status = 'failed';
    deploymentStatus.message = `部署失败: ${error.message}`;
    deploymentStatus.logs.push(`错误: ${error.message}`);
    console.error('部署错误:', error);
  }
}

// 生成Terraform配置
function generateTerraformConfig(config) {
  const { cloudProvider, region, az, vpc, subnet, components, componentProperties } = config;
  
  let terraformConfig = '';
  
  // 添加提供商配置
  switch(cloudProvider) {
    case 'aws':
      terraformConfig += `provider "aws" {
  region = "${region}"
}

`;
      break;
    case 'azure':
      terraformConfig += `provider "azurerm" {
  features {}
}

`;
      break;
    case 'alicloud':
      terraformConfig += `provider "alicloud" {
  region = "${region}"
}

`;
      break;
    // 其他云服务提供商...
    default:
      terraformConfig += `provider "${cloudProvider}" {
  region = "${region}"
}

`;
  }
  
  // 添加VPC配置
  if (cloudProvider === 'aws') {
    terraformConfig += `resource "aws_vpc" "${vpc.name}" {
  cidr_block           = "${vpc.cidr}"
  enable_dns_support   = ${vpc.enableDnsSupport || true}
  enable_dns_hostnames = ${vpc.enableDnsHostnames || true}
  
  tags = {
    Name = "${vpc.name}"
  }
}

`;
  } else if (cloudProvider<response clipped><NOTE>To save on context only part of this file has been shown to you. You should retry this tool after you have searched inside the file with `grep -n` in order to find the line numbers of what you are looking for.</NOTE>