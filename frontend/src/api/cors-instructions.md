# CORS 配置说明

## 前端CORS配置

前端代码已经添加了必要的CORS配置，包括：

1. 设置了 `mode: 'cors'` 明确指定跨域请求模式
2. 添加了 `credentials: 'include'` 允许跨域请求携带凭证
3. 添加了CORS相关请求头
4. 增强了错误处理，特别是针对CORS错误的处理

## 后端CORS配置

要完全解决CORS问题，后端也需要正确配置。以下是常见后端框架的CORS配置示例：

### Express.js (Node.js)

```javascript
const express = require('express');
const cors = require('cors');
const app = express();

// 配置CORS
app.use(cors({
  origin: 'http://前端域名', // 或使用 '*' 允许所有来源
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization'],
  credentials: true // 允许携带凭证
}));

// API路由
app.post('/api/deploy', (req, res) => {
  // 处理部署请求
});

app.listen(3000, () => {
  console.log('服务器运行在端口3000');
});
```

### Spring Boot (Java)

```java
@Configuration
public class CorsConfig implements WebMvcConfigurer {
    @Override
    public void addCorsMappings(CorsRegistry registry) {
        registry.addMapping("/**")
            .allowedOrigins("http://前端域名") // 或使用 "*" 允许所有来源
            .allowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
            .allowedHeaders("Content-Type", "Authorization")
            .allowCredentials(true);
    }
}
```

### Django (Python)

在settings.py中添加：

```python
INSTALLED_APPS = [
    # ...
    'corsheaders',
    # ...
]

MIDDLEWARE = [
    # ...
    'corsheaders.middleware.CorsMiddleware',
    'django.middleware.common.CommonMiddleware',
    # ...
]

CORS_ALLOW_ALL_ORIGINS = True  # 或者指定允许的域名
CORS_ALLOW_CREDENTIALS = True
CORS_ALLOW_METHODS = ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS']
CORS_ALLOW_HEADERS = ['Content-Type', 'Authorization']
```

## 测试CORS配置

1. 确保后端服务器正确配置了CORS头部
2. 启动前端应用，尝试点击"开始部署"按钮
3. 在浏览器开发者工具的网络面板中检查请求，确认没有CORS错误

如果仍然遇到CORS问题，请检查：
1. 后端服务器是否正确响应OPTIONS预检请求
2. 后端返回的Access-Control-Allow-Origin头部是否匹配前端域名
3. 如果使用了credentials: 'include'，后端必须设置具体的Origin而不能使用通配符'*'
