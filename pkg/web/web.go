package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mogensen/cert-checker/pkg/models"
	"github.com/sirupsen/logrus"
)

// UI exposes an html endpoint for certificates
type UI struct {
	log *logrus.Entry

	certService certProvider
	webAddress  string
	server      *http.Server
}

type certProvider interface {
	Certs() []models.Certificate
}

// New returns a new configured instance of the UI struct
func New(certService certProvider, webAddress string, log *logrus.Entry) *UI {
	return &UI{
		webAddress:  webAddress,
		certService: certService,
		log:         log,
	}
}

// Run will run the ui server
func (u *UI) Run(ctx context.Context) error {
	router := http.NewServeMux()
	router.Handle("/", u.handleFunc())

	u.server = &http.Server{
		Addr:         u.webAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		u.log.Infof("serving ui on %s", u.server.Addr)
		if err := u.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			u.log.Fatalf("Could not listen on %s: %v\n", u.webAddress, err)
		}
	}()

	return nil
}

// Shutdown closes the ui server gracefully
func (u *UI) Shutdown() error {
	// If ui server is not started than exit early
	if u.server == nil {
		return nil
	}

	u.log.Info("shutting down ui server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := u.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("ui server shutdown failed: %s", err)
	}

	u.log.Info("ui server gracefully stopped")

	return nil
}

func (u *UI) handleFunc() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		certs := u.certService.Certs()

		err := templateHTML(certs, w)
		if err != nil {
			logrus.Printf("Error templating: %v\n", err)
		}
	})
}
