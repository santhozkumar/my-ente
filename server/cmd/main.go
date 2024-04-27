package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GoKillers/libsodium-go/sodium"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/santhozkumar/my-ente/pkg/api"
	"github.com/santhozkumar/my-ente/pkg/middleware"
	"github.com/santhozkumar/my-ente/pkg/utils/config"

	"github.com/patrickmn/go-cache"

	ginprometheus "github.com/zsais/go-gin-prometheus"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func setLoggerSlog(environment string) {
	var handler slog.Handler
	if environment == "local" || environment == "" {
		handler = slog.NewTextHandler(os.Stdout, nil)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	// setLogger("local")
	setLoggerSlog("local")
	err := config.ConfigureViper("local")
	if err != nil {
		panic(err)
	}

	db := setupDatabase()
	defer db.Close()

	sodium.Init()

	os.Hostname()

	server := gin.New()
	server.HandleMethodNotAllowed = true

	authCache := cache.New(1*time.Minute, 15*time.Minute)
	authMiddleware := m
	rateLimiter := middleware.NewRateLimitMiddlware()

	p := ginprometheus.NewPrometheus("museum")
	p.Use(server)

	server.Use(requestid.New())
	server.Use(requestid.New(), middleware.Logger(urlSanitizer))

	publicAPI := server.Group("/")
	publicAPI.Use(rateLimiter.APIRateLimitMiddleWare(urlSanitizer))

	privateAPI := server.Group("/")
	privateAPI.Use()

	adminAPI := server.Group("/admin")
	adminAPI.Use()

	server.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!\n")
	})

	healthCheckHandler := api.HealthCheckHandler{DB: db}
	publicAPI.GET("/ping", timeout.New(
		timeout.WithTimeout(5*time.Second),
		timeout.WithHandler(healthCheckHandler.Ping),
		timeout.WithResponse(timeoutResponse)))

	publicAPI.GET("/ping-db", timeout.New(
		timeout.WithTimeout(5*time.Second),
		timeout.WithHandler(healthCheckHandler.PintDBStats),
		timeout.WithResponse(timeoutResponse)))

	// publicAPI.POST("/method-check", MethodCheck)
	// publicAPI.DELETE("/method-check", MethodCheckDelete)

	publicAPI.
		server.NoMethod(MethodNotAllowedMiddleWare(server.Routes()))

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	server.Run()
}

func timeoutResponse(c *gin.Context) {
	c.JSON(http.StatusRequestTimeout, gin.H{"handler": true})
}

func setupDatabase() *sql.DB {
	log.Println("Setting up db")
	log.Println(config.GetPGInfo())

	db, err := sql.Open("postgres", config.GetPGInfo())
	if err != nil {
		log.Panic(err)
		panic(err)
	}
	log.Println("Connected to DB")
	err = db.Ping()
	if err != nil {
		log.Panic(err)
		panic(err)
	}
	log.Println("Pinged DB")
	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	log.Println(driver)

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", "postgres", driver)

	log.Println(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Panic(err)
		panic(err)

	}
	db.SetMaxIdleConns(6)
	db.SetMaxOpenConns(30)

	log.Println("database configured successfully")
	return db
}

// func setLogger(environment string) {
// 	log.SetReportCaller(true)
// 	callerPrettyfier := func(f *runtime.Frame) (string, string) {
// 		s := strings.Split(f.Function, ".")
// 		funcName := s[len(s)-1]
// 		return funcName, fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
// 	}
//
// 	if environment == "local" {
// 		log.SetFormatter(&log.TextFormatter{
// 			CallerPrettyfier: callerPrettyfier,
// 			DisableQuote:     true,
// 			ForceColors:      true,
// 		})
// 	}
// }

func urlSanitizer(c *gin.Context) string {
	if c.Request.Method == http.MethodOptions {
		return "/options"
	}

	u := c.Request.URL
	u.RawQuery = ""
	uri := u.RequestURI()

	for _, p := range c.Params {
		uri = strings.Replace(uri, p.Value, fmt.Sprintf(":%s", p.Key), 1)
	}
	return c.Request.URL.Path
}

func MethodCheck(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}
	result := make(map[string]string)
	result["id"] = os.Getenv("GIT_COMMIT")
	result["message"] = "pong"
	c.JSON(http.StatusOK, result)
}

func MethodCheckDelete(c *gin.Context) {
	if c.Request.Method != http.MethodDelete {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}
	result := make(map[string]string)
	result["id"] = os.Getenv("GIT_COMMIT")
	result["message"] = "pong"
	c.JSON(http.StatusOK, result)
}

// func isAvailableMethods(c.r) {
// }
func MethodNotAllowedMiddleWare(routes gin.RoutesInfo) gin.HandlerFunc {
	log.Print("routes", routes)
	return func(c *gin.Context) {
		allowedMethods := []string{}
		for _, route := range routes {
			if route.Path == c.Request.URL.Path {
				allowedMethods = append(allowedMethods, route.Method)
			}
		}
		if len(allowedMethods) > 0 {
			c.Writer.Header().Add("Allow", strings.Join(allowedMethods, ", "))
			// for _, allowedMethod := range allowedMethods {
			//     c.Writer.Header().Add("Allow", allowedMethod)
			// }
		}
	}
}
