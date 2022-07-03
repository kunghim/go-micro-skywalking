# go-micro-skywalking
记录个人在学习 opentracing 时做的 DEMO。
目前已经初步了解了 go-micro 集成 jaeger 以及 skywalking，还有很多需要学习的地方期待与大家共同成长。

[go-micro-skywalking](https://github.com/kunghim/go-micro-skywalking)

[go-micro-jaeger](https://github.com/kunghim/go-micro-jaeger)

## 项目背景
公司的项目是用 GoLang 开发的，没有集成链路追踪。作为一个 GoLang 语言的初学者，同时也对公司 Go 项目的业务有诸多不懂，遇到非本人开发的 BUG，定位问题的时候找的我头皮发麻。尤其是线上问题，虽然说都是小问题，但是如果不马上处理就感觉特别慌。
所以迫切需要学习一下 GoLang 集成 opentracing，看到时候能不能在公司项目上引用。

## 安装
将 go-micro-skywalking clone 到本地
```git
git clone https://github.com/kunghim/go-micro-skywalking.git
```
![image](https://user-images.githubusercontent.com/104054614/177053528-9e18b7db-ca09-40bd-877d-27c3ed27031c.png)

## 目录结构描述
- go-micro-skywalking
  - .run 程序运行配置文件
  - cmd 服务启动入口
    - client 客户端（用于向 hello 服务发起请求）
    - hello hello 服务端（接收 client 的请求，处理业务并向 notice 服务端发起请求，实现三级链路）
    - notice notice 服务端（接收 hello 服务端的请求，处理业务并模拟执行失败）
  - constant 静态参数配置（项目启动前修改 SwServerAddr 为自己的 skywalking agent 的连接地址）
  - handler 业务处理模块
  - proto protoc 生成的 go 及 go-micro 文件（主要是 hello 服务以及 notice 服务）

## 使用
### 环境要求
1. golang 版本：1.17
2. 准备好 skywalking 服务
3. go path 已下载好 go-micro（v3）以及 go-gen-proto 等

本项目为 go-micro 集成 skywalking，用到的第三方组件为 [go2sky](https://github.com/SkyAPM/go2sky)，所以你可能需要下载 go2sky
```git
go get -u github.com/asim/go-micro/v3@latest
go get -u github.com/SkyAPM/go2sky@latest
go get -u github.com/SkyAPM/go2sky-plugins/micro@latest
```

本项目已经写好启动配置文件，所以当你打开项目后可以直接运行

![image](https://user-images.githubusercontent.com/104054614/177054099-ceb14126-260d-48ae-bb84-74cf5763034f.png)

当然你可能需要修改 skywalking agent service 的地址，在 constant/const.go 中修改 SwServerAddr 变量
![image](https://user-images.githubusercontent.com/104054614/177054144-fef33f7d-7377-41c9-9340-8107446716be.png)

或者你可以通过携带 skywalking agent service 地址参数的命令启动
进入 /cmd/notice 目录，通过以下命令启动 notice 服务
```shell
go run app.go -a="127.0.0.1:11800"
```
其中 -a 为skywalking agent service 的地址

hello 服务和 client 客户端同理
先启动 notice，然后启动 hello，最后启动 client 进行测试
