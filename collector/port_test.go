package collector

import (
	"exporter/metrics/port"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"testing"
)

func TestNewPortCollector(t *testing.T) {
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
		fmt.Printf("%s , %v  , %v \n",prometheus.BuildFQName(HostName, " ", exec),fmt.Sprintf("%s. state %s ", "tcp "+saddr, e.State), e.Process.PidValue())
	}
}
func TestCollectorPorts(t *testing.T) {
	tabs, err := port.TCPSocks(func(s *port.SockTabEntry) bool {
		return s.State == port.Listen
	})
	if err == nil {
		displaySockInfo("tcp", tabs)
	}
}


func displaySockInfo(proto string, s []port.SockTabEntry) {
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

	for _, e := range s {
		p := ""
		if e.Process != nil {
			p = e.Process.String()
		}
		saddr := lookup(e.LocalAddr)
		fmt.Printf("%-5s %-23.23s  %-12s %-16s\n", proto, saddr,  e.State, p)
	}
}