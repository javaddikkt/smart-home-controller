package models

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func TestSensorEvent_Validate(t *testing.T) {
	valid := &SensorEvent{
		Payload:            swag.Int64(10),
		SensorSerialNumber: swag.String("1234567890"),
	}
	tests := []struct {
		name    string
		mutate  func(e *SensorEvent)
		wantErr bool
	}{
		{"valid event", func(_ *SensorEvent) {}, false},
		{"missing payload", func(e *SensorEvent) { e.Payload = nil }, true},
		{"missing serial number", func(e *SensorEvent) { e.SensorSerialNumber = nil }, true},
		{"invalid serial number format", func(e *SensorEvent) { e.SensorSerialNumber = swag.String("abcd") }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := *valid
			tt.mutate(&e)
			err := e.Validate(strfmt.Default)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
