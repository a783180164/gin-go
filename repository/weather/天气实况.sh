#!/usr/bin/env bash

ollama run qwen:1.8b "
<|system|>
你是一个天气专家，可以告诉人们未来的天气。你熟悉使用 Seniverse 气象台的 API。

<|user|>
天气实况接口：https://api.seniverse.com/v3/weather/now.json  
请求方法：GET  
通用参数：  
- key: S2IiZWW4U4Zq-9jrx  
- location: 查询地点的汉语拼音  

下面是该接口返回的 JSON 模板，请把原始字段封装成更可读的文本并返回给用户。例如输出格式可以这样：
  
“📍 城市：西雅图 (US)  
🌡️ 温度：14°C (体感 14°C)  
☁️ 天气：多云  
💧 湿度：76%  
💨 风：西北 8.05 km/h (2 级)  
…  
数据更新时间：2015-09-25T22:45:00-07:00”

返回前，请先请求接口并解析数据，然后生成可读的输出。
  
以下是示例返回的数据结构：
\`\`\`json
{
  "results": [
    {
      "location": {
        "id": "C23NB62W20TF",
        "name": "西雅图",
        "country": "US",
        "path": "西雅图,华盛顿州,美国",
        "timezone": "America/Los_Angeles",
        "timezone_offset": "-07:00"
      },
      "now": {
        "text": "多云",
        "code": "4",
        "temperature": "14",
        "feels_like": "14",
        "pressure": "1018",
        "humidity": "76",
        "visibility": "16.09",
        "wind_direction": "西北",
        "wind_direction_degree": "340",
        "wind_speed": "8.05",
        "wind_scale": "2",
        "clouds": "90",
        "dew_point": "-12"
      },
      "last_update": "2015-09-25T22:45:00-07:00"
    }
  ]
}
\`\`\`
"
