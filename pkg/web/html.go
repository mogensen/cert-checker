package web

import (
	"embed"
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/mogensen/cert-checker/pkg/models"
)

//go:embed views/*
var views embed.FS

// Expire in 30 days is warning
// TODO make configurable
var minExpireDays = 30

// templateHTML generates an html representation for the given certs, and writes the result to the io.writer
func templateHTML(certs []models.Certificate, w io.Writer) error {

	sum := internalSummery{}

	for _, c := range certs {
		if c.Info == nil {
			// Cert is not proccessed yet
			continue
		}
		uiC := uiCert{
			DNS:               c.DNS,
			Issuer:            c.Info.Issuer,
			NotAfter:          c.Info.Detail().NotAfter.Format("2006-01-02"),
			NotBefore:         c.Info.Detail().NotBefore.Format("2006-01-02"),
			MinimumTLSVersion: c.Info.MinimumTLSVersion,
			Warning:           warning(c),
			Error:             c.Info.Error,
		}

		if uiC.Error != "" {
			sum.BadCerts = append(sum.BadCerts, uiC)
		} else if uiC.Warning != "" {
			sum.WarningCerts = append(sum.WarningCerts, uiC)
		} else {
			sum.GoodCerts = append(sum.GoodCerts, uiC)

		}
	}

	t := template.Must(template.New("index.html").Funcs(getFunctions()).ParseFS(views, "views/*"))
	err := t.Execute(w, sum)
	if err != nil {
		return err
	}
	return nil
}

func warning(cert models.Certificate) string {
	warnings := []string{}
	switch cert.Info.MinimumTLSVersion {
	case "SSLv3":
		fallthrough
	case "TLS 1.0":
		fallthrough
	case "TLS 1.1":
		warnings = append(warnings, "TLS version deprecated")
	}
	if cert.Info.Detail().NotAfter.Before(time.Now().AddDate(0, 0, minExpireDays)) {
		warnings = append(warnings, "Certificate is about to expire")
	}
	return strings.Join(warnings, ", ")
}

type uiCert struct {
	DNS               string `json:"dns"`
	Issuer            string `json:"issuer"`
	NotBefore         string `json:"notBefore"`
	NotAfter          string `json:"notAfter"`
	MinimumTLSVersion string `json:"minimumTLSVersion"`
	Warning           string `json:"warning"`
	Error             string `json:"error"`
}

type internalSummery struct {
	GoodCerts    []uiCert
	WarningCerts []uiCert
	BadCerts     []uiCert
}
