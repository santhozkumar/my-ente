package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoKillers/libsodium-go/sodium"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/santhozkumar/my-ente/pkg/api"
	"github.com/santhozkumar/my-ente/pkg/utils/config"
	// log "github.com/sirupsen/logrus"
)

func main() {
	// setLogger("local")
	err := config.ConfigureViper("local")
	if err != nil {
		panic(err)
	}

	db := setupDatabase()
	defer db.Close()

	sodium.Init()

	os.Hostname()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!\n")
	})

	healthCheckHandler := api.HealthCheckHandler{DB: db}
	// r.GET("/ping", healthCheckHandler.Ping)
	r.GET("/ping", timeout.New(
		timeout.WithTimeout(5*time.Second),
		timeout.WithHandler(healthCheckHandler.Ping),
		timeout.WithResponse(timeoutResponse)))

	r.GET("/ping-db", timeout.New(
		timeout.WithTimeout(5*time.Second),
		timeout.WithHandler(healthCheckHandler.PintDBStats),
		timeout.WithResponse(timeoutResponse)))
	r.Run()

}

func timeoutResponse(c *gin.Context) {
    c.JSON(http.StatusRequestTimeout, gin.H{"handler":true})
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
