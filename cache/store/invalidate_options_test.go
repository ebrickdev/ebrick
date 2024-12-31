package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidateOptionsTags(t *testing.T) {
	tests := []struct {
		name     string
		input    InvalidateOptions
		expected []string
	}{
		{
			name: "Single tag",
			input: InvalidateOptions{
				Tags: []string{"tag1"},
			},
			expected: []string{"tag1"},
		},
		{
			name: "Multiple tags",
			input: InvalidateOptions{
				Tags: []string{"tag1", "tag2", "tag3"},
			},
			expected: []string{"tag1", "tag2", "tag3"},
		},
		{
			name: "No tags",
			input: InvalidateOptions{
				Tags: []string{},
			},
			expected: []string{},
		},
		{
			name: "Nil tags",
			input: InvalidateOptions{
				Tags: nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.Tags)
		})
	}
}
