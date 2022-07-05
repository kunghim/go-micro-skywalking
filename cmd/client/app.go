package main

import (
	"context"
	"flag"
	"github.com/SkyAPM/go2sky"
	micro4sky "github.com/SkyAPM/go2sky-plugins/micro"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/asim/go-micro/v3"
	cons "go-micro-skywalking/constant"
	"go-micro-skywalking/proto/hello"
	"log"
	"time"
)

var skyWalkingUrl string

func init() {
	flag.StringVar(&skyWalkingUrl, "a", cons.SwServerAddr, "set your skywalking server address")
	flag.Parse()
}

func main() {
	//Use gRPC reporter for production
	r, err := reporter.NewGRPCReporter(skyWalkingUrl, reporter.WithCheckInterval(5*time.Second))
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer r.Close()

	tracer, err := go2sky.NewTracer(cons.ClientTracer, go2sky.WithReporter(r))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}
	// 创建 micro 服务
	service := micro.NewService(
		micro.Address(cons.ClientMicroAddress),
		// 设置 micro 服务名称
		micro.Name(cons.ClientMicroServer),
		// 加入 opentracing 的中间件
		micro.WrapClient(micro4sky.NewClientWrapper(tracer, micro4sky.WithClientWrapperReportTags("ClientMicroServer"))),
	)
	// 初始化 micro 服务
	service.Init()
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		log.Fatal("创建 ClientSpan 失败 -> ", err)
		return
	}
	span.SetOperationName("ClientSpan")
	defer span.End()
	name := "张三"
	span.Tag("name", name)
	helloWorldService := hello.NewHelloWorldService(cons.HelloMicroServer, service.Client())
	callResponse, err := helloWorldService.Call(ctx, &hello.CallRequest{Name: name})
	if err != nil {
		log.Fatal("调用 notice 服务的 send 接口失败 -> ", err)
		return
	}
	log.Println("执行成功 -> ", callResponse.Msg)
}
