<template>
  <div class="deployment-feedback-container">
    <el-card shadow="hover">
      <template #header>
        <h2>部署状态</h2>
      </template>
      <div class="deployment-content">
        <el-empty v-if="status === 'idle'" description="尚未开始部署，请在摘要页面点击'开始部署'按钮" />
        
        <div v-else-if="status === 'preparing' || status === 'deploying'" class="deploying-state">
          <el-progress 
            :percentage="progress" 
            :format="format => `${format}%`" 
            :stroke-width="20"
            status="primary"
          />
          
          <el-alert
            :title="message"
            type="info"
            :closable="false"
            class="mt-4"
          />
          
          <div class="deployment-logs">
            <el-divider content-position="left">部署日志</el-divider>
            <el-scrollbar height="300px">
              <div class="logs-container">
                <div v-for="(log, index) in logs" :key="index" class="log-entry">
                  {{ log }}
                </div>
              </div>
            </el-scrollbar>
          </div>
        </div>
        
        <div v-else-if="status === 'completed'" class="completed-state">
          <el-result
            icon="success"
            title="部署成功完成！"
            sub-title="您的云资源已成功部署"
          />
          
          <el-descriptions title="部署详情" :column="1" border>
            <el-descriptions-item label="部署ID">{{ result.deploymentId }}</el-descriptions-item>
            <el-descriptions-item label="云服务提供商">{{ getCloudProviderName(result.cloudProvider) }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ result.region }}</el-descriptions-item>
            <el-descriptions-item label="可用区">{{ result.az }}</el-descriptions-item>
          </el-descriptions>
          
          <div class="topology-container">
            <el-divider content-position="left">资源拓扑图</el-divider>
            <div class="topology-visualization">
              <div ref="topologyGraph" class="topology-graph"></div>
            </div>
          </div>
          
          <div class="deployment-logs">
            <el-divider content-position="left">部署日志</el-divider>
            <el-scrollbar height="300px">
              <div class="logs-container">
                <div v-for="(log, index) in logs" :key="index" class="log-entry">
                  {{ log }}
                </div>
              </div>
            </el-scrollbar>
          </div>
        </div>
        
        <div v-else-if="status === 'failed'" class="failed-state">
          <el-result
            icon="error"
            title="部署失败"
            sub-title="请查看错误详情并尝试重新部署"
          >
            <template #extra>
              <div class="retry-actions">
                <el-button type="primary" @click="retryDeployment">重试部署</el-button>
                <el-button @click="backToSummary">返回摘要页面</el-button>
              </div>
            </template>
          </el-result>
          
          <el-alert
            :title="message"
            type="error"
            :closable="false"
            class="mt-4"
            :description="message"
            show-icon
          />
          
          <div class="deployment-logs">
            <el-divider content-position="left">部署日志</el-divider>
            <el-scrollbar height="300px">
              <div class="logs-container">
                <div v-for="(log, index) in logs" :key="index" class="log-entry">
                  {{ log }}
                </div>
              </div>
            </el-scrollbar>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import * as d3 from 'd3';

export default {
  name: 'DeploymentFeedback',
  props: {
    status: {
      type: String,
      default: 'idle' // idle, preparing, deploying, completed, failed
    },
    progress: {
      type: Number,
      default: 0
    },
    message: {
      type: String,
      default: ''
    },
    logs: {
      type: Array,
      default: () => []
    },
    result: {
      type: Object,
      default: () => ({})
    },
    topology: {
      type: Object,
      default: () => ({ nodes: [], edges: [] })
    }
  },
  data() {
    return {
      refreshInterval: null,
      simulation: null,
      svg: null
    }
  },
  mounted() {
    // 如果状态是部署中或准备中，则启动自动刷新
    if (this.status === 'preparing' || this.status === 'deploying') {
      this.startRefreshing();
    }
    
    // 如果状态是已完成，则渲染拓扑图
    if (this.status === 'completed' && this.topology && this.topology.nodes) {
      this.$nextTick(() => {
        this.renderTopology();
      });
    }
  },
  beforeDestroy() {
    this.stopRefreshing();
  },
  watch: {
    status(newStatus) {
      if (newStatus === 'preparing' || newStatus === 'deploying') {
        this.startRefreshing();
      } else {
        this.stopRefreshing();
        
        if (newStatus === 'completed' && this.topology && this.topology.nodes) {
          this.$nextTick(() => {
            this.renderTopology();
          });
        }
      }
    },
    topology: {
      deep: true,
      handler(newTopology) {
        if (this.status === 'completed' && newTopology && newTopology.nodes) {
          this.$nextTick(() => {
            this.renderTopology();
          });
        }
      }
    }
  },
  methods: {
    getCloudProviderName(provider) {
      const providerMap = {
        'aws': 'AWS',
        'azure': 'Azure',
        'alicloud': '阿里云',
        'baidu': '百度云',
        'huawei': '华为云',
        'tencent': '腾讯云',
        'volcengine': '火山云'
      };
      
      return providerMap[provider] || provider;
    },
    startRefreshing() {
      // 每3秒刷新一次部署状态
      this.refreshInterval = setInterval(() => {
        this.$emit('refresh-status');
      }, 3000);
    },
    stopRefreshing() {
      if (this.refreshInterval) {
        clearInterval(this.refreshInterval);
        this.refreshInterval = null;
      }
    },
    retryDeployment() {
      this.$emit('retry-deployment');
    },
    backToSummary() {
      this.$emit('back-to-summary');
    },
    renderTopology() {
      if (!this.topology || !this.topology.nodes || !this.topology.edges || !this.$refs.topologyGraph) {
        return;
      }
      
      // 清除之前的图形
      d3.select(this.$refs.topologyGraph).selectAll('*').remove();
      
      const width = this.$refs.topologyGraph.clientWidth;
      const height = 500;
      
      // 创建SVG
      this.svg = d3.select(this.$refs.topologyGraph)
        .append('svg')
        .attr('width', width)
        .attr('height', height);
      
      // 定义节点图标
      const nodeIcons = {
        'cloud': '\uf0c2', // cloud
        'vpc': '\uf233', // server
        'subnet': '\uf6ff', // network-wired
        'load-balancer': '\uf013', // cog
        'object-storage': '\uf1c0', // database
        'rds': '\uf1c0', // database
        'compute': '\uf108', // desktop
        'cdn': '\uf0e8', // sitemap
        'default': '\uf111' // circle
      };
      
      // 定义节点颜色
      const nodeColors = {
        'cloud': '#3498db',
        'vpc': '#2ecc71',
        'subnet': '#1abc9c',
        'load-balancer': '#e74c3c',
        'object-storage': '#9b59b6',
        'rds': '#8e44ad',
        'compute': '#f39c12',
        'cdn': '#d35400',
        'default': '#7f8c8d'
      };
      
      // 创建箭头标记
      this.svg.append('defs').append('marker')
        .attr('id', 'arrowhead')
        .attr('viewBox', '-0 -5 10 10')
        .attr('refX', 20)
        .attr('refY', 0)
        .attr('orient', 'auto')
        .attr('markerWidth', 6)
        .attr('markerHeight', 6)
        .attr('xoverflow', 'visible')
        .append('svg:path')
        .attr('d', 'M 0,-5 L 10 ,0 L 0,5')
        .attr('fill', '#999')
        .style('stroke', 'none');
      
      // 创建力导向图
      this.simulation = d3.forceSimulation()
        .force('link', d3.forceLink().id(d => d.id).distance(100))
        .force('charge', d3.forceManyBody().strength(-300))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .force('collision', d3.forceCollide().radius(50));
      
      // 创建连线
      const link = this.svg.append('g')
        .selectAll('line')
        .data(this.topology.edges)
        .enter().append('line')
        .attr('stroke', '#999')
        .attr('stroke-opacity', 0.6)
        .attr('stroke-width', 2)
        .attr('marker-end', 'url(#arrowhead)');
      
      // 创建节点组
      const node = this.svg.append('g')
        .selectAll('.node')
        .data(this.topology.nodes)
        .enter().append('g')
        .attr('class', 'node')
        .call(d3.drag()
          .on('start', dragstarted)
          .on('drag', dragged)
          .on('end', dragended));
      
      // 添加节点圆形背景
      node.append('circle')
        .attr('r', 20)
        .attr('fill', d => nodeColors[d.type] || nodeColors.default);
      
      // 添加节点图标
      node.append('text')
        .attr('text-anchor', 'middle')
        .attr('dy', '0.35em')
        .attr('font-family', 'FontAwesome')
        .attr('font-size', '16px')
        .attr('fill', 'white')
        .text(d => nodeIcons[d.type] || nodeIcons.default);
      
      // 添加节点标签
      node.append('text')
        .attr('dy', 30)
        .attr('text-anchor', 'middle')
        .text(d => d.name)
        .attr('font-size', '12px')
        .attr('fill', '#333');
      
      // 更新力导向图
      this.simulation
        .nodes(this.topology.nodes)
        .on('tick', ticked);
      
      this.simulation.force('link')
        .links(this.topology.edges);
      
      // 定义拖拽行为
      const that = this;
      function dragstarted(event, d) {
        if (!event.active) that.simulation.alphaTarget(0.3).restart();
        d.fx = d.x;
        d.fy = d.y;
      }
      
      function dragged(event, d) {
        d.fx = event.x;
        d.fy = event.y;
      }
      
      function dragended(event, d) {
        if (!event.active) that.simulation.alphaTarget(0);
        d.fx = null;
        d.fy = null;
      }
      
      // 定义tick函数
      function ticked() {
        link
          .attr('x1', d => d.source.x)
          .attr('y1', d => d.source.y)
          .attr('x2', d => d.target.x)
          .attr('y2', d => d.target.y);
        
        node
          .attr('transform', d => `translate(${d.x},${d.y})`);
      }
    }
  }
}
</script>

<style scoped>
.deployment-feedback-container {
  margin-bottom: 2rem;
}

.deployment-content {
  margin-top: 1rem;
}

.mt-4 {
  margin-top: 1rem;
}

.deployment-logs {
  margin-top: 2rem;
}

.logs-container {
  padding: 1rem;
  background-color: #2c3e50;
  border-radius: 4px;
  color: #ecf0f1;
  font-family: monospace;
}

.log-entry {
  margin-bottom: 0.5rem;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
}

.topology-container {
  margin: 2rem 0;
}

.topology-visualization {
  padding: 1rem;
  background-color: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
}

.topology-graph {
  width: 100%;
  height: 500px;
  overflow: hidden;
}

.retry-actions {
  display: flex;
  justify-content: center;
  gap: 1rem;
}
</style>
