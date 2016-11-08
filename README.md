#   分词服务 

### 运行方式

```
    segment [dictPath] [ip:port]
```

---

### 获取分词

`POST  http://segment.services.com/segment`

> 参数

| 名称        | 类型           | 描述  |
| ------------- |:-------------:| -----:|
| content     | string | 分词内容 |
| postagging      | bool      |   是否标注词性 |

### 用例

`curl -d "content=我是中国人;postagging=true" http://127.0.0.1:6800/segment`
