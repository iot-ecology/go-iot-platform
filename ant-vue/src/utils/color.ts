// 将十六进制的颜色转为rgba格式
export const color2rgba = (color: string, alpha = 1) => {
  const arr = [];
  for (let i = 1; i < 7; i += 2) {
    arr.push(parseInt("0x" + color.slice(i, i + 2)));
  }
  return `rgba(${arr.join(",")},${alpha})`;
};
