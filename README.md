# 待办事项应用

使用Go语言开发的简单待办事项管理系统。

## 技术栈

后端：
- Go
- Gin框架
- SQLite数据库

前端：
- HTML
- CSS
- JavaScript

## 功能

- 查看待办事项列表
- 添加新的待办事项
- 更新待办事项状态
- 删除待办事项

## 使用方法

1. 启动服务器：
```
go run cmd/main.go
```

2. 访问 http://localhost:8080/static/index.html

## 目录结构

- cmd/main.go - 程序入口
- internal/ - 内部包
  - database/ - 数据库操作
  - handler/ - HTTP处理器
  - model/ - 数据模型
  - router/ - 路由配置
- static/ - 前端文件
