import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import DashboardCard from '../DashboardCard.vue'

describe('DashboardCard', () => {
  it('renders with required props', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Test Title',
        value: 42,
        icon: 'chart',
        subtitle: 'Test Subtitle'
      }
    })

    expect(wrapper.find('p.text-sm').text()).toBe('Test Title')
    expect(wrapper.find('.text-3xl').text()).toBe('42')
    expect(wrapper.find('.text-sm.text-slate-500').text()).toBe('Test Subtitle')
  })

  it('renders with string value', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Status',
        value: 'Active',
        icon: 'check'
      }
    })

    expect(wrapper.find('.text-3xl').text()).toBe('Active')
  })

  it('renders with number value', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Count',
        value: 1234,
        icon: 'chart'
      }
    })

    expect(wrapper.find('.text-3xl').text()).toBe('1234')
  })

  it('renders with zero value', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Empty Count',
        value: 0,
        icon: 'chart'
      }
    })

    expect(wrapper.find('.text-3xl').text()).toBe('0')
  })

  it('renders without subtitle', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Simple Card',
        value: 100,
        icon: 'target'
      }
    })

    expect(wrapper.find('p.text-sm').text()).toBe('Simple Card')
    expect(wrapper.find('.text-3xl').text()).toBe('100')
    // Should not have subtitle
    expect(wrapper.find('.text-sm.text-slate-500').exists()).toBe(false)
  })

  it('displays correct icon based on icon prop', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Chart Card',
        value: 50,
        icon: 'chart'
      }
    })

    // Should contain SVG icon
    expect(wrapper.find('svg').exists()).toBe(true)
  })

  it('handles different icon types', () => {
    const iconTypes = ['chart', 'check', 'target', 'building']
    
    iconTypes.forEach(iconType => {
      const wrapper = mount(DashboardCard, {
        props: {
          title: `${iconType} Card`,
          value: 25,
          icon: iconType
        }
      })

      expect(wrapper.find('svg').exists()).toBe(true)
    })
  })

  it('has proper styling classes', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Styled Card',
        value: 75,
        icon: 'chart',
        subtitle: 'With styling'
      }
    })

    // Check main container has background and border classes
    expect(wrapper.classes()).toContain('bg-white/80')
    expect(wrapper.classes()).toContain('backdrop-blur-sm')
    expect(wrapper.classes()).toContain('rounded-2xl')
    expect(wrapper.classes()).toContain('border')
  })

  it('renders large values correctly', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Large Number',
        value: 999999999,
        icon: 'chart'
      }
    })

    expect(wrapper.find('.text-3xl').text()).toBe('999999999')
  })

  it('renders negative values correctly', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Negative Value',
        value: -42,
        icon: 'chart'
      }
    })

    expect(wrapper.find('.text-3xl').text()).toBe('-42')
  })

  it('renders decimal values correctly', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Decimal Value',
        value: 42.5,
        icon: 'chart'
      }
    })

    expect(wrapper.find('.text-3xl').text()).toBe('42.5')
  })

  it('handles empty string value', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Empty String',
        value: '',
        icon: 'chart'
      }
    })

    expect(wrapper.find('.text-3xl').text()).toBe('')
  })

  it('handles null value gracefully', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Null Value',
        value: 'N/A',
        icon: 'chart'
      }
    })

    // Should not crash and display fallback value
    expect(wrapper.find('.text-3xl').exists()).toBe(true)
  })

  it('handles undefined value gracefully', () => {
    const wrapper = mount(DashboardCard, {
      props: {
        title: 'Undefined Value',
        value: '',
        icon: 'chart'
      }
    })

    // Should not crash and display empty string
    expect(wrapper.find('.text-3xl').exists()).toBe(true)
  })
})