import { shallowMount } from '@vue/test-utils'
import VPCConfig from '@/components/VPCConfig.vue'

describe('VPCConfig.vue', () => {
  it('renders VPC configuration form correctly', () => {
    const wrapper = shallowMount(VPCConfig, {
      props: {
        cloudProvider: 'aws',
        region: 'us-east-1'
      }
    })
    
    // 验证表单元素存在
    expect(wrapper.find('#vpc-name').exists()).toBe(true)
    expect(wrapper.find('#vpc-cidr').exists()).toBe(true)
    expect(wrapper.find('.vpc-options').exists()).toBe(true)
    
    // 验证DNS选项存在
    expect(wrapper.findAll('input[type="checkbox"]').length).toBe(2)
    expect(wrapper.text()).toContain('启用DNS支持')
    expect(wrapper.text()).toContain('启用DNS主机名')
  })
  
  it('sets default VPC name based on cloud provider and region', async () => {
    const wrapper = shallowMount(VPCConfig, {
      props: {
        cloudProvider: 'aws',
        region: 'us-east-1'
      }
    })
    
    // 验证默认VPC名称
    expect(wrapper.vm.vpcName).toBe('AWS-VPC-us-east-1')
    
    // 更改云服务提供商
    await wrapper.setProps({ cloudProvider: 'azure' })
    
    // 验证VPC名称更新
    expect(wrapper.vm.vpcName).toBe('Azure-VNet-us-east-1')
  })
  
  it('validates VPC name correctly', async () => {
    const wrapper = shallowMount(VPCConfig, {
      props: {
        cloudProvider: 'aws',
        region: 'us-east-1'
      }
    })
    
    // 设置空VPC名称
    await wrapper.setData({ vpcName: '' })
    wrapper.vm.validateForm()
    
    // 验证错误消息
    expect(wrapper.vm.errors.vpcName).toBe('VPC名称不能为空')
    
    // 设置过短的VPC名称
    await wrapper.setData({ vpcName: 'a' })
    wrapper.vm.validateForm()
    
    // 验证错误消息
    expect(wrapper.vm.errors.vpcName).toBe('VPC名称至少需要2个字符')
    
    // 设置有效的VPC名称
    await wrapper.setData({ vpcName: 'test-vpc' })
    wrapper.vm.validateForm()
    
    // 验证没有错误消息
    expect(wrapper.vm.errors.vpcName).toBe('')
  })
  
  it('validates CIDR block correctly', async () => {
    const wrapper = shallowMount(VPCConfig, {
      props: {
        cloudProvider: 'aws',
        region: 'us-east-1'
      }
    })
    
    // 设置空CIDR
    await wrapper.setData({ vpcCIDR: '' })
    wrapper.vm.validateForm()
    
    // 验证错误消息
    expect(wrapper.vm.errors.vpcCIDR).toBe('CIDR块不能为空')
    
    // 设置无效的CIDR格式
    await wrapper.setData({ vpcCIDR: '10.0.0.0' })
    wrapper.vm.validateForm()
    
    // 验证错误消息
    expect(wrapper.vm.errors.vpcCIDR).toBe('请输入有效的CIDR格式，例如: 10.0.0.0/16')
    
    // 设置无效的IP地址
    await wrapper.setData({ vpcCIDR: '256.0.0.0/16' })
    wrapper.vm.validateForm()
    
    // 验证错误消息
    expect(wrapper.vm.errors.vpcCIDR).toBe('IP地址部分无效，每个部分应在0-255之间')
    
    // 设置有效的CIDR
    await wrapper.setData({ vpcCIDR: '10.0.0.0/16' })
    wrapper.vm.validateForm()
    
    // 验证没有错误消息
    expect(wrapper.vm.errors.vpcCIDR).toBe('')
  })
  
  it('emits vpc-updated event when form is valid', async () => {
    const wrapper = shallowMount(VPCConfig, {
      props: {
        cloudProvider: 'aws',
        region: 'us-east-1'
      }
    })
    
    // 设置有效的表单数据
    await wrapper.setData({
      vpcName: 'test-vpc',
      vpcCIDR: '10.0.0.0/16',
      enableDnsSupport: true,
      enableDnsHostnames: true
    })
    
    // 验证表单
    wrapper.vm.validateForm()
    
    // 验证触发了vpc-updated事件
    expect(wrapper.emitted('vpc-updated')).toBeTruthy()
    
    // 验证事件参数
    expect(wrapper.emitted('vpc-updated')[0][0]).toEqual({
      name: 'test-vpc',
      cidr: '10.0.0.0/16',
      enableDnsSupport: true,
      enableDnsHostnames: true
    })
  })
})
