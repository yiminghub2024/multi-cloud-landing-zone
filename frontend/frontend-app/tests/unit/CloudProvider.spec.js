import { shallowMount } from '@vue/test-utils'
import CloudProvider from '@/components/CloudProvider.vue'

describe('CloudProvider.vue', () => {
  it('renders cloud provider options correctly', () => {
    const wrapper = shallowMount(CloudProvider)
    
    // 验证组件包含7个云服务提供商选项
    expect(wrapper.findAll('.cloud-option').length).toBe(7)
    
    // 验证包含AWS选项
    expect(wrapper.text()).toContain('AWS')
    
    // 验证包含Azure选项
    expect(wrapper.text()).toContain('Azure')
    
    // 验证包含阿里云选项
    expect(wrapper.text()).toContain('阿里云')
    
    // 验证包含百度云选项
    expect(wrapper.text()).toContain('百度云')
    
    // 验证包含华为云选项
    expect(wrapper.text()).toContain('华为云')
    
    // 验证包含腾讯云选项
    expect(wrapper.text()).toContain('腾讯云')
    
    // 验证包含火山云选项
    expect(wrapper.text()).toContain('火山云')
  })
  
  it('selects cloud provider when clicked', async () => {
    const wrapper = shallowMount(CloudProvider)
    
    // 初始状态下没有选中的云服务提供商
    expect(wrapper.vm.selectedCloud).toBe('')
    
    // 点击AWS选项
    await wrapper.findAll('.cloud-option').at(0).trigger('click')
    
    // 验证选中的云服务提供商是AWS
    expect(wrapper.vm.selectedCloud).toBe('aws')
    
    // 验证AWS选项有selected类
    expect(wrapper.findAll('.cloud-option').at(0).classes()).toContain('selected')
  })
  
  it('emits cloud-selected event when cloud provider is selected', async () => {
    const wrapper = shallowMount(CloudProvider)
    
    // 点击Azure选项
    await wrapper.findAll('.cloud-option').at(1).trigger('click')
    
    // 验证触发了cloud-selected事件
    expect(wrapper.emitted('cloud-selected')).toBeTruthy()
    
    // 验证事件参数是azure
    expect(wrapper.emitted('cloud-selected')[0]).toEqual(['azure'])
  })
  
  it('emits regions-loaded event after cloud provider is selected', async () => {
    const wrapper = shallowMount(CloudProvider)
    
    // 模拟fetchCloudRegions方法
    wrapper.vm.fetchCloudRegions = jest.fn()
    
    // 点击阿里云选项
    await wrapper.findAll('.cloud-option').at(2).trigger('click')
    
    // 验证调用了fetchCloudRegions方法
    expect(wrapper.vm.fetchCloudRegions).toHaveBeenCalledWith('alicloud')
  })
})
