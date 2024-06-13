/**
 * 获取元素的样式属性
 * @param {HTMLElement} el 需要获取样式的dom
 * @param {string} attr 需要获取的属性名
 */
export const getElStyle = (el: HTMLElement, attr: string): any => {
  const style = window.getComputedStyle(el, null) as Record<string, any>;
  return style[attr];
};

/**
 * 获取元素的padding信息
 * @param {HTMLElement} el 需要获取padding的dom
 */
export type GetElPadding = (el: HTMLElement) => { left: number; right: number; top: number; bottom: number };
export const getElPadding: GetElPadding = (el) => {
  const style = window.getComputedStyle(el, null);
  const left = Number.parseInt(style.paddingLeft, 10) || 0;
  const right = Number.parseInt(style.paddingRight, 10) || 0;
  const top = Number.parseInt(style.paddingTop, 10) || 0;
  const bottom = Number.parseInt(style.paddingBottom, 10) || 0;
  return { left, right, top, bottom };
};

/**
 * 获取元素的margin信息
 * @param {HTMLElement} el 需要获取margin的dom
 */
export type GetElMargin = (el: HTMLElement) => { left: number; right: number; top: number; bottom: number };
export const getElMargin: GetElMargin = (el) => {
  const style = window.getComputedStyle(el, null);
  const left = Number.parseInt(style.marginLeft, 10) || 0;
  const right = Number.parseInt(style.marginRight, 10) || 0;
  const top = Number.parseInt(style.marginTop, 10) || 0;
  const bottom = Number.parseInt(style.marginBottom, 10) || 0;
  return { left, right, top, bottom };
};

/**
 * 判断元素内文本是否溢出了
 * @param {HTMLElement} el 需要判断的dom
 */
export const isOverflow = (el: HTMLElement): boolean => {
  const range = document.createRange();
  range.setStart(el, 0);
  range.setEnd(el, el.childNodes.length);
  const rangeWidth = Math.round(range.getBoundingClientRect().width);
  const { left, right } = getElPadding(el);
  const horizontalPadding = left + right;
  if (rangeWidth + horizontalPadding > el.offsetWidth) {
    return true;
  } else if (el.scrollWidth > el.offsetWidth) {
    return true;
  } else {
    return false;
  }
};
