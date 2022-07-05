package main

import (
	"context"
	micro4sky "github.com/SkyAPM/go2sky-plugins/micro"
	"github.com/asim/go-micro/v3"
	cons "go-micro-skywalking/constant"
	"go-micro-skywalking/go2sky"
	"go-micro-skywalking/proto/hello"
	"log"
)

func main() {
	trace, err := go2sky.NewGRPCReporter(cons.SwServerAddr, cons.ClientTracer)
	if err != nil {
		log.Fatalf("new trace error %v \n", err)
		return
	}
	defer trace.Reporter.Close()
	// 创建 micro 服务
	service := micro.NewService(
		micro.Address(cons.ClientMicroAddress),
		// 设置 micro 服务名称
		micro.Name(cons.ClientMicroServer),
		// 加入 opentracing 的中间件
		micro.WrapClient(micro4sky.NewClientWrapper(trace.Tracer, micro4sky.WithClientWrapperReportTags("ClientMicroServer"))),
	)
	// 初始化 micro 服务
	service.Init()
	span, ctx, err := trace.CreateSpan(context.Background(), "ClientSpan")
	if err != nil {
		log.Fatal("创建 ClientSpan 失败 -> ", err)
		return
	}
	defer span.End()
	name := "张三"
	span.Tag("name", name)
	helloWorldService := hello.NewHelloWorldService(cons.HelloMicroServer, service.Client())
	callResponse, err := helloWorldService.Call(ctx, &hello.CallRequest{Name: name})
	if err != nil {
		log.Fatal("调用 hello 服务的 call 接口失败 -> ", err)
		return
	}
	log.Println("执行成功 -> ", callResponse.Msg)
}
