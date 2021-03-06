package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// "udp://host:1234", 80 => udp host:1234 host 1234
// "host:1234", 80       => tcp host:1234 host 1234
// "host", 80            => tcp host:80   host 80
func parseAddr(addr string, defaultPort int) (network, address string, err error) {
	u, err := url.Parse(strings.ToLower(addr))
	if err != nil {
		// It's possible that we're dealing with IPv6 localhost here!
		if strings.HasPrefix(addr, "[::]") {
			_, port, e := net.SplitHostPort(addr)
			if e != nil {
				return network, address, e
			}
			if port == "" {
				port = strconv.Itoa(defaultPort)
			}
			return "tcp", fmt.Sprintf("0.0.0.0:%s", port), nil
		}

		return network, address, err
	}

	switch {
	case u.Scheme == "" && u.Opaque == "" && u.Host == "" && u.Path != "": // "host"
		u.Scheme, u.Opaque, u.Host, u.Path = "tcp", "", net.JoinHostPort(u.Path, strconv.Itoa(defaultPort)), ""
	case u.Scheme != "" && u.Opaque != "" && u.Host == "" && u.Path == "": // "host:1234"
		u.Scheme, u.Opaque, u.Host, u.Path = "tcp", "", net.JoinHostPort(u.Scheme, u.Opaque), ""
	case u.Scheme != "" && u.Opaque == "" && u.Host != "" && u.Path == "": // "tcp://host[:1234]"
		if _, _, err := net.SplitHostPort(u.Host); err != nil {
			u.Host = net.JoinHostPort(u.Host, strconv.Itoa(defaultPort))
		}
	default:
		return network, address, errors.Errorf("%s: unsupported address format", addr)
	}

	return u.Scheme, u.Host, nil
}

func registerMetrics(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}

func registerProfile(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}

func interrupt(cancel <-chan struct{}, dst chan<- struct{}) error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		select {
		case <-c:
			dst <- struct{}{}
		case <-cancel:
			return errors.New("canceled")
		}
	}
}
