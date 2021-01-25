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

			nlog := logrus.New()
			nlog.SetOutput(os.Stdout)
			nlog.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

			log := logrus.NewEntry(nlog)

			configFile, err := cmd.Flags().GetString("config")
			if err != nil {
				nlog.Fatalf("Could not ger configuration file: %v", err)
			}
			opts, err := newOptionsFromFile(configFile)
			if err != nil {
				nlog.Fatalf("Could not start server: %v", err)
			}

			logLevel, err := logrus.ParseLevel(opts.LogLevel)
			if err != nil {
				return fmt.Errorf("failed to parse  loglevel %q: %s",
					opts.LogLevel, err)
			}
			nlog.SetLevel(logLevel)

			metrics := metrics.New(log)
			if err := metrics.Run("0.0.0.0:8080"); err != nil {
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
