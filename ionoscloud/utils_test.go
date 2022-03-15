package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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

	utils.SetPropWithNilCheck(m, "bool_nil", pBnil)
	utils.SetPropWithNilCheck(m, "string_nil", pSnil)
	utils.SetPropWithNilCheck(m, "bool_ok", pB)
	utils.SetPropWithNilCheck(m, "string_ok", pS)
	utils.SetPropWithNilCheck(m, "string_simple", s2)
	utils.SetPropWithNilCheck(m, "bool_simple", b2)

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
func testNotEmptySlice(resource, attribute string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resource {
				continue
			}

			lengthOfSlice := rs.Primary.Attributes[attribute]

			if lengthOfSlice == "0" {
				return fmt.Errorf("returned version slice is empty")
			}
		}
		return nil
	}
}

func testImageNotNull(resource, attribute string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resource {
				continue
			}

			image := rs.Primary.Attributes[attribute]

			if image == "" {
				return fmt.Errorf("%s is empty, expected an UUID", attribute)
			} else if !IsValidUUID(image) {
				return fmt.Errorf("%s should be a valid UUID, got: %#v", attribute, image)
			}

		}
		return nil
	}
}
