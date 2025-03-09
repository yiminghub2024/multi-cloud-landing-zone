# 多云着陆区部署平台部署文档

## 1. 部署概述

本文档提供了多云着陆区部署平台的部署指南，包括环境准备、安装步骤、配置说明和验证方法。

## 2. 系统要求

### 2.1 硬件要求

- CPU: 2核心或以上
- 内存: 4GB或以上
- 磁盘空间: 10GB或以上

### 2.2 软件要求

- 操作系统: Linux, macOS, 或 Windows
- Node.js: v14.0.0或以上
- npm: v6.0.0或以上
- Terraform: v1.0.0或以上

### 2.3 网络要求

- 互联网连接，用于下载依赖和访问云服务API
- 开放端口:
  - 前端: 8080 (开发环境) / 80 (生产环境)
  - 后端: 3000

## 3. 环境准备

### 3.1 安装Node.js和npm

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install nodejs npm
```

#### macOS
```bash
brew install node
```

#### Windows
从[Node.js官网](https://nodejs.org/)下载并安装最新LTS版本。

### 3.2 安装Terraform

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update && sudo apt-get install -y gnupg software-properties-common curl
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update && sudo apt-get install terraform
```

#### macOS
```bash
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

#### Windows
从[Terraform官网](https://www.terraform.io/downloads.html)下载并安装最新版本。

### 3.3 配置云服务提供商凭证

#### AWS
```bash
mkdir -p ~/.aws
cat > ~/.aws/credentials << EOF
[default]
aws_access_key_id = YOUR_ACCESS_KEY
aws_secret_access_key = YOUR_SECRET_KEY
EOF
```

#### Azure
```bash
az login
```

#### 阿里云
```bash
mkdir -p ~/.aliyun
cat > ~/.aliyun/config.json << EOF
{
  "current": "default",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "YOUR_ACCESS_KEY",
      "access_key_secret": "YOUR_SECRET_KEY"
    }
  ]
}
EOF
```

其他云服务提供商的凭证配置请参考各自的官方文档。

## 4. 安装步骤

### 4.1 获取源代码

```bash
git clone https://github.com/your-org/multi-cloud-landing-zone.git
cd multi-cloud-landing-zone
```

### 4.2 安装依赖

#### 安装前端依赖
```bash
cd frontend/frontend-app
npm install
```

#### 安装后端依赖
```bash
cd ../../backend
npm install
```

### 4.3 配置应用

#### 前端配置
创建或编辑 `frontend/frontend-app/.env` 文件:
```
VUE_APP_API_URL=http://localhost:3000/api
```

#### 后端配置
创建或编辑 `backend/.env` 文件:
```
PORT=3000
NODE_ENV=production
TERRAFORM_PATH=/usr/bin/terraform
```

## 5. 部署方式

### 5.1 开发环境部署

#### 启动后端服务
```bash
cd backend
npm run dev
```

#### 启动前端服务
```bash
cd ../frontend/frontend-app
npm run serve
```

访问 http://localhost:8080 即可使用应用。

### 5.2 生产环境部署

#### 构建前端
```bash
cd frontend/frontend-app
npm run build
```

#### 配置Nginx (可选)
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        root /path/to/multi-cloud-landing-zone/frontend/frontend-app/dist;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

#### 启动后端服务
```bash
cd backend
npm start
```

### 5.3 Docker部署

#### 构建Docker镜像
```bash
docker build -t multi-cloud-landing-zone .
```

#### 运行Docker容器
```bash
docker run -d -p 80:8080 -p 3000:3000 \
  -v ~/.aws:/root/.aws \
  -v ~/.azure:/root/.azure \
  -v ~/.aliyun:/root/.aliyun \
  --name multi-cloud-landing-zone \
  multi-cloud-landing-zone
```

### 5.4 Kubernetes部署

#### 创建Kubernetes配置文件
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: multi-cloud-landing-zone
spec:
  replicas: 1
  selector:
    matchLabels:
      app: multi-cloud-landing-zone
  template:
    metadata:
      labels:
        app: multi-cloud-landing-zone
    spec:
      containers:
      - name: multi-cloud-landing-zone
        image: multi-cloud-landing-zone:latest
        ports:
        - containerPort: 8080
        - containerPort: 3000
        volumeMounts:
        - name: aws-credentials
          mountPath: /root/.aws
        - name: azure-credentials
          mountPath: /root/.azure
        - name: aliyun-credentials
          mountPath: /root/.aliyun
      volumes:
      - name: aws-credentials
        secret:
          secretName: aws-credentials
      - name: azure-credentials
        secret:
          secretName: azure-credentials
      - name: aliyun-credentials
        secret:
          secretName: aliyun-credentials
---
apiVersion: v1
kind: Service
metadata:
  name: multi-cloud-landing-zone
spec:
  selector:
    app: multi-cloud-landing-zone
  ports:
  - name: frontend
    port: 80
    targetPort: 8080
  - name: backend
    port: 3000
    targetPort: 3000
  type: LoadBalancer
```

#### 部署到Kubernetes
```bash
kubectl apply -f kubernetes.yaml
```

## 6. 验证部署

### 6.1 验证前端

访问前端URL (开发环境为 http://localhost:8080，生产环境为您配置的域名)，确认能够正常加载页面。

### 6.2 验证后端

访问后端API (例如 http://localhost:3000/api/providers)，确认能够获取到云服务提供商列表。

### 6.3 验证Terraform

在后端服务器上执行以下命令，确认Terraform能够正常工作:
```bash
terraform version
```

## 7. 故障排除

### 7.1 前端无法访问

- 检查前端服务是否正在运行
- 检查端口是否被占用
- 检查浏览器控制台是否有错误信息

### 7.2 后端API错误

- 检查后端服务是否正在运行
- 检查环境变量配置是否正确
- 检查服务器日志是否有错误信息

### 7.3 Terraform执行错误

- 检查Terraform是否正确安装
- 检查云服务提供商凭证是否正确配置
- 检查Terraform执行日志是否有错误信息

## 8. 备份和恢复

### 8.1 备份

#### 备份配置文件
```bash
cp -r backend/terraform /backup/terraform
cp -r frontend/frontend-app/.env /backup/frontend.env
cp -r backend/.env /backup/backend.env
```

#### 备份部署状态
如果使用了Terraform远程状态存储，确保定期备份状态文件。

### 8.2 恢复

#### 恢复配置文件
```bash
cp -r /backup/terraform backend/
cp /backup/frontend.env frontend/frontend-app/.env
cp /backup/backend.env backend/.env
```

#### 恢复部署状态
将备份的Terraform状态文件恢复到相应位置。

## 9. 监控和日志

### 9.1 应用日志

- 前端日志: 浏览器控制台
- 后端日志: 标准输出或配置的日志文件

### 9.2 系统监控

建议使用以下工具监控应用:
- Prometheus + Grafana: 监控系统资源和应用性能
- ELK Stack: 集中管理日志
- Uptime Robot: 监控应用可用性

## 10. 安全建议

### 10.1 网络安全

- 使用HTTPS保护前端通信
- 配置防火墙，只开放必要端口
- 使用API网关保护后端API

### 10.2 认证和授权

- 实现用户认证系统
- 使用JWT或OAuth进行API认证
- 实现基于角色的访问控制

### 10.3 云凭证安全

- 使用环境变量或密钥管理系统存储云凭证
- 定期轮换访问密钥
- 使用最小权限原则配置云服务访问权限

## 11. 升级指南

### 11.1 前端升级

```bash
cd frontend/frontend-app
git pull
npm install
npm run build
```

### 11.2 后端升级

```bash
cd backend
git pull
npm install
pm2 restart server.js # 如果使用pm2
```

### 11.3 Terraform升级

```bash
sudo apt-get update && sudo apt-get install terraform # 对于Debian/Ubuntu
brew upgrade terraform # 对于macOS
```

## 12. 性能优化

### 12.1 前端优化

- 启用Gzip压缩
- 配置浏览器缓存
- 使用CDN分发静态资源

### 12.2 后端优化

- 增加服务器资源
- 实现API缓存
- 优化数据库查询（如果使用数据库）

### 12.3 Terraform优化

- 使用远程状态存储
- 启用并行执行
- 使用模块化结构提高复用性
