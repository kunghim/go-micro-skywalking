package main

import (
	micro4sky "github.com/SkyAPM/go2sky-plugins/micro"
	"github.com/asim/go-micro/v3"
	cons "go-micro-skywalking/constant"
	"go-micro-skywalking/go2sky"
	"go-micro-skywalking/handler"
	"go-micro-skywalking/proto/hello"
	"go-micro-skywalking/proto/notice"
	"log"
)

func main() {
	trace, err := go2sky.NewGRPCReporter(cons.SwServerAddr, cons.HelloTracer)
	if err != nil {
		log.Fatalf("new trace error %v \n", err)
		return
	}
	defer trace.Reporter.Close()
	// 创建 micro 服务
	service := micro.NewService(
		micro.Address(cons.HelloMicroAddress),
		// 设置 micro 服务名称
		micro.Name(cons.HelloMicroServer),
		// 加入 opentracing 的中间件
		micro.WrapHandler(micro4sky.NewHandlerWrapper(trace.Tracer, "HelloMicroServer")),
	)
	// 初始化 micro 服务
	service.Init()

	// 获取 micro-notice 服务的 noticeService，才能在 Call 中调用 notice send
	noticeService := notice.NewNoticeService(cons.NoticeMicroServer, service.Client())
	err = hello.RegisterHelloWorldHandler(service.Server(), handler.HelloService{Trace: trace, NoticeServer: noticeService})
	if err != nil {
		log.Fatal("注册 hello service 失败 -> ", err)
		return
	}

	// 启动服务
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
