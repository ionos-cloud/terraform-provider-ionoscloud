package mongov2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestFormatOptionalTime(t *testing.T) {
	ts := time.Date(2026, 11, 23, 9, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    *time.Time
		ok       bool
		expected types.String
	}{
		{
			name:     "set value is formatted as RFC3339",
			input:    &ts,
			ok:       true,
			expected: types.StringValue("2026-11-23T09:30:00Z"),
		},
		{
			name:     "not ok returns null",
			input:    nil,
			ok:       false,
			expected: types.StringNull(),
		},
		{
			name:     "ok but nil pointer returns null",
			input:    nil,
			ok:       true,
			expected: types.StringNull(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, formatOptionalTime(tt.input, tt.ok))
		})
	}
}
