package models

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func TestSensorToUserBinding_Validate(t *testing.T) {
	valid := &SensorToUserBinding{
		SensorID: swag.Int64(1),
	}
	tests := []struct {
		name    string
		mutate  func(b *SensorToUserBinding)
		wantErr bool
	}{
		{"valid binding", func(_ *SensorToUserBinding) {}, false},
		{"missing sensor ID", func(b *SensorToUserBinding) { b.SensorID = nil }, true},
		{"sensor ID less than minimum", func(b *SensorToUserBinding) { b.SensorID = swag.Int64(0) }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := *valid
			tt.mutate(&b)
			err := b.Validate(strfmt.Default)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
