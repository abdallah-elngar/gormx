package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateField_Valid(t *testing.T) {
	valid := []string{
		"name",
		"user_id",
		"first_name",
		"users.id",
		"t1.name",
		"field_1",
	}

	for _, f := range valid {
		t.Run(f, func(t *testing.T) {
			assert.NoError(t, ValidateField(f))
		})
	}
}

func TestValidateField_Invalid(t *testing.T) {
	invalid := []string{
		"",
		"name; DROP TABLE",
		"name'",
		"name\"",
		"name--",
		"1field",
		"field name",
		"select",
		"drop",
	}

	for _, f := range invalid {
		t.Run(f, func(t *testing.T) {
			assert.Error(t, ValidateField(f))
		})
	}
}

func TestValidateField_TooLong(t *testing.T) {
	long := make([]byte, 200)
	for i := range long {
		long[i] = 'a'
	}
	assert.Error(t, ValidateField(string(long)))
}

func TestValidateLookup(t *testing.T) {
	valid := []string{"gt", "gte", "lt", "lte", "in", "contains", "between"}
	for _, l := range valid {
		assert.NoError(t, ValidateLookup(l))
	}

	invalid := []string{"", "bad", "drop", "select"}
	for _, l := range invalid {
		assert.Error(t, ValidateLookup(l))
	}
}

func TestSplitFieldLookup(t *testing.T) {
	tests := []struct {
		input      string
		wantField  string
		wantLookup string
	}{
		{"age__gt", "age", "gt"},
		{"name", "name", ""},
		{"user__email__contains", "user__email", "contains"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f, l := SplitFieldLookup(tt.input)
			assert.Equal(t, tt.wantField, f)
			assert.Equal(t, tt.wantLookup, l)
		})
	}
}

func TestSanitizeDirection(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"asc", "ASC", false},
		{"DESC", "DESC", false},
		{"", "", false},
		{"bad", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := SanitizeDirection(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}