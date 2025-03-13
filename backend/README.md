# 多云着陆区部署平台 - Go后端

这是一个使用Go语言和Gin框架实现的多云着陆区部署平台后端。该平台支持AWS、Azure、阿里云、百度云、华为云、腾讯云和火山云等多种云环境的部署。

## 项目结构

```
go-backend/
├── controllers/     # API控制器
├── models/          # 数据模型
├── routes/          # 路由配置
├── utils/           # 工具函数
├── terraform/       # Terraform配置和部署
├── main.go          # 主程序入口
└── .env             # 环境变量配置
```

## 功能特性

1. 支持7种云服务提供商：AWS、Azure、阿里云、百度云、华为云、腾讯云和火山云
2. 提供云服务提供商、区域和可用区的选择
3. 支持VPC和子网配置
4. 支持云组件选择和配置
5. 使用Terraform进行基础设施部署
6. 提供部署状态监控和反馈
7. 生成资源拓扑图

## API接口

1. **获取云服务提供商列表**
   - 路径: `/api/providers`
   - 方法: GET
   - 功能: 返回支持的云服务提供商列表

2. **获取区域列表**
   - 路径: `/api/regions/:provider`
   - 方法: GET
   - 参数: provider - 云服务提供商标识
   - 功能: 返回指定云服务提供商的可用区域列表

3. **获取可用区列表**
   - 路径: `/api/azs/:provider/:region`
   - 方法: GET
   - 参数: provider - 云服务提供商标识, region - 区域标识
   - 功能: 返回指定云服务提供商和区域的可用区列表

4. **获取云组件列表**
   - 路径: `/api/components/:provider/:region`
   - 方法: GET
   - 参数: provider - 云服务提供商标识, region - 区域标识
   - 功能: 返回指定云服务提供商和区域可用的云组件列表

5. **执行部署**
   - 路径: `/api/deploy`
   - 方法: POST
   - 参数: 部署配置对象(包含云提供商、区域、可用区、VPC、子网和组件信息)
   - 功能: 异步执行部署过程

6. **获取部署状态**
   - 路径: `/api/deployment/status`
   - 方法: GET
   - 功能: 获取当前部署的状态信息

## 安装和运行

### 前提条件

- Go 1.18或更高版本
- Terraform CLI

### 安装步骤

1. 克隆代码库
   ```bash
   git clone <repository-url>
   cd go-backend
   ```

2. 安装依赖
   ```bash
   go mod tidy
   ```

3. 配置环境变量
   创建`.env`文件并设置以下变量：
   ```
   PORT=3000
   GIN_MODE=release  # 或debug用于开发环境
   ```

4. 构建和运行
   ```bash
   go build -o backend
   ./backend
   ```

## 部署配置示例

以下是一个部署配置的JSON示例：

```json
{
  "cloudProvider": "aws",
  "region": "us-east-1",
  "az": "us-east-1a",
  "vpc": {
    "name": "my-vpc",
    "cidr": "10.0.0.0/16",
    "enableDnsSupport": true,
    "enableDnsHostnames": true
  },
  "subnet": {
    "name": "my-subnet",
    "cidr": "10.0.1.0/24",
    "mapPublicIpOnLaunch": false
  },
  "components": [
    {
      "name": "负载均衡器",
      "value": "load-balancer"
    },
    {
      "name": "对象存储",
      "value": "object-storage"
    }
  ],
  "componentProperties": {
    "load-balancer": {
      "instance_count": "1",
      "listener_port": "80"
    },
    "object-storage": {
      "bucket_name": "my-unique-bucket",
      "storage_class": "Standard"
    }
  }
}
```

## 与原Express后端的区别

本项目是原Express后端的Go语言重写版本，保持了相同的API接口和功能，但使用了Go语言和Gin框架的特性进行了优化：

1. 使用Go的并发特性处理异步部署
2. 使用Go的类型系统提供更强的类型安全
3. 使用Gin框架提供更高的性能
4. 保持了与原Express后端相同的API接口，便于前端无缝切换

## 注意事项

- 确保Terraform已正确安装并配置
- 确保有适当的云服务提供商凭证配置
- 部署过程是异步的，需要通过状态API监控进度
