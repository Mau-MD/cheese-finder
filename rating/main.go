package main

import (
	"cheesefinder/db"
	"cheesefinder/rate"
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func main() {

	routerLogger := httplog.NewLogger("rating-service", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   false,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"version": "v1",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})

	slog.SetDefault(routerLogger.Logger)

	db.Init()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(routerLogger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Mount("/v1", routerV1())

	http.ListenAndServe(":3000", r)
}

func routerV1() http.Handler {
	r := chi.NewRouter()
	r.Route("/rate", func(r chi.Router) {
		r.Get("/{limit}", rate.GetObjectsToRate)
		r.Post("/", rate.RateObject)
	})
	return r
}
