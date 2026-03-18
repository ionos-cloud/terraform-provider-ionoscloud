package objectstorage

import (
	"encoding/json"
	"testing"
)

func TestBucketPolicyPrincipalMarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		principal bucketPolicyPrincipal
		want      string
	}{
		{
			name:      "single principal uses string AWS form",
			principal: bucketPolicyPrincipal{AWS: []string{"arn:aws:iam:::user/test"}},
			want:      `{"AWS":"arn:aws:iam:::user/test"}`,
		},
		{
			name:      "multiple principals use array AWS form",
			principal: bucketPolicyPrincipal{AWS: []string{"arn:aws:iam:::user/a", "arn:aws:iam:::user/b"}},
			want:      `{"AWS":["arn:aws:iam:::user/a","arn:aws:iam:::user/b"]}`,
		},
		{
			name:      "empty principal is deterministic empty array",
			principal: bucketPolicyPrincipal{},
			want:      `{"AWS":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, err := json.Marshal(&tt.principal)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if got := string(gotBytes); got != tt.want {
				t.Fatalf("Marshal() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPoliciesSemanticEqual_PrincipalFormatVariants(t *testing.T) {
	statePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":["arn:aws:iam:::user/test"]}]}`
	apiPolicy := `{"Statement":[{"Resource":["arn:aws:s3:::bucket/*"],"Action":["s3:GetObject"],"Effect":"Allow","Principal":{"AWS":"arn:aws:iam:::user/test"}}],"Version":"2012-10-17"}`

	if !PoliciesSemanticEqual(statePolicy, apiPolicy) {
		t.Fatalf("expected policies to be semantically equal")
	}
}

func TestPoliciesSemanticEqual_BareStringPrincipal(t *testing.T) {
	statePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":"*"}]}`
	apiPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::bucket/*"],"Principal":{"AWS":"*"}}]}`

	if !PoliciesSemanticEqual(statePolicy, apiPolicy) {
		t.Fatalf("expected bare-string and AWS-string principal forms to be semantically equal")
	}
}

func TestPoliciesSemanticEqual_InvalidJSON(t *testing.T) {
	invalidA := `{"Statement":[`
	invalidB := `{"Statement":]`

	if PoliciesSemanticEqual(invalidA, invalidB) {
		t.Fatalf("expected different invalid JSON values to be non-equal")
	}

	if !PoliciesSemanticEqual(invalidA, invalidA) {
		t.Fatalf("expected identical invalid JSON values to be equal by exact string fallback")
	}
}
