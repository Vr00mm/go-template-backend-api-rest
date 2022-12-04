package metrics

import (
	"net/http"
	"time"

	"github.com/nbari/violetear"
	"github.com/nbari/violetear/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func MetricsMWCallBack(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do something here
		next.ServeHTTP(w, r)
	})
}

func MetricsMW(c *prometheus.HistogramVec) middleware.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// do something here
			next.ServeHTTP(w, r)
			endpoint := violetear.GetRouteName(r)
			//method := prometheus.Labels{"method": r.Method}
			c.WithLabelValues(endpoint, r.Method).Observe(time.Since(start).Seconds())
			//c.With(method)

		})
	}
}
