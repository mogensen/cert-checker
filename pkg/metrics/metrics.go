package metrics

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/mogensen/cert-checker/pkg/models"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics exposes certificate checks as prometheus metrics
type Metrics struct {
	*http.Server

	registry       *prometheus.Registry
	certExpiration *prometheus.GaugeVec
	certValidity   *prometheus.GaugeVec
	log            *logrus.Entry

	// container cache stores a cache of a container's current image, version,
	// and the latest
	containerCache map[string]models.Certificate
	mu             sync.Mutex
}

// New returns a new configured instance of the Metrics server
func New(log *logrus.Entry) *Metrics {
	certValidity := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cert_checker",
			Name:      "is_valid",
			Help:      "Detailing if the certificate served by the server at the dns is valid",
		},
		[]string{
			"dns", "issuer", "not_before", "not_after", "cert_error",
		},
	)

	certExpiration := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cert_checker",
			Name:      "expire_time",
			Help:      "Detailing when a certificate is set to expire",
		},
		[]string{
			"dns", "issuer", "not_before", "not_after",
		},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(certValidity)
	registry.MustRegister(certExpiration)

	return &Metrics{
		log:            log,
		registry:       registry,
		certExpiration: certExpiration,
		certValidity:   certValidity,
		containerCache: make(map[string]models.Certificate),
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
func (m *Metrics) AddCertificateInfo(cer models.Certificate, isValid bool) {
	// Remove old certificate information if it exists
	m.RemoveCertificateInfo(cer.DNS)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.containerCache[cer.DNS] = cer

	isValidF := 0.0
	if isValid {
		isValidF = 1.0
	}

	m.certValidity.With(
		m.buildLabelsValidity(cer),
	).Set(isValidF)

	if !isValid {
		return
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", cer.Info.NotAfter)
	if err != nil {
		fmt.Println(err)
		return
	}

	m.certExpiration.With(
		m.buildLabelsExpiration(cer),
	).Set(float64(parsedTime.Unix()))
}

// RemoveCertificateInfo removed an existing certificate record
func (m *Metrics) RemoveCertificateInfo(dns string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.containerCache[dns]
	if !ok {
		m.log.Debugf("Did not find %s in cache", dns)
		return
	}

	m.certValidity.Delete(m.buildLabelsValidity(item))
	m.certExpiration.Delete(m.buildLabelsExpiration(item))

	delete(m.containerCache, dns)
}

func (m *Metrics) buildLabelsExpiration(cer models.Certificate) prometheus.Labels {
	return prometheus.Labels{
		"dns":        cer.DNS,
		"issuer":     cer.Info.Issuer,
		"not_before": cer.Info.NotBefore,
		"not_after":  cer.Info.NotAfter,
	}
}

func (m *Metrics) buildLabelsValidity(cer models.Certificate) prometheus.Labels {
	return prometheus.Labels{
		"dns":        cer.DNS,
		"issuer":     cer.Info.Issuer,
		"not_before": cer.Info.NotBefore,
		"not_after":  cer.Info.NotAfter,
		"cert_error": cer.Info.Error,
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
