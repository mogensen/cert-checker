package metrics

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics exposes certificate checks as prometheus metrics
type Metrics struct {
	*http.Server

	registry        *prometheus.Registry
	dnsCertValidity *prometheus.GaugeVec
	log             *logrus.Entry

	// container cache stores a cache of a container's current image, version,
	// and the latest
	containerCache map[string]cacheItem
	mu             sync.Mutex
}

type cacheItem struct {
	issuer    string
	notBefore string
	notAfter  string
}

// New returns a new configured instance of the Metrics server
func New(log *logrus.Entry) *Metrics {
	containerImageVersion := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cert_checker",
			Name:      "is_valid",
			Help:      "Detailing if the certificate served by the server at the dns is valid",
		},
		[]string{
			"dns", "issuer", "not_before", "not_after",
		},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(containerImageVersion)

	return &Metrics{
		log:             log,
		registry:        registry,
		dnsCertValidity: containerImageVersion,
		containerCache:  make(map[string]cacheItem),
	}
}

// Run will run the metrics server
func (m *Metrics) Run(servingAddress string) error {
	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))

	ln, err := net.Listen("tcp", servingAddress)
	if err != nil {
		return err
	}

	m.Server = &http.Server{
		Addr:           ln.Addr().String(),
		ReadTimeout:    8 * time.Second,
		WriteTimeout:   8 * time.Second,
		MaxHeaderBytes: 1 << 15, // 1 MiB
		Handler:        router,
	}

	go func() {
		m.log.Infof("serving metrics on %s/metrics", ln.Addr())

		if err := m.Serve(ln); err != nil {
			m.log.Errorf("failed to serve prometheus metrics: %s", err)
			return
		}
	}()

	return nil
}

// AddCertificateInfo registers a new or updates and existing certificate record
func (m *Metrics) AddCertificateInfo(dns, issuer, notAfter, notBefore string, isValid bool) {
	// Remove old certificate information if it exists
	m.RemoveCertificateInfo(dns)

	m.mu.Lock()
	defer m.mu.Unlock()

	isValidF := 0.0
	if isValid {
		isValidF = 1.0
	}

	m.dnsCertValidity.With(
		m.buildLabels(dns, issuer, notAfter, notBefore),
	).Set(isValidF)

	m.containerCache[dns] = cacheItem{
		issuer:    issuer,
		notAfter:  notAfter,
		notBefore: notBefore,
	}
}

// RemoveCertificateInfo removed an existing certificate record
func (m *Metrics) RemoveCertificateInfo(dns string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.containerCache[dns]
	if !ok {
		return
	}

	m.dnsCertValidity.Delete(
		m.buildLabels(dns, item.issuer, item.notBefore, item.notAfter),
	)
	delete(m.containerCache, dns)
}

func (m *Metrics) buildLabels(dns, issuer, notBefore, notAfter string) prometheus.Labels {
	return prometheus.Labels{
		"dns":        dns,
		"issuer":     issuer,
		"not_before": notBefore,
		"not_after":  notAfter,
	}
}

// Shutdown closes the metrics server gracefully
func (m *Metrics) Shutdown() error {
	// If metrics server is not started than exit early
	if m.Server == nil {
		return nil
	}

	m.log.Info("shutting down prometheus metrics server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := m.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("prometheus metrics server shutdown failed: %s", err)
	}

	m.log.Info("prometheus metrics server gracefully stopped")

	return nil
}
