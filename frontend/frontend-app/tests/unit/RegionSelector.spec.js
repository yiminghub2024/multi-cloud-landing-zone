import { shallowMount } from '@vue/test-utils'
import RegionSelector from '@/components/RegionSelector.vue'

describe('RegionSelector.vue', () => {
  it('displays loading state when loading is true', () => {
    const wrapper = shallowMount(RegionSelector, {
      props: {
        cloudProvider: 'aws',
        initialRegions: []
      },
      data() {
        return {
          loading: true
        }
      }
    })
    
    // 验证显示加载指示器
    expect(wrapper.find('.loading-indicator').exists()).toBe(true)
    expect(wrapper.find('.spinner').exists()).toBe(true)
    expect(wrapper.text()).toContain('正在加载区域数据')
  })
  
  it('displays no data message when regions array is empty', () => {
    const wrapper = shallowMount(RegionSelector, {
      props: {
        cloudProvider: 'aws',
        initialRegions: []
      },
      data() {
        return {
          loading: false,
          regions: []
        }
      }
    })
    
    // 验证显示无数据消息
    expect(wrapper.find('.no-data').exists()).toBe(true)
    expect(wrapper.text()).toContain('没有找到区域数据')
  })
  
  it('renders region options correctly when regions are provided', () => {
    const regions = [
      { name: '美国东部 (弗吉尼亚北部)', value: 'us-east-1' },
      { name: '美国东部 (俄亥俄)', value: 'us-east-2' },
      { name: '美国西部 (加利福尼亚北部)', value: 'us-west-1' }
    ]
    
    const wrapper = shallowMount(RegionSelector, {
      props: {
        cloudProvider: 'aws',
        initialRegions: regions
      },
      data() {
        return {
          loading: false,
          regions: regions
        }
      }
    })
    
    // 验证显示区域选项
    expect(wrapper.findAll('.region-option').length).toBe(3)
    expect(wrapper.text()).toContain('美国东部 (弗吉尼亚北部)')
    expect(wrapper.text()).toContain('美国东部 (俄亥俄)')
    expect(wrapper.text()).toContain('美国西部 (加利福尼亚北部)')
  })
  
  it('selects region when clicked and emits region-selected event', async () => {
    const regions = [
      { name: '美国东部 (弗吉尼亚北部)', value: 'us-east-1' },
      { name: '美国东部 (俄亥俄)', value: 'us-east-2' }
    ]
    
    const wrapper = shallowMount(RegionSelector, {
      props: {
        cloudProvider: 'aws',
        initialRegions: regions
      },
      data() {
        return {
          loading: false,
          regions: regions
        }
      }
    })
    
    // 初始状态下没有选中的区域
    expect(wrapper.vm.selectedRegion).toBe('')
    
    // 点击第一个区域选项
    await wrapper.findAll('.region-option').at(0).trigger('click')
    
    // 验证选中的区域是us-east-1
    expect(wrapper.vm.selectedRegion).toBe('us-east-1')
    
    // 验证第一个区域选项有selected类
    expect(wrapper.findAll('.region-option').at(0).classes()).toContain('selected')
    
    // 验证触发了region-selected事件
    expect(wrapper.emitted('region-selected')).toBeTruthy()
    expect(wrapper.emitted('region-selected')[0]).toEqual(['us-east-1'])
  })
  
  it('loads regions when cloudProvider changes', async () => {
    const wrapper = shallowMount(RegionSelector, {
      props: {
        cloudProvider: '',
        initialRegions: []
      }
    })
    
    // 模拟loadRegions方法
    wrapper.vm.loadRegions = jest.fn()
    
    // 更新cloudProvider属性
    await wrapper.setProps({ cloudProvider: 'aws' })
    
    // 验证调用了loadRegions方法
    expect(wrapper.vm.loadRegions).toHaveBeenCalledWith('aws')
  })
  
  it('emits azs-loaded event after region is selected', async () => {
    const regions = [
      { name: '美国东部 (弗吉尼亚北部)', value: 'us-east-1' }
    ]
    
    const wrapper = shallowMount(RegionSelector, {
      props: {
        cloudProvider: 'aws',
        initialRegions: regions
      },
      data() {
        return {
          loading: false,
          regions: regions
        }
      }
    })
    
    // 模拟fetchAvailabilityZones方法
    wrapper.vm.fetchAvailabilityZones = jest.fn()
    
    // 点击区域选项
    await wrapper.findAll('.region-option').at(0).trigger('click')
    
    // 验证调用了fetchAvailabilityZones方法
    expect(wrapper.vm.fetchAvailabilityZones).toHaveBeenCalledWith('aws', 'us-east-1')
  })
})
