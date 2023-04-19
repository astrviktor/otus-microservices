package middleware

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"time"
)

func Logging(log *zap.Logger, h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		h(ctx)

		log.Info(fmt.Sprintf("[RequestMiddleware] %s %s response", ctx.Method(), ctx.Path()),
			zap.ByteString("body", ctx.Response.Body()),
			zap.ByteString("response_headers", ctx.Response.Header.Header()),
			zap.ByteString("request_headers", ctx.Request.Header.Header()),
			zap.Int("code", ctx.Response.StatusCode()),
			zap.ByteString("request_body", ctx.Request.Body()),
			zap.Int64("duration", time.Since(start).Milliseconds()),
		)

		//log.Println(ip, r.Method, r.RequestURI, r.Proto, recorder.Status, duration, userAgent)
	}
}
