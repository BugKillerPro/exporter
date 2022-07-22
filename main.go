package main

import (
	"exporter/collector"
	"fmt"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	prom.MustRegister(collector.NewPortCollector())
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8088", nil); err != nil {
		fmt.Printf("Error occur when start custom collector on %v %v",collector.HostName, err)
	}
}