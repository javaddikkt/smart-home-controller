package models

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"testing"
)

func TestUserToCreate_Validate(t *testing.T) {
	valid := &UserToCreate{
		Name: swag.String("John Doe"),
	}
	tests := []struct {
		name    string
		mutate  func(u *UserToCreate)
		wantErr bool
	}{
		{"valid user to create", func(_ *UserToCreate) {}, false},
		{"missing name", func(u *UserToCreate) { u.Name = nil }, true},
		{"empty name", func(u *UserToCreate) { u.Name = swag.String("") }, true},
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
