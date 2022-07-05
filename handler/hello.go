package handler

import (
	"context"
	"go-micro-skywalking/go2sky"
	"go-micro-skywalking/proto/hello"
	"go-micro-skywalking/proto/notice"
	"log"
	"time"
)

type HelloService struct {
	NoticeServer notice.NoticeService
	Trace        *go2sky.Trace
}

func (h HelloService) Call(ctx context.Context, request *hello.CallRequest, response *hello.CallResponse) error {
	log.Println("this is HelloService.Call")
	response.Msg = "Hello, " + request.Name
	span, ctx, err := h.Trace.CreateSpan(ctx, "HelloSpan")
	if err != nil {
		log.Fatal("创建 HelloSpan 失败 -> ", err)
		return err
	}
	defer span.End()
	span.Tag("CallRequest.Name", request.Name)
	span.Log(time.Now(), "this is HelloService.Call", "receive client request success")
	// 调用 notice 服务的 send 接口
	sendResponse, err := h.NoticeServer.Send(ctx, &notice.SendRequest{Name: request.Name})
	if err != nil {
		return err
	}
	log.Println("执行成功 -> ", sendResponse.Msg)
	// 模拟执行成功（可忽略）
	time.Sleep(1 * time.Second)
	return nil
}
