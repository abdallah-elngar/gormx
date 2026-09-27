package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"FirstName", "first_name"},
		{"UserID", "user_id"},
		{"HTTPServer", "http_server"},
		{"APIKey", "api_key"},
		{"ID", "id"},
		{"", ""},
		{"already_snake", "already_snake"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToSnake(tt.input))
		})
	}
}

func TestToCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"first_name", "firstName"},
		{"user_id", "userID"},
		{"", ""},
		{"single", "single"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToCamel(tt.input))
		})
	}
}

func TestToPascal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"first_name", "FirstName"},
		{"user_id", "UserID"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToPascal(tt.input))
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "users"},
		{"box", "boxes"},
		{"city", "cities"},
		{"person", "people"},
		{"man", "men"},
		{"child", "children"},
		{"sheep", "sheep"},
		{"fish", "fish"},
		{"leaf", "leaves"},
		{"knife", "knives"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, Pluralize(tt.input))
		})
	}
}

func TestSingularize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"users", "user"},
		{"boxes", "box"},
		{"cities", "city"},
		{"people", "person"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, Singularize(tt.input))
		})
	}
}