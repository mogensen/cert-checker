package app

import (
	"context"
	"fmt"
	"os"

	"github.com/mogensen/cert-checker/pkg/controller"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mogensen/cert-checker/pkg/metrics"
)

const (
	helpOutput = "Certificate monitoring utility for watching tls certificates and reporting the result as metrics."
)

// NewCommand sets up the version-checker command and all dependencies
func NewCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version-checker",
		Short: helpOutput,
		Long:  helpOutput,
		RunE: func(cmd *cobra.Command, args []string) error {

			opts, err := newOptionsFromFile("config.yaml")
			if err != nil {
				panic(fmt.Sprintf("Could not start server: %v", err))
			}

			logLevel, err := logrus.ParseLevel(opts.LogLevel)
			if err != nil {
				return fmt.Errorf("failed to parse  loglevel %q: %s",
					opts.LogLevel, err)
			}

			nlog := logrus.New()
			nlog.SetOutput(os.Stdout)
			nlog.SetLevel(logLevel)
			nlog.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

			log := logrus.NewEntry(nlog)

			metrics := metrics.New(log)
			if err := metrics.Run("0.0.0.0:9090"); err != nil {
				return fmt.Errorf("failed to start metrics server: %s", err)
			}

			defer func() {
				if err := metrics.Shutdown(); err != nil {
					log.Error(err)
				}
			}()

			c := controller.New(opts.IntervalDuration, metrics, log, opts.Certificates)

			return c.Run(ctx)
		},
	}

	return cmd
}
