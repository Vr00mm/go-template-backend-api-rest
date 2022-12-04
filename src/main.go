package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"application-template/logger"
	"application-template/metrics"

	"github.com/heptiolabs/healthcheck"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nbari/violetear"
	"github.com/nbari/violetear/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	db        *sqlx.DB
	once      sync.Once
	metricsMw middleware.Chain
)

/*
export PGUSER=postgres
export PGDB=postgres
export PGHOST=127.0.0.1
export PGPORT=5432
export PGPASS=mynewpassword
*/
func MustGetConnection() *sqlx.DB {
	once.Do(func() {
		pguser := os.Getenv("PGUSER")
		pgdb := os.Getenv("PGDB")
		pghost := os.Getenv("PGHOST")
		pgport := os.Getenv("PGPORT")
		pgpass := os.Getenv("PGPASS")
		dbURI := fmt.Sprintf("user=%s dbname=%s host=%s port=%v sslmode=disable", pguser, pgdb, pghost, pgport)
		if pgpass != "" {
			dbURI += " password=" + pgpass
		}
		var err error
		db, err = sqlx.Connect("postgres", dbURI)
		if err != nil {
			panic(fmt.Sprintf("Unable to connection to database: %v\n", err))
		}
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(10)
	})
	return db
}

func main() {

	router := violetear.New()
	router.LogRequests = true
	router.RequestID = "REQUEST_LOG_ID"
	router.Logger = logger.JSONLogger
	router.Verbose = true

	conn := MustGetConnection()
	err := conn.Ping()
	if err != nil {
		log.Fatalln("Cannot Ping DATABASE: ", err)
	}

	healthz := healthcheck.NewHandler()
	healthz.AddLivenessCheck("database-"+os.Getenv("PGDB"), func() error {
		err := db.Ping()
		if err != nil {
			return fmt.Errorf("Cannot connect to database")
		}
		return nil
	})
	healthz.AddReadinessCheck("database-limit-connection", func() error {
		stats := db.Stats().OpenConnections
		if stats >= 10 {
			return fmt.Errorf("Too much open connection to DB")
		}
		return nil
	})

	counter := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "myAPI",
			Name:      "requests_total",
			Help:      "Total number of requests.",
		}, []string{"endpoint", "method"})

	r := prometheus.NewRegistry()
	r.MustRegister(counter)

	metricsMw = middleware.New(metrics.MetricsMW(counter), metrics.MetricsMWCallBack)

	router.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}), "GET").Name("/mectrics")
	router.Handle("/live", metricsMw.ThenFunc(healthz.LiveEndpoint), "GET").Name("/live")
	router.Handle("/ready", metricsMw.ThenFunc(healthz.ReadyEndpoint), "GET").Name("/ready")

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   7 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(srv.ListenAndServe())

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
