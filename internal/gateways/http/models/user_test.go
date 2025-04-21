package models

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func TestUser_Validate(t *testing.T) {
	valid := &User{
		ID:   swag.Int64(1),
		Name: swag.String("John Doe"),
	}
	tests := []struct {
		name    string
		mutate  func(u *User)
		wantErr bool
	}{
		{"valid user", func(_ *User) {}, false},
		{"missing ID", func(u *User) { u.ID = nil }, true},
		{"ID less than minimum", func(u *User) { u.ID = swag.Int64(0) }, true},
		{"missing name", func(u *User) { u.Name = nil }, true},
		{"empty name", func(u *User) { u.Name = swag.String("") }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := *valid
			tt.mutate(&u)
			err := u.Validate(strfmt.Default)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
