package web

import (
	"testing"
	"time"
)

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		wantLayout string
		wantErr    bool
	}{
		{
			name:       "default empty",
			in:         "",
			wantLayout: "2006-01-02",
		},
		{
			name:       "default explicit",
			in:         "YYYY-MM-DD",
			wantLayout: "2006-01-02",
		},
		{
			name:       "dots with day first",
			in:         "dd.MM.YYYY",
			wantLayout: "02.01.2006",
		},
		{
			name:       "with clock",
			in:         "YYYY-MM-DD HH:mm:ss",
			wantLayout: "2006-01-02 15:04:05",
		},
		{
			name:    "unsupported token",
			in:      "QQ.MM.YYYY",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := newFormatter(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q", tt.in)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if f.GoLayout != tt.wantLayout {
				t.Fatalf("expected go layout %q, got %q", tt.wantLayout, f.GoLayout)
			}
		})
	}
}

func TestFormatterFormat(t *testing.T) {
	f, err := newFormatter("dd.MM.YYYY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	in := time.Date(2026, 4, 15, 14, 30, 0, 0, time.UTC)
	got := f.Format(in)
	if got != "15.04.2026" {
		t.Fatalf("expected 15.04.2026, got %q", got)
	}
}
