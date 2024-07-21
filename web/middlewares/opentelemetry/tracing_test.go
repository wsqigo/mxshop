package opentelemetry

import (
	"awesomeProject/web"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"testing"
	"time"
)

func TestOpenTelemetry(t *testing.T) {
	ctx := context.Background()
	tracer := otel.GetTracerProvider().Tracer("awesomeProject/web")

	// 如果 ctx 已经和一个 span 绑定了，那么新的 span 就是老的 span 的子 span
	ctx, span := tracer.Start(ctx, "opentelemetry-demo",
		trace.WithAttributes(attribute.String("version", "1")))
	defer span.End()

	// 重置名字
	span.SetName("otel-name")
	span.SetAttributes(attribute.Int("status", 200))
	span.AddEvent("马老师，发生什么事了")
}

func initZipkin(t *testing.T) {
	// 开启 zipkin
	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		t.Fatal(err)
	}

	// 注册 exporter
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("awesomeProject/web"),
		)),
		sdktrace.WithBatcher(exporter), // 传入 exporter, span 数据会被导出到 zipkin
	)

	// 设置全局的 TracerProvider
	otel.SetTracerProvider(provider)
}

func initJeager(t *testing.T) {
	// 开启 jaeger
	url := "http://localhost:14268/api/traces"
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		t.Fatal(err)
	}
	// 注册 exporter
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exporter),
		// Record information about this application in an Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("awesomeProject/web"),
			attribute.String("environment", "test"),
			attribute.Int64("ID", 1),
		)),
	)

	otel.SetTracerProvider(tp)
}

type User struct {
	ID   int
	Name string
}

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer("awesomeProject/web")
	builder := &MiddlewareBuilder{tracer: tracer}
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))

	// Add route
	server.Get("/user", func(ctx *web.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		secondC, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		_, third1 := tracer.Start(secondC, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()
		_, third2 := tracer.Start(secondC, "third_layer_2")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()

		_, span2 := tracer.Start(ctx.Req.Context(), "first_layer_2")
		defer span2.End()
		ctx.RespJson(http.StatusOK, &User{
			ID:   1,
			Name: "Amos",
		})
	})

	initZipkin(t)

	server.Start(":8081")
}
