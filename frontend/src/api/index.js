// API配置文件
// 定义后端API的基本URL和端点

// 后端API基本URL
const API_BASE_URL = 'http://10.168.0.5:3000';

// API端点
const API_ENDPOINTS = {
  deploy: '/api/deploy'
};

// 完整API URL
const API_URLS = {
  deploy: `${API_BASE_URL}${API_ENDPOINTS.deploy}`
};

// 使用fetch API发送请求的通用函数
async function fetchApi(url, method = 'GET', data = null) {
  const options = {
    method,
    headers: {
      'Content-Type': 'application/json'
    },
    // 简化CORS配置
    mode: 'cors',
    credentials: 'include'
  };

  if (data) {
    options.body = JSON.stringify(data);
  }

  try {
    // 添加详细的调试日志
    console.log('准备发送API请求:', {
      url,
      method,
      headers: options.headers,
      mode: options.mode,
      credentials: options.credentials,
      hasBody: !!options.body
    });
    
    // 对于非简单请求，先发送预检请求
    if (method !== 'GET' && method !== 'HEAD') {
      console.log('发送跨域请求，可能需要预检请求...');
    }
    
    const response = await fetch(url, options);
    
    // 记录响应信息
    console.log('收到API响应:', {
      url,
      status: response.status,
      statusText: response.statusText,
      headers: Object.fromEntries([...response.headers.entries()]),
      ok: response.ok
    });
    
    // 检查响应状态
    if (!response.ok) {
      throw new Error(`API请求失败: ${response.status} ${response.statusText}`);
    }
    
    // 解析JSON响应
    const result = await response.json();
    console.log('API响应数据:', result);
    return result;
  } catch (error) {
    console.error('API请求错误:', error);
    
    // 特别处理CORS错误
    if (error.message.includes('CORS') || error.message.includes('跨源')) {
      console.error('CORS错误详情:', {
        message: error.message,
        url,
        method,
        mode: options.mode,
        credentials: options.credentials
      });
      console.error('CORS错误排查提示:');
      console.error('1. 检查后端CORS配置是否正确');
      console.error('2. 确认前端源与后端允许的Origin匹配');
      console.error('3. 使用credentials: "include"时，后端不能使用通配符"*"作为Access-Control-Allow-Origin');
      console.error('4. 检查网络连接是否正常');
      
      throw new Error('跨域请求被拒绝，请查看控制台获取详细错误信息和排查提示。详细错误: ' + error.message);
    }
    throw error;
  }
}

// 部署API函数
export async function deployInfrastructure(deploymentData) {
  return fetchApi(API_URLS.deploy, 'POST', deploymentData);
}

// 导出API配置
export default {
  baseUrl: API_BASE_URL,
  endpoints: API_ENDPOINTS,
  urls: API_URLS,
  deployInfrastructure
};
