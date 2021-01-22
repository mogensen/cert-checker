package controller

import (
	"context"
	"time"

	"github.com/mogensen/cert"
	"github.com/mogensen/cert-checker/pkg/metrics"
	"github.com/mogensen/cert-checker/pkg/models"
	"github.com/sirupsen/logrus"
)

// Controller probes certificates and registers the result in the metrics server
type Controller struct {
	log *logrus.Entry

	metrics  *metrics.Metrics
	certs    []models.Certificate
	interval time.Duration
}

// New returns a new configured instance of the Controller struct
func New(interval time.Duration, metrics *metrics.Metrics, log *logrus.Entry, certs []models.Certificate) *Controller {
	return &Controller{
		certs:    certs,
		metrics:  metrics,
		interval: interval,
		log:      log,
	}
}

// Run starts the main loop that will call ProbeAll regularly.
func (c *Controller) Run(ctx context.Context) error {
	// Start by probing all certificates before starting the ticker
	c.probeAll()

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		//select as usual
		select {
		case <-ctx.Done():
			c.log.Info("Stopping controller..")
			return nil
		case <-ticker.C:
			//give priority to a possible concurrent Done() event non-blocking way
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			c.probeAll()
		}
	}
}

// probeAll triggers the Probe function for each registered service in the manager.
// Everything is done asynchronously.
func (c *Controller) probeAll() {
	c.log.Debug("Probing all")

	for id, cer := range c.certs {
		c.log.Debugf("Probing: %s", cer.DNS)

		cer.Info = cert.NewCert(cer.DNS)
		c.certs[id] = cer

		isValid := cer.Info.Error == ""

		c.metrics.AddCertificateInfo(cer.DNS, cer.Info.Issuer, cer.Info.NotBefore, cer.Info.NotAfter, isValid)
	}
}
