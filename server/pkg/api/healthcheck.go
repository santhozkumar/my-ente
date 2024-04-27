package api

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/santhozkumar/my-ente/ente"
)

type HealthCheckHandler struct {
	DB *sql.DB
}

func (h *HealthCheckHandler) Ping(c *gin.Context) {
	res := 0
	err := h.DB.QueryRowContext(c, `SELECT 1`).Scan(&res)
	// t := time.NewTimer(8 * time.Second)
	// <-t.C
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ente.NewInternalError("Ping failed"))
		return
	}
	result := make(map[string]string)
	result["id"] = os.Getenv("GIT_COMMIT")
	result["message"] = "pong"
	c.JSON(http.StatusOK, result)
}

func (h *HealthCheckHandler) PintDBStats(c *gin.Context) {
	_ = h.DB.Stats()
	stats := h.DB.Stats()
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// logger.Handler().WithAttrs(attrs []slog.Attr)
	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"DB Status",
		slog.Int("MaxOpenConnections", stats.MaxOpenConnections),
		slog.Int("OpenConnections", stats.OpenConnections),
		slog.Int("InUse", stats.InUse),
		slog.Int("Idle", stats.Idle),
		slog.Int("WaitCount", int(stats.WaitCount)),
		slog.String("WaitDuration", stats.WaitDuration.String()),
		slog.Int("MaxIdleClosed", int(stats.MaxIdleClosed)),
		slog.Int("MaxIdleTimeClosed", int(stats.MaxIdleTimeClosed)),
		slog.Int("MaxLifetimeClosed", int(stats.MaxLifetimeClosed)),
	)
	slog.Info("DB ping start")
	err := h.DB.Ping()
	slog.Info("DB ping end")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ente.NewInternalError(""))
		return
	}
	c.Status(http.StatusOK)

}
