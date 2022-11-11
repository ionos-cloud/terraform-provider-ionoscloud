package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CreateTransport - creates customizable transport for http clients
func CreateTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		DisableKeepAlives:     true,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   3,
		MaxConnsPerHost:       3,
	}
}

func DiffSlice(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

// DiffSliceOneWay returns the elements in `a` that aren't in `b`.
func DiffSliceOneWay(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func GenerateSetError(resource, field string, err error) error {
	return fmt.Errorf("an error occured while setting %s property for %s, %w", field, resource, err)
}

func SetPropWithNilCheck(m map[string]interface{}, prop string, v interface{}) {

	rVal := reflect.ValueOf(v)
	if rVal.Kind() == reflect.Ptr {
		if !rVal.IsNil() {
			m[prop] = rVal.Elem().Interface()
		}
	} else {
		m[prop] = v
	}
}

func GenerateEmail() string {
	email := fmt.Sprintf("terraform_test-%d@mailinator.com", time.Now().UnixNano())
	return email
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func TestNotEmptySlice(resource, attribute string) resource.TestCheckFunc {
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

func TestValueInSlice(resource, attribute, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resource {
				continue
			}

			lengthOfSlice, err := strconv.Atoi(rs.Primary.Attributes[attribute])

			if err != nil {
				return err
			} else if lengthOfSlice <= 0 {
				return fmt.Errorf("returned %s slice is empty", attribute)
			} else {
				for i := 0; i < lengthOfSlice; i++ {
					attribute = attribute[:len(attribute)-1] + strconv.Itoa(i)
					if rs.Primary.Attributes[attribute] == value {
						return nil
					}
				}
			}

		}
		return fmt.Errorf("value %s not in %s slice", value, attribute)
	}
}

func TestImageNotNull(resource, attribute string) resource.TestCheckFunc {
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

func CheckFileExists(filePath string) bool {
	_, err := os.Open(filePath) // For read access.
	return err == nil
}

// WriteToFile - creates the file and writes 'value' to it.
func WriteToFile(name, value string) error {
	file, err := os.Create(name)
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("[DEBUG] could not close file %v", err)
		}
	}()

	if err != nil {
		return err
	}
	_, err = file.WriteString(value)
	return err
}

// DiffWithoutNewLines terraform suppress differences between newlines
func DiffWithoutNewLines(_, old, new string, _ *schema.ResourceData) bool {
	old = RemoveNewLines(old)
	new = RemoveNewLines(new)
	return strings.EqualFold(old, new)
}

func RemoveNewLines(s string) string {
	newlines := regexp.MustCompile(`(?:\r\n?|\n)*\z`)
	return newlines.ReplaceAllString(s, "")
}

// DiffToLower terraform suppress differences between lower and upper
func DiffToLower(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
