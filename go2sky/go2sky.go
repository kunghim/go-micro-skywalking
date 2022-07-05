package go2sky

import (
	"context"
	"errors"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"time"
)

const (
	IntervalTime = 5 * time.Second

	EmptySpanName = ""
)

type Trace struct {
	Reporter go2sky.Reporter
	Tracer   *go2sky.Tracer
}

func NewGRPCReporter(skyAgentAddr, traceName string) (*Trace, error) {
	if len(skyAgentAddr) == 0 || len(traceName) == 0 {
		return nil, errors.New("不能为空")
	}
	grpcReporter, err := reporter.NewGRPCReporter(skyAgentAddr, reporter.WithCheckInterval(IntervalTime))
	if err != nil {
		return nil, err
	}
	tracer, err := go2sky.NewTracer(traceName, go2sky.WithReporter(grpcReporter))
	if err != nil {
		return nil, err
	}
	go2sky.SetGlobalTracer(tracer)
	return &Trace{Reporter: grpcReporter, Tracer: tracer}, nil
}

func (t *Trace) CreateSpan(context context.Context, spanName string) (span go2sky.Span, ctx context.Context, err error) {
	if len(spanName) == 0 {
		spanName = EmptySpanName
	}
	span, ctx, err = t.Tracer.CreateLocalSpan(context)
	if err != nil {
		return nil, nil, err
	}
	span.SetOperationName(spanName)
	return
}
