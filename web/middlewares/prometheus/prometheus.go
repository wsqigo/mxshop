package prometheus

import (
	"awesomeProject/web"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Name        string
	SubSystem   string
	ConstLabels map[string]string
	Help        string
}

func NewMiddlewareBuilder(name, subsystem string, constLabels map[string]string, help string) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		Name:        name,
		SubSystem:   subsystem,
		ConstLabels: constLabels,
		Help:        help,
	}
}

func (b *MiddlewareBuilder) Build() web.Middleware {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        b.Name,
		Subsystem:   b.SubSystem,
		ConstLabels: b.ConstLabels,
		Help:        b.Help,
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(summaryVec)
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx *web.Context) {
			startTime := time.Now()
			go func() {
				dur := time.Since(startTime)
				summaryVec.WithLabelValues(ctx.MatchedRoute, ctx.Req.Method, strconv.Itoa(ctx.RespStatusCode)).
					Observe(float64(dur.Milliseconds()))
			}()
			next(ctx)
		}
	}
}
