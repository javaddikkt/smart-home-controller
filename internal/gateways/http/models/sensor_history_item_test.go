package models

import (
	"github.com/go-openapi/strfmt"
	"testing"
	"time"
)

func TestSensorHistoryItem_Validate(t *testing.T) {
	now := strfmt.DateTime(time.Now())
	valid := &SensorHistoryItem{
		Payload:   42,
		Timestamp: now,
	}
	tests := []struct {
		name    string
		mutate  func(h *SensorHistoryItem)
		wantErr bool
	}{
		{"valid history item", func(_ *SensorHistoryItem) {}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := *valid
			tt.mutate(&h)
			err := h.Validate(strfmt.Default)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
