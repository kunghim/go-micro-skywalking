package handler

import (
	"context"
	"github.com/SkyAPM/go2sky"
	"go-micro-skywalking/proto/notice"
	"log"
	"time"
)

type NoticeService struct {
	Tracer *go2sky.Tracer
}

func (n NoticeService) Send(ctx context.Context, request *notice.SendRequest, response *notice.SendResponse) error {
	log.Println("this is NoticeService.Send")
	response.Msg = "NoticeService 接收到请求啦"
	span, ctx, err := n.Tracer.CreateLocalSpan(ctx)
	if err != nil {
		log.Fatal("创建 NoticeSpan 失败 -> ", err)
		return err
	}
	defer span.End()
	span.SetOperationName("NoticeSpan")
	// 设置一个 tag
	span.Tag("SendRequest.Name", request.Name)
	span.Log(time.Now(), "this is NoticeService.Send", "receive hello request success")
	span.Error(time.Now(), "测试执行失败", "睡眠 3s...")
	time.Sleep(3 * time.Second)
	return nil
}
