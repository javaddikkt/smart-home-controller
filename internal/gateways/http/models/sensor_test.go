package models

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func TestSensor_Validate(t *testing.T) {
	now := strfmt.DateTime(time.Now())

	validSensor := &Sensor{
		ID:           swag.Int64(1),
		SerialNumber: swag.String("1234567890"),
		Type:         swag.String(SensorTypeCc),
		Description:  swag.String("Test sensor"),
		IsActive:     swag.Bool(true),
		CurrentState: swag.Int64(0),
		RegisteredAt: &now,
		LastActivity: &now,
	}

	tests := []struct {
		name    string
		mutate  func(s *Sensor)
		wantErr bool
	}{
		{"valid sensor", func(_ *Sensor) {}, false},
		{"missing serial number", func(s *Sensor) { s.SerialNumber = nil }, true},
		{"invalid serial number format", func(s *Sensor) { s.SerialNumber = swag.String("abcd") }, true},
		{"missing current state", func(s *Sensor) { s.CurrentState = nil }, true},
		{"missing activity state", func(s *Sensor) { s.IsActive = nil }, true},
		{"missing last activity", func(s *Sensor) { s.LastActivity = nil }, true},
		{"missing registered at", func(s *Sensor) { s.RegisteredAt = nil }, true},
		{"missing description", func(s *Sensor) { s.Description = nil }, true},
		{"invalid sensor type", func(s *Sensor) { s.Type = swag.String("invalid-type") }, true},
		{"missing ID", func(s *Sensor) { s.ID = nil }, true},
		{"ID less than minimum", func(s *Sensor) { s.ID = swag.Int64(0) }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sensor := *validSensor
			tt.mutate(&sensor)
			err := sensor.Validate(strfmt.Default)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
