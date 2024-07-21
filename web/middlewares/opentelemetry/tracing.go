package opentelemetry

import (
	"awesomeProject/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "awesomeProject/web/middlewares/opentelemetry"
)

type MiddlewareBuilder struct {
	tracer trace.Tracer
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (b *MiddlewareBuilder) Build() web.Middleware {
	if b.tracer == nil {
		b.tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			// 从请求头中提取 traceid 和 spanid
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			// 创建 span
			reqCtx, span := b.tracer.Start(reqCtx, "unknown", trace.WithAttributes())
			// span.End 执行之后，就意味着 span 本身已经确定无疑了，将不能再变化了
			defer span.End()

			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("peer.home", ctx.Req.Host))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr)) // 从请求头中提取客户端信息
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))

			ctx.Req = ctx.Req.WithContext(reqCtx)
			next(ctx)

			if ctx.MatchedRoute != "" {
				span.SetName(ctx.MatchedRoute)
			}
			span.SetAttributes(attribute.Int("http.status_code", ctx.RespStatusCode))
		}
	}
}
