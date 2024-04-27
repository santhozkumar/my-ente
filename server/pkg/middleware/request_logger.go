package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/santhozkumar/my-ente/pkg/utils/network"
	timeutil "github.com/santhozkumar/my-ente/pkg/utils/time"
)

var latency = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "museum_latency",
	Help:    "museum_latency",
	Buckets: []float64{10, 50, 100, 200, 500, 1000, 10000, 30000, 60000, 120000, 600000},
}, []string{"code", "method", "host", "url"})

func Logger(urlSanitizer func(_ *gin.Context) string) gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		reqID := requestid.Get(c)
		reqMethod := c.Request.Method
		clientip := network.GetClientIP(c)
		clientPackage := c.GetHeader("X-Client-Package")
		clientVersion := c.GetHeader("X-Client-Version")
		queryValues := c.Request.URL.Query()
		if queryValues.Has("token") {
			queryValues.Set("token", "redacted-token")
		}
		queryParamsForLog := queryValues.Encode()
		reqContextLogger := slog.New(slog.Default().Handler().WithAttrs(
			[]slog.Attr{
				slog.String("client_ip", clientip),
				slog.String("client_package", clientPackage),
				slog.String("client_version", clientVersion),
				slog.String("request_id", reqID),
				slog.String("request_method", reqMethod),
				slog.String("request_uri", c.Request.URL.Path),
				slog.String("query_params", queryParamsForLog),
				slog.String("user_agent", c.GetHeader("User-Agent")),
			}))

		reqContextLogger.Log(context.Background(), slog.LevelInfo, "Incoming")
		c.Next()
		latencyTime := time.Since(start)
		statusCode := c.Writer.Status()
		reqURI := urlSanitizer(c)
		if reqMethod != http.MethodOptions {
			slog.Info("sending prometheus")
			latency.WithLabelValues(
				strconv.Itoa(statusCode),
				reqMethod,
				c.Request.Host,
				reqURI).Observe(float64(latencyTime.Milliseconds()))
			// do the prometheus logging
		}

		reqContextLogger.LogAttrs(
			context.Background(),
			slog.LevelInfo, "Incoming",
			slog.String("latency_time", latencyTime.String()),
			slog.String("h_latency", timeutil.HumanFriendlyDuration(latencyTime)),
			slog.Int("status_code", statusCode),
		)
	}
}
