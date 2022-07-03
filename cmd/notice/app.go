package main

import (
	"flag"
	"github.com/SkyAPM/go2sky"
	micro4sky "github.com/SkyAPM/go2sky-plugins/micro"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/asim/go-micro/v3"
	cons "go-micro-skywalking/constant"
	"go-micro-skywalking/handler"
	"go-micro-skywalking/proto/notice"
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

	tracer, err := go2sky.NewTracer(cons.NoticeTracer, go2sky.WithReporter(r))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}

	// 创建 micro 服务
	service := micro.NewService(
		// 设置 micro 服务名称
		micro.Name(cons.NoticeMicroServer),
		// 加入 opentracing 的中间件
		micro.WrapHandler(micro4sky.NewHandlerWrapper(tracer, "NoticeMicroServer")),
	)
	// 初始化 micro 服务
	service.Init()

	err = notice.RegisterNoticeServiceHandler(service.Server(), handler.NoticeService{Tracer: tracer})
	if err != nil {
		log.Fatal("注册 notice service 失败 -> ", err)
		return
	}

	// 启动服务
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
