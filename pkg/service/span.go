package service

import (
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type SpanTag map[string]interface{}

func StartfromContext(ctx context.Context, operationName string) (tracer.Span, context.Context) {
	return tracer.StartSpanFromContext(ctx, operationName)
}

func StartSpan(operationName string) tracer.Span {
	return tracer.StartSpan(operationName)
}

func Span(ctx context.Context, operationName string, tags ...SpanTag) (tracer.Span, func(options ...tracer.FinishOption)) {
	sp, _ := StartfromContext(ctx, operationName)
	for _, tag := range tags {
		for key, value := range tag {
			sp.SetTag(key, value)
		}
	}
	return sp, sp.Finish
}
