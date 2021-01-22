package models

import "github.com/mogensen/cert"

// Certificate stores the dns name, the information about the cert, and possible errors
type Certificate struct {
	DNS  string     `json:"dns"`
	Info *cert.Cert `json:"info"`
}
