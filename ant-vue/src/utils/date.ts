import dayjs from "dayjs";

// 获取日期时间范围
type GetDateRangeType = <T>(config: {
  size: number; // 时间大小
  unit: string; // 时间单位
  start?: boolean; // 是否取时间单位的开始
  end?: boolean; // 是否取时间单位的末尾 false=当前时间
  future?: boolean; // 结束时间是否是未来时间
  output?: string; // 输出格式默认时间戳 dayjs/timestamp/日期格式
}) => T[];
export const getDateRange: GetDateRangeType = (config) => {
  const size = config.size;
  const unit = config.unit as dayjs.ManipulateType;
  const start = Boolean(config.start);
  const end = Boolean(config.end);
  const future = Boolean(config.future);
  const output = config.output ?? "timestamp";
  // 开始时间
  let startTime: any = dayjs().subtract(size, unit);
  // 是否取单位时间的开始
  if (start) startTime = startTime.startOf(unit);
  // 结束时间
  let endTime: any = dayjs();
  // 结束时间是否也是过去时间
  if (end && !future) endTime = endTime.subtract(size, unit);
  // 是否取单位时间的结束
  if (end) endTime = endTime.endOf(unit);
  // 设置输出格式
  switch (output) {
    case "dayjs":
      break;
    case "timestamp":
      startTime = startTime.valueOf();
      endTime = endTime.valueOf();
      break;
    default:
      startTime = startTime.format(output);
      endTime = endTime.format(output);
      break;
  }
  return [startTime, endTime];
};

/**
 * 转换日期格式
 * @param source 源数据   支持格式Array/String/Number/Dayjs
 * @param output 输出格式 默认string string/timestamp/dayjs
 * @param format 时间格式 默认YYYY-MM-DD HH:mm:ss
 */
export const dateFormatter = <T>(source: any, output?: "string" | "timestamp" | "dayjs", format?: string): T => {
  const isArray = Array.isArray(source);
  let result = [];
  const _source = isArray ? source : [source];
  const _output = output ?? "string";
  const _format = format ?? "YYYY-MM-DD HH:mm:ss";
  result = _source.map((item) => {
    const _item = isNaN(Number(item)) ? item : Number(item);
    if (!_item) return "";
    if (_output === "string") {
      return dayjs(_item).format(_format);
    } else if (_output === "timestamp") {
      return dayjs(_item).valueOf();
    } else {
      return dayjs(_item);
    }
  });
  return isArray ? result : (result[0] as any);
};
