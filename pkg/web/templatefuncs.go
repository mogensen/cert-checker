package web

import (
	"html/template"
	"time"
)

func getFunctions(f formatter) template.FuncMap {
	return template.FuncMap{
		"timeNow": func() string {
			return templateTimeNow(f)
		},
		"tlsToCSS": templateTLSToCSS,
		"expireTimeToCSS": func(ts string) string {
			return templateExpireTimeToCSS(ts, f)
		},
	}
}

func templateTimeNow(f formatter) string {
	return f.Format(time.Now())
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

func parseTime(s string, f formatter) *time.Time {
	parsedTime, err := time.Parse(f.GoLayout, s)
	if err != nil || (parsedTime == time.Time{}) {
		return nil
	}
	return &parsedTime
}

func templateExpireTimeToCSS(ts string, f formatter) string {
	parsedTime := parseTime(ts, f)
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
