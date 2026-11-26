package logging

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogElem(tag string, labels ...string) map[string]any {
	labelsAny := make([]any, len(labels))
	for i, l := range labels {
		labelsAny[i] = l
	}
	return map[string]any{
		"source":   "kubernetes",
		"tag":      tag,
		"protocol": "http",
		"labels":   schema.NewSet(schema.HashString, labelsAny),
		"destinations": []any{
			map[string]any{"type": "loki", "retention_in_days": 7},
		},
	}
}

func TestExpandPipelineLogs(t *testing.T) {
	tests := []struct {
		name    string
		logs    []any
		oldLogs []any
		// wantLabelsJSON is the JSON value of 'labels' for each log, "" when the field must be omitted
		wantLabelsJSON []string
	}{
		{
			name:           "create without labels omits the field",
			logs:           []any{testLogElem("a")},
			wantLabelsJSON: []string{""},
		},
		{
			name:           "create with labels sends them sorted",
			logs:           []any{testLogElem("a", "zone", "app")},
			wantLabelsJSON: []string{`["app","zone"]`},
		},
		{
			name:           "update removing labels sends an empty list",
			logs:           []any{testLogElem("a")},
			oldLogs:        []any{testLogElem("a", "app")},
			wantLabelsJSON: []string{`[]`},
		},
		{
			name:           "update without labels before or after omits the field",
			logs:           []any{testLogElem("a")},
			oldLogs:        []any{testLogElem("a")},
			wantLabelsJSON: []string{""},
		},
		{
			name:           "update changing labels sends the new labels",
			logs:           []any{testLogElem("a", "env")},
			oldLogs:        []any{testLogElem("a", "app")},
			wantLabelsJSON: []string{`["env"]`},
		},
		{
			name:           "update adding a log without labels omits the field for it",
			logs:           []any{testLogElem("a", "app"), testLogElem("b")},
			oldLogs:        []any{testLogElem("a", "app")},
			wantLabelsJSON: []string{`["app"]`, ""},
		},
		{
			// regression: matching by index compared new 'b' with old 'a' and did not clear b's labels
			name:           "update removing a log before one whose labels are removed clears them",
			logs:           []any{testLogElem("b")},
			oldLogs:        []any{testLogElem("a"), testLogElem("b", "app")},
			wantLabelsJSON: []string{`[]`},
		},
		{
			name:           "update reordering logs and removing labels clears only the matching log",
			logs:           []any{testLogElem("b"), testLogElem("a")},
			oldLogs:        []any{testLogElem("a"), testLogElem("b", "app")},
			wantLabelsJSON: []string{`[]`, ""},
		},
		{
			// matching by index would send an empty list for 'b' because old 'a' had labels
			name:           "update removing a log with labels does not clear the log that shifts into its position",
			logs:           []any{testLogElem("b")},
			oldLogs:        []any{testLogElem("a", "app"), testLogElem("b")},
			wantLabelsJSON: []string{""},
		},
		{
			name:           "update renaming the tag of a log with labels treats it as a new log",
			logs:           []any{testLogElem("renamed")},
			oldLogs:        []any{testLogElem("a", "app")},
			wantLabelsJSON: []string{""},
		},
		{
			name:           "missing labels key is handled",
			logs:           []any{map[string]any{"source": "docker", "tag": "a", "protocol": "tcp", "destinations": []any{}}},
			wantLabelsJSON: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandPipelineLogs(tt.logs, tt.oldLogs)
			require.Len(t, got, len(tt.wantLabelsJSON))
			for i, logElem := range got {
				serialized, err := logElem.ToMap()
				require.NoError(t, err)
				labels, present := serialized["labels"]
				if tt.wantLabelsJSON[i] == "" {
					assert.False(t, present, "log %d: labels must be omitted, got %v", i, labels)
					continue
				}
				require.True(t, present, "log %d: labels must be sent", i)
				labelsJSON, err := json.Marshal(labels)
				require.NoError(t, err)
				assert.JSONEq(t, tt.wantLabelsJSON[i], string(labelsJSON), "log %d", i)
			}
		})
	}
}

func TestExpandPipelineLogsFields(t *testing.T) {
	got := expandPipelineLogs([]any{testLogElem("mytag", "app")}, nil)
	require.Len(t, got, 1)
	assert.Equal(t, "kubernetes", got[0].Source)
	assert.Equal(t, "mytag", got[0].Tag)
	assert.Equal(t, "http", got[0].Protocol)
	require.Len(t, got[0].Destinations, 1)
	assert.Equal(t, "loki", got[0].Destinations[0].Type)
	assert.Equal(t, int32(7), got[0].Destinations[0].RetentionInDays)
}
