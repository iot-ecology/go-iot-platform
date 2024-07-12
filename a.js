var b = "function main(jsonData) {\n" + "    var c = []\n" + "    for (var jsonDatum of jsonData) {\n" + "        var time = jsonDatum.Time;\n" + "        var arr = []\n" + "        var timeField = {\n" + "            \"FieldName\": \"time\",\n" + "            \"Value\": time\n" + "        }\n" + "\n" + "        arr.push(timeField)\n" + "        var idd = {\n" + "            \"FieldName\": \"id\",\n" + "            \"Value\": time\n" + "        }\n" + "        arr.push(idd)\n" + "        for (var e of jsonDatum.DataRows) {\n" + "            if (e.Name == \"a\") {\n" + "                var aField = {\n" + "                    \"FieldName\": \"name\",\n" + "                    \"Value\": e.Value\n" + "                }\n" + "                arr.push(aField)\n" + "            }\n" + "        }\n" + "        c.push(arr)\n" + "    }\n" + "    return c;\n" + "}"

var a = [{
    "time": 1720751132, "device_uid": "111", "data": [{
        "name": "a", "value": "测试"
    }], "nc": "111"
}];
var param = {"data_row_list": a, "script": b}

console.log(param)