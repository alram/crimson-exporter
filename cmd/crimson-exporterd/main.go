package main

import (
	"alram/crimson-exporter/crimson"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	e := crimson.NewCrimsonExporter()
	prometheus.MustRegister(e)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9092", nil)

}
