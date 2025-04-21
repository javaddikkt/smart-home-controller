package models

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"testing"
	"time"
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

	t.Run("valid sensor", func(t *testing.T) {
		err := validSensor.Validate(strfmt.Default)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("missing serial number", func(t *testing.T) {
		sensor := *validSensor
		sensor.SerialNumber = nil
		err := sensor.Validate(strfmt.Default)
		if err == nil {
			t.Errorf("expected error for missing serial number")
		}
	})

	t.Run("invalid serial number format", func(t *testing.T) {
		sensor := *validSensor
		sensor.SerialNumber = swag.String("abcd")
		err := sensor.Validate(strfmt.Default)
		if err == nil {
			t.Errorf("expected error for invalid serial number format")
		}
	})

	t.Run("invalid sensor type", func(t *testing.T) {
		sensor := *validSensor
		sensor.Type = swag.String("invalid-type")
		err := sensor.Validate(strfmt.Default)
		if err == nil {
			t.Errorf("expected error for invalid type enum")
		}
	})

	t.Run("missing ID", func(t *testing.T) {
		sensor := *validSensor
		sensor.ID = nil
		err := sensor.Validate(strfmt.Default)
		if err == nil {
			t.Errorf("expected error for missing ID")
		}
	})

	t.Run("ID less than minimum", func(t *testing.T) {
		sensor := *validSensor
		sensor.ID = swag.Int64(0)
		err := sensor.Validate(strfmt.Default)
		if err == nil {
			t.Errorf("expected error for ID < 1")
		}
	})
}
