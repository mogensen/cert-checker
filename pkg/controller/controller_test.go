package controller

import (
	"context"
	"testing"
	"time"

	"github.com/mogensen/cert-checker/pkg/metrics"
	"github.com/mogensen/cert-checker/pkg/models"
	"github.com/sirupsen/logrus"
)

func TestController_Run_StopsWhenContextIsCanceled(t *testing.T) {
	type fields struct {
		log      *logrus.Entry
		metrics  *metrics.Metrics
		certs    []models.Certificate
		interval time.Duration
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		interval time.Duration
	}{
		{interval: time.Millisecond * 100},
		{interval: time.Hour * 1},
		{interval: time.Second * 2},
	}
	for _, tt := range tests {
		t.Run("TestController_Run_StopsWhenContextIsCanceled "+tt.interval.String(), func(t *testing.T) {

			log := logrus.NewEntry(logrus.New())
			c := New(tt.interval, metrics.New(log), log, []models.Certificate{})

			timeout := time.After(2 * time.Second)
			done := make(chan bool)
			go func() {
				ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
				if err := c.Run(ctx); err != nil {
					t.Errorf("Controller.Run() error = %v", err)
				}
				done <- true
			}()

			select {
			case <-timeout:
				t.Fatal("Test didn't finish in time")
			case <-done:
			}

		})
	}
}
