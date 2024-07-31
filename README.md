# Fuxi

Fuxi 是一款 Gin 框架开发辅助工具。

> Fuxi（伏羲）名字来源于中国古代神话传说。

## Goals
通过Fuxi工具，帮助开发者在使用 Gin 框架进行项目开发时，能够更加得心应手，
提升开发效率和代码质量, 从而能够助力开发者创造出卓越的软件产品。

## Getting Started
### Required
- [golang](https://golang.org/dl/)

### Installing
##### go install 安装：
```
go install github.com/go-zxb/fuxi@latest
fuxi --version
```
##### 源码编译安装：
```
git clone https://github.com/go-zxb/fuxi
cd fuxi
go mod tidy
go install
```

### Create a project
```go

# 创建项目模板
fuxi project -n helloworld

# 生成xxCRUD源码
fuxi api:new -n user -q 帮我设计一个用户管理模块
# -n user 模块名
# -q 你的需求, 例如：帮我设计一个用户管理模块(提交后会请求xx大模型生成相关数据)

# 添加一个xx接口源码
# -a getUserName 路由地址
# -f getUserName 方法名
fuxi api:add -n user -a getUserName -f getUserName

# 运行程序
fuxi run
```

## Wechat
![Fuxi](docs/images/wechat.jpg)

## 贡献
我们欢迎任何形式的贡献，包括但不限于代码提交、问题反馈、功能建议等。

## 许可证
本项目采用 MIT 许可证，详情请参见 LICENSE 文件。