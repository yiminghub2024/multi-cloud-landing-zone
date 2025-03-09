import { shallowMount } from '@vue/test-utils'
import DeploymentFeedback from '@/components/DeploymentFeedback.vue'

describe('DeploymentFeedback.vue', () => {
  it('displays idle state when status is idle', () => {
    const wrapper = shallowMount(DeploymentFeedback, {
      props: {
        status: 'idle',
        progress: 0,
        message: '',
        logs: [],
        result: {},
        topology: { nodes: [], edges: [] }
      }
    })
    
    // 验证显示空闲状态
    expect(wrapper.find('.idle-state').exists()).toBe(true)
    expect(wrapper.text()).toContain('尚未开始部署')
  })
  
  it('displays deploying state with progress bar when status is preparing or deploying', () => {
    const wrapper = shallowMount(DeploymentFeedback, {
      props: {
        status: 'deploying',
        progress: 45,
        message: '正在执行Terraform部署...',
        logs: ['开始部署过程...', 'Terraform初始化完成'],
        result: {},
        topology: { nodes: [], edges: [] }
      }
    })
    
    // 验证显示部署中状态
    expect(wrapper.find('.deploying-state').exists()).toBe(true)
    expect(wrapper.find('.progress-bar').exists()).toBe(true)
    expect(wrapper.find('.progress').attributes('style')).toContain('width: 45%')
    expect(wrapper.find('.progress-text').text()).toBe('45%')
    expect(wrapper.find('.status-message').text()).toBe('正在执行Terraform部署...')
    
    // 验证显示日志
    expect(wrapper.find('.deployment-logs').exists()).toBe(true)
    expect(wrapper.findAll('.log-entry').length).toBe(2)
    expect(wrapper.findAll('.log-entry').at(0).text()).toBe('开始部署过程...')
    expect(wrapper.findAll('.log-entry').at(1).text()).toBe('Terraform初始化完成')
  })
  
  it('displays completed state with success message when status is completed', () => {
    const wrapper = shallowMount(DeploymentFeedback, {
      props: {
        status: 'completed',
        progress: 100,
        message: '部署完成',
        logs: ['开始部署过程...', 'Terraform部署执行完成'],
        result: {
          deploymentId: '1234567890',
          cloudProvider: 'aws',
          region: 'us-east-1',
          az: 'us-east-1a'
        },
        topology: { 
          nodes: [
            { id: 'vpc', type: 'vpc', name: 'test-vpc' }
          ], 
          edges: [] 
        }
      },
      methods: {
        renderTopology: jest.fn() // 模拟renderTopology方法
      }
    })
    
    // 验证显示完成状态
    expect(wrapper.find('.completed-state').exists()).toBe(true)
    expect(wrapper.find('.success-message').exists()).toBe(true)
    expect(wrapper.find('.success-text').text()).toBe('部署成功完成！')
    
    // 验证显示部署详情
    expect(wrapper.find('.deployment-details').exists()).toBe(true)
    expect(wrapper.text()).toContain('部署ID:')
    expect(wrapper.text()).toContain('1234567890')
    expect(wrapper.text()).toContain('云服务提供商:')
    expect(wrapper.text()).toContain('区域:')
    expect(wrapper.text()).toContain('us-east-1')
    
    // 验证显示拓扑图容器
    expect(wrapper.find('.topology-container').exists()).toBe(true)
    expect(wrapper.find('.topology-graph').exists()).toBe(true)
  })
  
  it('displays failed state with error message when status is failed', () => {
    const wrapper = shallowMount(DeploymentFeedback, {
      props: {
        status: 'failed',
        progress: 60,
        message: '部署失败: Terraform配置验证失败',
        logs: ['开始部署过程...', '错误: Terraform配置验证失败'],
        result: {},
        topology: { nodes: [], edges: [] }
      }
    })
    
    // 验证显示失败状态
    expect(wrapper.find('.failed-state').exists()).toBe(true)
    expect(wrapper.find('.error-message').exists()).toBe(true)
    expect(wrapper.find('.error-text').text()).toBe('部署失败')
    
    // 验证显示错误详情
    expect(wrapper.find('.error-details').exists()).toBe(true)
    expect(wrapper.find('.error-container').text()).toBe('部署失败: Terraform配置验证失败')
    
    // 验证显示重试按钮
    expect(wrapper.find('.retry-actions').exists()).toBe(true)
    expect(wrapper.find('.retry-button').exists()).toBe(true)
    expect(wrapper.find('.back-button').exists()).toBe(true)
  })
  
  it('emits refresh-status event when startRefreshing is called', async () => {
    jest.useFakeTimers()
    
    const wrapper = shallowMount(DeploymentFeedback, {
      props: {
        status: 'deploying',
        progress: 30,
        message: '正在初始化Terraform...',
        logs: ['开始部署过程...'],
        result: {},
        topology: { nodes: [], edges: [] }
      }
    })
    
    // 调用startRefreshing方法
    wrapper.vm.startRefreshing()
    
    // 快进3秒
    jest.advanceTimersByTime(3000)
    
    // 验证触发了refresh-status事件
    expect(wrapper.emitted('refresh-status')).toBeTruthy()
    
    // 清理
    wrapper.vm.stopRefreshing()
    jest.useRealTimers()
  })
  
  it('emits retry-deployment event when retry button is clicked', async () => {
    const wrapper = shallowMount(DeploymentFeedback, {
      props: {
        status: 'failed',
        progress: 60,
        message: '部署失败: Terraform配置验证失败',
        logs: ['开始部署过程...', '错误: Terraform配置验证失败'],
        result: {},
        topology: { nodes: [], edges: [] }
      }
    })
    
    // 点击重试按钮
    await wrapper.find('.retry-button').trigger('click')
    
    // 验证触发了retry-deployment事件
    expect(wrapper.emitted('retry-deployment')).toBeTruthy()
  })
  
  it('emits back-to-summary event when back button is clicked', async () => {
    const wrapper = shallowMount(DeploymentFeedback, {
      props: {
        status: 'failed',
        progress: 60,
        message: '部署失败: Terraform配置验证失败',
        logs: ['开始部署过程...', '错误: Terraform配置验证失败'],
        result: {},
        topology: { nodes: [], edges: [] }
      }
    })
    
    // 点击返回按钮
    await wrapper.find('.back-button').trigger('click')
    
    // 验证触发了back-to-summary事件
    expect(wrapper.emitted('back-to-summary')).toBeTruthy()
  })
})
