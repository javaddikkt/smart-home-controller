package models

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"testing"
)

func TestSensorToCreate_Validate(t *testing.T) {
	valid := &SensorToCreate{
		Description:  swag.String("Test sensor"),
		IsActive:     swag.Bool(true),
		SerialNumber: swag.String("1234567890"),
		Type:         swag.String(SensorToCreateTypeCc),
	}
	tests := []struct {
		name    string
		mutate  func(s *SensorToCreate)
		wantErr bool
	}{
		{"valid sensor to create", func(_ *SensorToCreate) {}, false},
		{"missing description", func(s *SensorToCreate) { s.Description = nil }, true},
		{"missing is_active", func(s *SensorToCreate) { s.IsActive = nil }, true},
		{"missing serial number", func(s *SensorToCreate) { s.SerialNumber = nil }, true},
		{"invalid serial number format", func(s *SensorToCreate) { s.SerialNumber = swag.String("abcd") }, true},
		{"missing type", func(s *SensorToCreate) { s.Type = nil }, true},
		{"invalid type enum", func(s *SensorToCreate) { s.Type = swag.String("invalid") }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := *valid
			tt.mutate(&s)
			err := s.Validate(strfmt.Default)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
