import { color2rgba } from "@/utils/color";

// 主题色 #1677FF
const colorPrimary = "44, 135, 255";
// 成功色 #11CABE
const colorSuccess = "17, 202, 190";
// 错误色 #DF4A4A
const colorError = "223, 74, 74";
// 信息色 #1890FF
const colorInfo = "24, 144, 255";
// 警戒色 #FAAD14
const colorWarning = "250, 173, 20";
// 字体系列
const fontFamily = "'Pingfang SC', 'SF UI Text', 'Helvetica Neue', 'Consolas'";

// 项目本身需要的变量
export const cssVariableLight: Record<string, string> = {
  "--g-color-primary": `rgb(${colorPrimary})`, // 主题色
  "--g-color-primary-rgb": colorPrimary, // 主题色三原色
  "--g-color-success": `rgb(${colorSuccess})`, // 成功色
  "--g-color-success-rgb": colorSuccess, // 成功色三原色
  "--g-color-error": `rgb(${colorError})`, // 失败色
  "--g-color-error-rgb": colorError, // 失败色三原色
  "--g-color-info": `rgb(${colorInfo})`, // 信息色
  "--g-color-info-rgb": colorInfo, // 信息色三原色
  "--g-color-warning": `rgb(${colorWarning})`, // 警戒色
  "--g-color-warning-rgb": colorWarning, // 警戒色三原色
  "--g-font-family": fontFamily, // 字体系列
  "--g-font-size": "12px", // 文字大小
  "--g-color-title": color2rgba("#000000", 0.9), // 标题文字颜色
  "--g-color-text": color2rgba("#000000", 0.88), // 正文文字颜色
  "--g-color-desc": color2rgba("#000000", 0.7), // 描述文字颜色
  "--g-color-divider-line": "#F0F0F0", // 分割线颜色

  "--g-layout-bg": "#F5F5F5", // 布局背景色
  "--g-container-bg": "#FFFFFF", // 组件容器背景色
  "--g-container-border": "#F0F0F0", // 组件容器边框色
  "--g-container-border-radius": "4px", // 组件容器边框圆角

  "--g-scrollbar-thumb-bg": "#B4B4B4", // 滚动条颜色
  "--g-scrollbar-thumb-hover-bg": "#9D9D9D", // 滚动条悬停颜色

  "--vxe-font-family": fontFamily, // 字体系列
  "--vxe-font-size": "12px", // 字体大小
  "--vxe-font-color": color2rgba("#000000", 0.88), // 字体颜色
  "--vxe-table-header-font-color": color2rgba("#000000", 0.88), // 表格表头字体颜色
  "--vxe-table-header-font-weight": "600", // 表格表头字体粗细
  "--vxe-table-header-background-color": "transparents", // 表格表头背景颜色
  "--vxe-table-body-background-color": "transparents", // 表格内容区域背景颜色
  "--vxe-table-border-radius": "0px", // 表格边框圆角半径
  "--vxe-table-border-width": "1px", // 表格边框宽度
  "--vxe-table-border-color": "#EBEEF5", // 表格边框宽度
  "--vxe-table-resizable-line-color": "#C0C4CC", // 表格可调整大小线条颜色
  "--vxe-table-resizable-drag-line-color": `rgb(${colorPrimary})`, // 表格可调整大小拖动线条颜色
  "--vxe-table-row-height-default": "38px", // 表格行高度（默认高度）
  "--vxe-table-row-line-height": "38px", // 表格行行高
  "--vxe-table-row-current-background-color": "#ECF5FF", // 表格当前行背景颜色
  "--vxe-table-row-hover-background-color": "#F5F7FA", // 表格鼠标悬停行背景颜色
  "--vxe-table-row-hover-current-background-color": "#F5F7FA", // 表格鼠标悬停当前行背景颜色
  "--vxe-table-row-hover-striped-background-color": "#F5F7FA", // 表格鼠标悬停斑马纹行背景颜色
  "--vxe-table-row-hover-checkbox-checked-background-color": "#F5F7FA", // 表格鼠标悬停复选框选中行背景颜色
  "--vxe-table-row-striped-background-color": "#FAFAFA", // 表格斑马纹行背景颜色
  "--vxe-table-row-checkbox-checked-background-color": "#ECF5FF", // 表格复选框选中行背景颜色
  "--vxe-table-column-padding-default": "8px", // 表格列默认内边距
  "--vxe-table-cell-padding-left": "8px", // 表格单元格左边距
  "--vxe-table-cell-padding-right": "8px", // 表格单元格右边距
  "--vxe-loading-background-color": "#F5F5F5", // 表格loading背景色
};

// ant-design-vue变量
export const antDesignVueLight = {
  token: {
    borderRadius: 4, // 圆角大小
    colorPrimary: `rgb(${colorPrimary})`, // 主题色
    colorSuccess: `rgb(${colorSuccess})`, // 成功色
    colorError: `rgb(${colorError})`, // 错误色
    colorInfo: `rgb(${colorInfo})`, // 信息色
    colorWarning: `rgb(${colorWarning})`, // 警戒色
    fontFamily, // 字体系列
    fontSize: 12, // 字体大小
  },
};
