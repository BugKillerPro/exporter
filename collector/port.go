package collector

import (
	"exporter/metrics/port"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"os"
)

var (
	HostName string
)

func init() {
	HostName, _ = os.Hostname()
}

type PortCollector struct {
	portDesc *prometheus.Desc
}

func NewPortCollector() PortCollector {
	return PortCollector{portDesc: prometheus.NewDesc(
		"port_listening_state",
		"help",
		nil,
		nil)}

}

// Describe implements the prometheus.Collector interface.
func (p PortCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.portDesc
}

// Collect implements the prometheus.Collector interface.This method may be called concurrently and must therefore be
// implemented in a concurrency safe way
func (p PortCollector) Collect(ch chan<- prometheus.Metric) {
	tabs, err := port.TCPSocks(func(s *port.SockTabEntry) bool {
		return s.State == port.Listen
	})
	if err != nil {
		return
	}
	lookup := func(skaddr *port.SockAddr) string {
		const IPv4Strlen = 17
		addr := skaddr.IP.String()
		names, err := net.LookupAddr(addr)
		if err == nil && len(names) > 0 {
			addr = names[0]
		}
		if len(addr) > IPv4Strlen {
			addr = addr[:IPv4Strlen]
		}
		return fmt.Sprintf("%s:%d", addr, skaddr.Port)
	}

	for _, e := range tabs {
		exec := ""
		if e.Process != nil {
			exec = e.Process.ExecName()
		}
		saddr := lookup(e.LocalAddr)
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				"custom_port_exporter",
				"custom_port_exporter scrape node port state info",
				[]string{"host","app","port","state","pid"},
				nil,
			),
			prometheus.UntypedValue,
			e.Process.PidFloatValue(),
			[]string{ HostName, exec, "tcp "+saddr, e.State.String(), e.Process.PidValue()}...,
		)
	}
}
