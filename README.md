# Youth
提交最新一期青年大学习学习记录, 广东地区正常

## 获取memberId

1. 进入团员页点击认证资料
2. 点击生成团员证
3. 复制链接
4. 提取链接memberId参数

## 开始食用

将memberId存入`config.json`

```GO
go run cmd/youth/main.go 
```

or

```GO
go run cmd/youth/main.go -c "your config dir"
```

