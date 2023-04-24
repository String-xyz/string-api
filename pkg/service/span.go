package service

import (
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func StartfromContext(ctx context.Context, operationName string) (tracer.Span, context.Context) {
	return tracer.StartSpanFromContext(ctx, operationName)
}

func StartSpan(operationName string) tracer.Span {
	return tracer.StartSpan(operationName)
}

func Span(ctx context.Context, operationName string, tag string, tagValue string) (tracer.Span, func(options ...tracer.FinishOption)) {
	sp, _ := StartfromContext(ctx, operationName)
	sp.SetTag(tag, tagValue)
	return sp, sp.Finish
}
