#!/bin/bash

# 测试CORS配置的脚本
# 此脚本模拟前端请求，测试后端CORS配置是否正确

echo "===== 多云着陆区部署平台 CORS 测试工具 ====="
echo "测试后端API的CORS配置是否正确"
echo ""

# 默认设置
API_HOST="10.168.0.5"
API_PORT="3000"
ORIGIN="http://10.168.0.5:8080"
ENDPOINT="/api/deploy"
METHOD="POST"
CONTENT_TYPE="application/json"
TEST_DATA='{"cloudProvider":"aws","region":"us-east-1","az":"us-east-1a","vpc":"vpc-default","subnet":"subnet-default","components":["ec2","s3"]}'

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印帮助信息
print_help() {
  echo "用法: $0 [选项]"
  echo ""
  echo "选项:"
  echo "  -h, --host HOST      设置API主机 (默认: $API_HOST)"
  echo "  -p, --port PORT      设置API端口 (默认: $API_PORT)"
  echo "  -o, --origin ORIGIN  设置Origin头部 (默认: $ORIGIN)"
  echo "  -e, --endpoint PATH  设置API端点 (默认: $ENDPOINT)"
  echo "  -m, --method METHOD  设置HTTP方法 (默认: $METHOD)"
  echo "  -d, --data DATA      设置请求数据 (默认: 简单的部署配置)"
  echo "  --help               显示此帮助信息"
  echo ""
  echo "示例:"
  echo "  $0 --host 10.168.0.5 --port 3000 --origin http://10.168.0.5:8080"
  echo ""
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
  case $1 in
    -h|--host)
      API_HOST="$2"
      shift 2
      ;;
    -p|--port)
      API_PORT="$2"
      shift 2
      ;;
    -o|--origin)
      ORIGIN="$2"
      shift 2
      ;;
    -e|--endpoint)
      ENDPOINT="$2"
      shift 2
      ;;
    -m|--method)
      METHOD="$2"
      shift 2
      ;;
    -d|--data)
      TEST_DATA="$2"
      shift 2
      ;;
    --help)
      print_help
      exit 0
      ;;
    *)
      echo -e "${RED}错误: 未知选项 $1${NC}"
      print_help
      exit 1
      ;;
  esac
done

API_URL="http://${API_HOST}:${API_PORT}${ENDPOINT}"

echo -e "${BLUE}测试配置:${NC}"
echo "API URL: $API_URL"
echo "Origin: $ORIGIN"
echo "Method: $METHOD"
echo "Content-Type: $CONTENT_TYPE"
echo "测试数据: $TEST_DATA"
echo ""

# 测试1: 发送预检请求
echo -e "${YELLOW}测试1: 发送预检请求 (OPTIONS)${NC}"
echo "发送OPTIONS请求到 $API_URL..."

OPTIONS_RESPONSE=$(curl -s -i -X OPTIONS \
  -H "Origin: $ORIGIN" \
  -H "Access-Control-Request-Method: $METHOD" \
  -H "Access-Control-Request-Headers: Content-Type" \
  "$API_URL")

echo -e "${BLUE}预检请求响应:${NC}"
echo "$OPTIONS_RESPONSE"
echo ""

# 检查预检响应中的CORS头部
if echo "$OPTIONS_RESPONSE" | grep -q "Access-Control-Allow-Origin"; then
  echo -e "${GREEN}✓ 预检响应包含 Access-Control-Allow-Origin 头部${NC}"
else
  echo -e "${RED}✗ 预检响应缺少 Access-Control-Allow-Origin 头部${NC}"
fi

if echo "$OPTIONS_RESPONSE" | grep -q "Access-Control-Allow-Credentials"; then
  echo -e "${GREEN}✓ 预检响应包含 Access-Control-Allow-Credentials 头部${NC}"
else
  echo -e "${RED}✗ 预检响应缺少 Access-Control-Allow-Credentials 头部${NC}"
fi

if echo "$OPTIONS_RESPONSE" | grep -q "Access-Control-Allow-Methods"; then
  echo -e "${GREEN}✓ 预检响应包含 Access-Control-Allow-Methods 头部${NC}"
else
  echo -e "${RED}✗ 预检响应缺少 Access-Control-Allow-Methods 头部${NC}"
fi

if echo "$OPTIONS_RESPONSE" | grep -q "Access-Control-Allow-Headers"; then
  echo -e "${GREEN}✓ 预检响应包含 Access-Control-Allow-Headers 头部${NC}"
else
  echo -e "${RED}✗ 预检响应缺少 Access-Control-Allow-Headers 头部${NC}"
fi

echo ""

# 测试2: 发送实际请求
echo -e "${YELLOW}测试2: 发送实际请求 ($METHOD)${NC}"
echo "发送$METHOD请求到 $API_URL..."

ACTUAL_RESPONSE=$(curl -s -i -X "$METHOD" \
  -H "Origin: $ORIGIN" \
  -H "Content-Type: $CONTENT_TYPE" \
  -d "$TEST_DATA" \
  "$API_URL")

echo -e "${BLUE}实际请求响应:${NC}"
echo "$ACTUAL_RESPONSE"
echo ""

# 检查实际响应中的CORS头部
if echo "$ACTUAL_RESPONSE" | grep -q "Access-Control-Allow-Origin"; then
  echo -e "${GREEN}✓ 实际响应包含 Access-Control-Allow-Origin 头部${NC}"
else
  echo -e "${RED}✗ 实际响应缺少 Access-Control-Allow-Origin 头部${NC}"
fi

if echo "$ACTUAL_RESPONSE" | grep -q "Access-Control-Allow-Credentials"; then
  echo -e "${GREEN}✓ 实际响应包含 Access-Control-Allow-Credentials 头部${NC}"
else
  echo -e "${RED}✗ 实际响应缺少 Access-Control-Allow-Credentials 头部${NC}"
fi

echo ""

# 测试3: 检查API响应状态
if echo "$ACTUAL_RESPONSE" | grep -q "HTTP/1.1 200"; then
  echo -e "${GREEN}✓ API响应状态码为 200 OK${NC}"
else
  echo -e "${RED}✗ API响应状态码不是 200 OK${NC}"
fi

# 检查API响应内容
if echo "$ACTUAL_RESPONSE" | grep -q "\"success\":true"; then
  echo -e "${GREEN}✓ API响应包含 success:true${NC}"
else
  echo -e "${RED}✗ API响应不包含 success:true${NC}"
fi

echo ""
echo -e "${BLUE}CORS测试完成${NC}"
echo "如果所有检查都通过，则CORS配置正确。"
echo "如果有任何失败的检查，请参考错误信息修复CORS配置。"
