# 读取 txt 文件，每一行都是字符串，加入到数组中,去掉换行符
def read_txt(file_path):
    with open(file_path, "r") as f:
        lines = f.readlines()
        for i in range(len(lines)):
            lines[i] = lines[i].strip()
        return lines

nodebind = read_txt("1.txt")

mqtt_config = read_txt("2.txt")

# 对比 nodebind 和 mqtt_config 求差异元素
def compare(nodebind, mqtt_config):
    nodebind_set = set(nodebind)
    print(nodebind_set)
    mqtt_config_set = set(mqtt_config)
    print(mqtt_config_set)
    return mqtt_config_set.difference(nodebind_set)

print(compare(nodebind, mqtt_config))