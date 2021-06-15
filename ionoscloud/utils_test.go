package ionoscloud

import (
	"testing"
)

func TestSetPropWithNilCheck(t *testing.T) {
	m := map[string]interface{}{}

	var pBnil *bool
	var pSnil *string

	s := "foo"
	s2 := "bar"
	b := true
	b2 := false
	pS := &s
	pB := &b

	setPropWithNilCheck(m, "bool_nil", pBnil)
	setPropWithNilCheck(m, "string_nil", pSnil)
	setPropWithNilCheck(m, "bool_ok", pB)
	setPropWithNilCheck(m, "string_ok", pS)
	setPropWithNilCheck(m, "string_simple", s2)
	setPropWithNilCheck(m, "bool_simple", b2)

	if _, ok := m["bool_nil"]; ok {
		t.Error("bool_nil was set")
	}

	if _, ok := m["string_nil"]; ok {
		t.Error("string_nil was set")
	}

	if m["bool_ok"] != b {
		t.Errorf("bool_ok != %+v", b)
	}

	if m["string_ok"] != s {
		t.Errorf("string_ok: %+v != %+v", m["string_ok"], s)
	}

	if m["bool_simple"] != b2 {
		t.Errorf("bool_simple != %+v", b2)
	}

	if m["string_simple"] != s2 {
		t.Errorf("string_simple != %+v", s2)
	}

}
