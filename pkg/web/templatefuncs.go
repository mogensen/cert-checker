package web

import (
	"html/template"
	"time"
)

func getFunctions() template.FuncMap {
	return template.FuncMap{
		"timeNow":         templateTimeNow,
		"tlsToCSS":        templateTLSToCSS,
		"expireTimeToCSS": templateExpireTimeToCSS,
	}
}

func templateTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func templateTLSToCSS(tlsStr string) string {
	switch tlsStr {
	case "SSLv3":
		return "failed"
	case "TLS 1.0":
		return "failed"
	case "TLS 1.1":
		return "warning"
	case "TLS 1.2":
		return "success"
	case "TLS 1.3":
		return "success"
	}
	return "failed"
}

func parseTime(s string) *time.Time {
	parsedTime, err := time.Parse("2006-01-02", s)
	if err != nil || (parsedTime == time.Time{}) {
		return nil
	}
	return &parsedTime
}

func templateExpireTimeToCSS(ts string) string {
	parsedTime := parseTime(ts)
	if parsedTime == nil {
		return "hidden"
	}
	if parsedTime.Before(time.Now()) {
		return "failed"
	}
	if parsedTime.Before(time.Now().AddDate(0, 0, minExpireDays)) {
		return "warning"
	}
	return "success"
}
