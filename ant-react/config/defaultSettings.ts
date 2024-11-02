import { ProLayoutProps } from '@ant-design/pro-components';

/**
 * @name
 */
const Settings: ProLayoutProps & {
  pwa?: boolean;
  logo?: string;
} = {
  navTheme: 'light',
  // 拂晓蓝
  colorPrimary: '#1890ff',
  layout: 'mix',
  contentWidth: 'Fluid',
  fixedHeader: false,
  fixSiderbar: true,
  colorWeak: false,
  title: 'Go IoT 开发平台',
  pwa: true,
  logo: 'https://gw.alipayobjects.com/zos/rmsportal/KDpgvguMpGfqaHPjicRK.svg',
  iconfontUrl: '',
  token: {// 布局背景颜色
    bgLayout: '#eaeff3', // 浅蓝灰色
// 左侧菜单栏配置
    sider: {
      // 折叠图标背景颜色
      colorBgCollapsedButton: '#ffffff', // 白色
      // 鼠标移动到箭头上悬浮颜色
      colorTextCollapsedButtonHover: '#6c63ff', // 明亮的紫蓝色
      // 折叠箭头图标颜色
      colorTextCollapsedButton: '#333333', // 深灰色
      // 菜单背景颜色
      colorMenuBackground: '#ffffff', // 白色
      // 菜单分割线的颜色
      colorMenuItemDivider: '#d1d1d1', // 浅灰色
      // 菜单悬浮颜色
      colorBgMenuItemHover: '#f0f8ff', // 浅蓝色
      // 菜单处于激活状态的颜色
      colorBgMenuItemActive: '#d1e7fd', // 浅蓝色
      // 菜单选中后的颜色
      colorBgMenuItemSelected: '#d1e7fd', // 浅蓝色
      // 菜单选中后字体颜色
      colorTextMenuSelected: '#333333', // 深灰色
      // 鼠标悬浮后item字体颜色
      colorTextMenuItemHover: '#6c63ff', // 明亮的紫蓝色
      // 菜单处于激活状态的字体颜色
      colorTextMenuActive: '#333333', // 深灰色
      // 菜单字体颜色
      colorTextMenu: '#333333', // 深灰色
      /**
       * menu 顶部 title 的字体颜色
       */
      colorTextMenuTitle: '#333333', // 深灰色
    },
// 头部颜色配置
    header: {
      // 头部背景颜色
      colorBgHeader: '#ffffff', // 白色
      // 头部标题颜色
      colorHeaderTitle: '#333333', // 深灰色
      // 字体颜色
      colorTextMenu: '#333333', // 深灰色
      // 二级菜单字体颜色
      colorTextMenuSecondary: '#6c63ff', // 明亮的紫蓝色
      // 右上角鼠标悬浮字体背景颜色
      colorBgRightActionsItemHover: '#f0f8ff', // 浅蓝色
      // 右上角字体颜色
      colorTextRightActionsItem: '#333333', // 深灰色
    },
    /**
     * 内容主体颜色配置
     */
    pageContainer: {
      /**
       * pageContainer 的背景颜色
       */
      colorBgPageContainer: '#f9f9f9', // 浅灰色
      /**
       * pageContainer 自带的 margin inline
       * @deprecated 请使用 paddingInlinePageContainerContent
       */
      colorBgPageContainerFixed: '#d1d1d1', // 浅灰色
    },
  },
};

export default Settings;
