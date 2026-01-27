package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const TraceIdKey = "traceId"

func TraceMiddleware() func(httpCtx interface{}, next func()) {
	rand.Seed(time.Now().UnixNano())

	return func(httpCtx interface{}, next func()) {
		ctx := httpCtx.(context.Context)

		//尝试从请求头获取traceId,否则生成一个新的
		traceId := generateTraceId()

		//将traceId存入context中
		ctx = context.WithValue(ctx, TraceIdKey, traceId)

		next()
	}
}

func generateTraceId() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
}
