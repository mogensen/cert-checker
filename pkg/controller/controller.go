package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/genkiroid/cert"
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
func New(interval time.Duration, servingAddress string, log *logrus.Entry, certs []models.Certificate) *Controller {
	metrics := metrics.New(log)
	if err := metrics.Run(servingAddress); err != nil {
		log.Errorf("failed to start metrics server: %s", err)
		return nil
	}
	return &Controller{
		certs:    certs,
		metrics:  metrics,
		interval: interval,
		log:      log,
	}
}

// Certs exposes certificate info to external services
func (c *Controller) Certs() []models.Certificate {
	return c.certs
}

// Run starts the main loop that will call ProbeAll regularly.
func (c *Controller) Run(ctx context.Context) error {
	// Start by probing all certificates before starting the ticker
	c.probeAll(ctx)

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
			c.probeAll(ctx)
		}
	}
}

// probeAll triggers the Probe function for each registered service in the manager.
// Everything is done asynchronously.
func (c *Controller) probeAll(ctx context.Context) {
	c.log.Debug("Probing all")

	for id, cer := range c.certs {
		if ctx.Err() != nil {
			return
		}
		c.log.Debugf("Probing: %s", cer.DNS)

		cer.Info = cert.NewCert(cer.DNS)
		// For now we will ignore dial up errors
		if strings.HasPrefix(cer.Info.Error, "dial tcp") {
			c.log.Warnf("Problem checking %s : %s", cer.DNS, cer.Info.Error)
			continue
		}

		c.certs[id] = cer

		isValid := cer.Info.Error == ""

		if !isValid {
			c.log.Debugf(" - Found error for %s : %s", cer.DNS, cer.Info.Error)
		}
		c.metrics.AddCertificateInfo(cer, isValid)
	}
}

// Shutdown closes the metrics server gracefully
func (c *Controller) Shutdown() error {
	// If metrics server is not started than exit early
	if c.metrics == nil {
		return nil
	}

	c.log.Info("shutting down metrics server...")

	if err := c.metrics.Shutdown(); err != nil {
		return fmt.Errorf("metrics server shutdown failed: %s", err)
	}

	c.log.Info("metrics server gracefully stopped")

	return nil
}
