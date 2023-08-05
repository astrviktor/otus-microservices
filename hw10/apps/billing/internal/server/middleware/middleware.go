package middleware

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"otus-microservices/billing/internal/server/prometheus"
	"strconv"
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

		//metrics.GetOrCreateHistogram(
		//	fmt.Sprintf("request_duration_ms{method=%q,path=%q,status=%q}",
		//		string(ctx.Method()),
		//		string(ctx.Path()),
		//		strconv.Itoa(ctx.Response.StatusCode()),
		//	),
		//).UpdateDuration(start)

		prometheus.Metrics.ResponseTime.WithLabelValues(
			string(ctx.Method()),
			string(ctx.Path()),
			strconv.Itoa(ctx.Response.StatusCode()),
		).Observe(float64(time.Since(start).Milliseconds()) / 1000)

		//log.Println(ip, r.Method, r.RequestURI, r.Proto, recorder.Status, duration, userAgent)
	}
}
