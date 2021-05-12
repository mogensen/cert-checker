package app

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/mogensen/cert-checker/pkg/controller"
	"github.com/mogensen/cert-checker/pkg/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

			// create a WaitGroup
			wg := new(sync.WaitGroup)
			wg.Add(2)

			// Metrics

			metricsAddress := fmt.Sprintf("%s:%d", "0.0.0.0", opts.Port)
			c := controller.New(opts.IntervalDuration, metricsAddress, log, opts.Certificates)

			go func() {
				<-ctx.Done()
				if err := c.Shutdown(); err != nil {
					log.Error(err)
				}
			}()

			go func() {
				c.Run(ctx)
				wg.Done()
			}()

			// Web UI

			webAddress := fmt.Sprintf("%s:%d", "0.0.0.0", opts.WebPort)
			ui := web.New(c, webAddress, log)

			go func() {
				<-ctx.Done()
				if err := ui.Shutdown(); err != nil {
					log.Error(err)
				}
			}()

			go func() {
				ui.Run(ctx)
				wg.Done()
			}()

			// wait until WaitGroup is done
			wg.Wait()
			log.Infof("Everything is successfully stopped")

			return nil
		},
	}

	return cmd
}
