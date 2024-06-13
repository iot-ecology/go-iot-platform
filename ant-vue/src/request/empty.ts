// 不需要处理空值情况的数组
export const notDealEmpty: string[] = [];

// 处理空值情况
// [] => null、{} => null、"" => null
export const replaceEmptyValue = (value: any): any => {
  if (value === null) return null;
  else if (value === undefined) return undefined;
  else if (typeof value === "string" && !value) return null;
  else if (Array.isArray(value)) {
    if (value.length === 0) {
      return null;
    } else {
      return value.map((item) => replaceEmptyValue(item));
    }
  } else if (typeof value === "object") {
    const _length = Object.keys(value).length;
    // 兼容就是需要 {} 这种情况
    if (_length === 1 && value.__empty__ === true) {
      return {};
    } else if (_length === 0) {
      return null;
    } else {
      const newValue: Record<any, any> = {};
      for (const key in value) {
        if (Object.prototype.hasOwnProperty.call(value, key)) {
          newValue[key] = replaceEmptyValue(value[key]);
        }
      }
      return newValue;
    }
  } else return value;
};
