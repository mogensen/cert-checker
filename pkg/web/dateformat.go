package web

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	defaultDateFormat = "YYYY-MM-DD"
)

var unsupportedDateTokens = regexp.MustCompile(`[A-Za-z]`)

type formatter struct {
	ConfigLayout string
	GoLayout     string
}

func defaultFormatter() formatter {
	return formatter{
		ConfigLayout: defaultDateFormat,
		GoLayout:     "2006-01-02",
	}
}

func newFormatter(layout string) (formatter, error) {
	if layout == "" {
		return defaultFormatter(), nil
	}

	goLayout := strings.NewReplacer(
		"yyyy", "2006",
		"YYYY", "2006",
		"YY", "06",
		"yy", "06",
		"DD", "02",
		"MM", "01",
		"dd", "02",
		"HH", "15",
		"mm", "04",
		"ss", "05",
	).Replace(layout)

	if unsupportedDateTokens.MatchString(goLayout) {
		return formatter{}, fmt.Errorf("unsupported date format token in %q", layout)
	}

	if _, err := time.Parse(goLayout, time.Now().Format(goLayout)); err != nil {
		return formatter{}, fmt.Errorf("invalid date format %q: %w", layout, err)
	}

	return formatter{
		ConfigLayout: layout,
		GoLayout:     goLayout,
	}, nil
}

func (f formatter) Format(t time.Time) string {
	return t.Format(f.GoLayout)
}
