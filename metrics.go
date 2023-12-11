package diameter

import (
	"time"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

type DiameterMetrics struct {
	RequestDuration    *metrics.Metric
	RequestCount       *metrics.Metric
	FailedRequestCount *metrics.Metric
}

func registerMetrics(vu modules.VU) DiameterMetrics {
	registry := vu.InitEnv().Registry
	metrics := DiameterMetrics{
		RequestDuration:    registry.MustNewMetric("diameter_req_duration", metrics.Trend, metrics.Time),
		RequestCount:       registry.MustNewMetric("diameter_req_count", metrics.Counter, metrics.Default),
		FailedRequestCount: registry.MustNewMetric("diameter_req_failed", metrics.Rate, metrics.Default),
	}
	return metrics
}

func (c *DiameterClient) reportMetric(metric *metrics.Metric, now time.Time, value float64, tags map[string]string) {
	state := c.vu.State()
	ctx := c.vu.Context()
	if state == nil || ctx == nil {
		return
	}

	metrics.PushIfNotDone(ctx, state.Samples, metrics.Sample{
		Time: now,
		TimeSeries: metrics.TimeSeries{
			Metric: metric,
			Tags:   metrics.NewRegistry().RootTagSet().WithTagsFromMap(tags),
		},
		Value: value,
	})
}
