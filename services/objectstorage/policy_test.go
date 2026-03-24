package objectstorage

import (
	"encoding/json"
	"testing"
)

func TestBucketPolicyPrincipal_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "object with array",
			input:    `{"AWS": ["arn:aws:iam:::user/123", "arn:aws:iam:::user/456"]}`,
			expected: []string{"arn:aws:iam:::user/123", "arn:aws:iam:::user/456"},
		},
		{
			name:     "object with single string",
			input:    `{"AWS": "arn:aws:iam:::user/123"}`,
			expected: []string{"arn:aws:iam:::user/123"},
		},
		{
			name:     "object with wildcard string",
			input:    `{"AWS": "*"}`,
			expected: []string{"*"},
		},
		{
			name:     "flat array",
			input:    `["arn:aws:iam:::user/123", "arn:aws:iam:::user/456"]`,
			expected: []string{"arn:aws:iam:::user/123", "arn:aws:iam:::user/456"},
		},
		{
			name:     "flat array with wildcard",
			input:    `["*"]`,
			expected: []string{"*"},
		},
		{
			name:     "bare wildcard string",
			input:    `"*"`,
			expected: []string{"*"},
		},
		{
			name:     "bare arn string",
			input:    `"arn:aws:iam:::user/123"`,
			expected: []string{"arn:aws:iam:::user/123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p bucketPolicyPrincipal
			if err := json.Unmarshal([]byte(tt.input), &p); err != nil {
				t.Fatalf("UnmarshalJSON returned error: %v", err)
			}
			if len(p.AWS) != len(tt.expected) {
				t.Fatalf("expected %d AWS entries, got %d: %v", len(tt.expected), len(p.AWS), p.AWS)
			}
			for i, v := range tt.expected {
				if p.AWS[i] != v {
					t.Errorf("AWS[%d] = %q, want %q", i, p.AWS[i], v)
				}
			}
		})
	}
}

func TestBucketPolicyPrincipal_UnmarshalJSON_Error(t *testing.T) {
	inputs := []string{
		`123`,
		`true`,
		`{"AWS": 123}`,
	}
	for _, input := range inputs {
		var p bucketPolicyPrincipal
		if err := json.Unmarshal([]byte(input), &p); err == nil {
			t.Errorf("expected error for input %s, got nil (AWS=%v)", input, p.AWS)
		}
	}
}

func TestBucketPolicyPrincipal_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    bucketPolicyPrincipal
		expected string
	}{
		{
			name:     "single principal marshals as string",
			input:    bucketPolicyPrincipal{AWS: []string{"arn:aws:iam:::user/123"}},
			expected: `{"AWS":"arn:aws:iam:::user/123"}`,
		},
		{
			name:     "wildcard marshals as string",
			input:    bucketPolicyPrincipal{AWS: []string{"*"}},
			expected: `{"AWS":"*"}`,
		},
		{
			name:     "multiple principals marshal as array",
			input:    bucketPolicyPrincipal{AWS: []string{"arn:aws:iam:::user/123", "arn:aws:iam:::user/456"}},
			expected: `{"AWS":["arn:aws:iam:::user/123","arn:aws:iam:::user/456"]}`,
		},
		{
			name:     "nil principals marshal as empty array",
			input:    bucketPolicyPrincipal{AWS: nil},
			expected: `{"AWS":[]}`,
		},
		{
			name:     "empty principals marshal as empty array",
			input:    bucketPolicyPrincipal{AWS: []string{}},
			expected: `{"AWS":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&tt.input)
			if err != nil {
				t.Fatalf("MarshalJSON returned error: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("got %s, want %s", string(data), tt.expected)
			}
		})
	}
}

func TestBucketPolicyPrincipal_RoundTrip(t *testing.T) {
	inputs := []string{
		`{"AWS":"*"}`,
		`{"AWS":"arn:aws:iam:::user/123"}`,
		`{"AWS":["arn:aws:iam:::user/123","arn:aws:iam:::user/456"]}`,
	}
	for _, input := range inputs {
		var p bucketPolicyPrincipal
		if err := json.Unmarshal([]byte(input), &p); err != nil {
			t.Fatalf("UnmarshalJSON(%s) error: %v", input, err)
		}
		data, err := json.Marshal(&p)
		if err != nil {
			t.Fatalf("MarshalJSON error: %v", err)
		}
		if string(data) != input {
			t.Errorf("round-trip mismatch: input=%s, output=%s", input, string(data))
		}
	}
}

func TestPoliciesSemanticEqual(t *testing.T) {
	tests := []struct {
		name      string
		a         string
		b         string
		expected  bool
		expectErr bool
	}{
		{
			name: "identical policies",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"}}]}`,
			expected: true,
		},
		{
			name: "principal flat array vs object string",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":["arn:aws:iam:::user/123"]}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"arn:aws:iam:::user/123"}}]}`,
			expected: true,
		},
		{
			name: "principal wildcard string vs flat array",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":"*"}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":["*"]}]}`,
			expected: true,
		},
		{
			name: "principal object array vs object string",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":["arn:aws:iam:::user/123"]}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"arn:aws:iam:::user/123"}}]}`,
			expected: true,
		},
		{
			name:     "different key ordering",
			a:        `{"Statement":[{"Effect":"Allow","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}}],"Version":"2012-10-17"}`,
			b:        `{"Version":"2012-10-17","Statement":[{"Principal":{"AWS":"*"},"Effect":"Allow","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::bucket/*"]}]}`,
			expected: true,
		},
		{
			name: "different effect",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Deny","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"}}]}`,
			expected: false,
		},
		{
			name: "different actions",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:PutObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"}}]}`,
			expected: false,
		},
		{
			name: "different principals",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"arn:aws:iam:::user/123"}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"arn:aws:iam:::user/456"}}]}`,
			expected: false,
		},
		{
			name:      "invalid json identical strings",
			a:         `{invalid json`,
			b:         `{invalid json`,
			expectErr: true,
		},
		{
			name:      "invalid json different strings",
			a:         `{invalid json a`,
			b:         `{invalid json b`,
			expectErr: true,
		},
		{
			name:      "one valid one invalid json",
			a:         `{"Version":"2012-10-17","Statement":[]}`,
			b:         `{invalid json`,
			expectErr: true,
		},
		{
			name: "multiple principals object array vs flat array",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":["arn:aws:iam:::user/123","arn:aws:iam:::user/456"]}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":["arn:aws:iam:::user/123","arn:aws:iam:::user/456"]}]}`,
			expected: true,
		},
		{
			name: "multiple principals same entries different order",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":["arn:aws:iam:::user/123","arn:aws:iam:::user/456"]}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":["arn:aws:iam:::user/456","arn:aws:iam:::user/123"]}}]}`,
			expected: true,
		},
		{
			name: "multiple principals different sets",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":["arn:aws:iam:::user/123","arn:aws:iam:::user/456"]}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":["arn:aws:iam:::user/123","arn:aws:iam:::user/789"]}}]}`,
			expected: false,
		},
		{
			name: "actions different order",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:PutObject","s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject","s3:PutObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"}}]}`,
			expected: true,
		},
		{
			name: "resources different order",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/b","arn:aws:s3:::bucket/a"],
				"Principal":{"AWS":"*"}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Allow","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/a","arn:aws:s3:::bucket/b"],
				"Principal":{"AWS":"*"}}]}`,
			expected: true,
		},
		{
			name: "statements different order",
			a: `{"Version":"2012-10-17","Statement":[
				{"Effect":"Deny","Action":["s3:DeleteObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}},
				{"Effect":"Allow","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}}
			]}`,
			b: `{"Version":"2012-10-17","Statement":[
				{"Effect":"Allow","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}},
				{"Effect":"Deny","Action":["s3:DeleteObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}}
			]}`,
			expected: true,
		},
		{
			name: "statements different order with different content is not equal",
			a: `{"Version":"2012-10-17","Statement":[
				{"Effect":"Deny","Action":["s3:DeleteObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}},
				{"Effect":"Allow","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}}
			]}`,
			b: `{"Version":"2012-10-17","Statement":[
				{"Effect":"Allow","Action":["s3:PutObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}},
				{"Effect":"Deny","Action":["s3:DeleteObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}}
			]}`,
			expected: false,
		},
		{
			name: "condition IPs different order",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Deny","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"},
				"Condition":{"IpAddress":["192.168.1.1","10.0.0.1"]}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Deny","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"},
				"Condition":{"IpAddress":["10.0.0.1","192.168.1.1"]}}]}`,
			expected: true,
		},
		{
			name: "condition ExcludedIPs different order",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Deny","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"},
				"Condition":{"NotIpAddress":["172.16.0.1","10.0.0.1"]}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Deny","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"},
				"Condition":{"NotIpAddress":["10.0.0.1","172.16.0.1"]}}]}`,
			expected: true,
		},
		{
			name: "condition IPs different values",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Deny","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"},
				"Condition":{"IpAddress":["192.168.1.1","10.0.0.1"]}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Deny","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"},
				"Condition":{"IpAddress":["192.168.1.1","10.0.0.2"]}}]}`,
			expected: false,
		},
		{
			name: "condition both IPs and ExcludedIPs different order",
			a: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Deny","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"},
				"Condition":{"IpAddress":["192.168.1.1","10.0.0.1"],"NotIpAddress":["172.16.0.1","10.10.0.1"]}}]}`,
			b: `{"Version":"2012-10-17","Statement":[{
				"Effect":"Deny","Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::bucket/*"],
				"Principal":{"AWS":"*"},
				"Condition":{"IpAddress":["10.0.0.1","192.168.1.1"],"NotIpAddress":["10.10.0.1","172.16.0.1"]}}]}`,
			expected: true,
		},
		{
			name:      "fails typed unmarshal (Statement as object instead of array)",
			a:         `{"Statement":{"Effect":"Allow"}}`,
			b:         `{ "Statement": { "Effect": "Allow" } }`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PoliciesSemanticEqual(tt.a, tt.b)
			if tt.expectErr {
				if err == nil {
					t.Errorf("PoliciesSemanticEqual() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("PoliciesSemanticEqual() unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("PoliciesSemanticEqual() = %v, want %v", result, tt.expected)
			}
		})
	}
}
