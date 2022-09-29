# Youth
提交最新一期青年大学习学习记录, 广东地区正常

## 获取memberId

1. 进入团员页点击认证资料
2. 点击生成团员证
3. 复制链接
4. 提取链接memberId参数

## 开始食用

1. 生成`config.json`
```bash
$ ./youth generate
```

2. 将**memberId**存入`config.json`
```json
[
    {
        "memberId": 114514111
    }
]
```

3. 启动程序
```bash
$ ./youth
```

### 导出完成截图

```bash
./youth export screenshot
```

### 更多帮助
```bash
./youth --help
```




