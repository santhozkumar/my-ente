package healthcheck

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/santhozkumar/my-ente/ente"
	// "github.com/sirupsen/logrus"
)

type HealthCheckHandler struct {
	DB *sql.DB
}

func (h *HealthCheckHandler) Ping(c *gin.Context) {
	res := 0
	err := h.DB.QueryRowContext(c, `SELECT 1`).Scan(&res)
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
	// stats := h.DB.Stats()
	// logrus.WithFields(logrus.Fields{
	// 	"MaxOpenConnections": stats.MaxOpenConnections, // Maximum number of open connections to the database.
	// 	"OpenConnections":    stats.OpenConnections,    // The number of established connections both in use and idle.
	// 	"InUse":              stats.InUse,              // The number of connections currently in use.
	// 	"Idle":               stats.Idle,               // The number of idle connections.
	//
	// 	// Counters
	// 	"WaitCount":         stats.WaitCount,             // The total number of connections waited for.
	// 	"WaitDuration":      stats.WaitDuration.String(), // The total time blocked waiting for a new connection.
	// 	"MaxIdleClosed":     stats.MaxIdleClosed,         // The total number of connections closed due to SetMaxIdleConns.
	// 	"MaxIdleTimeClosed": stats.MaxIdleClosed,         // The total number of connections closed due to SetConnMaxIdleTime.
	// 	"MaxLifetimeClosed": stats.MaxLifetimeClosed,     // The total number of connections closed due to SetConnMaxLifetime.
	// }).Info("DB Status")
	// logrus.Info("DB ping start")
	err := h.DB.Ping()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ente.NewInternalError(""))
		return
	}
	c.Status(http.StatusOK)

}
