package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
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

const DefaultTimeout = 60 * time.Minute

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
	return fmt.Errorf("occured while setting %s property for %s, %w", field, resource, err)
}

func GenerateImmutableError(resource, field string) error {
	return fmt.Errorf("%s property is immutable for %s", field, resource)
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

// DiffEmptyIps suppress difference when empty value for array is overwritten by API and assigned an actual IP address
func DiffEmptyIps(_, old, new string, _ *schema.ResourceData) bool {
	if old != "" && new == "" {
		return true
	}
	return false
}

// ApiResponseInfo - interface over different ApiResponse types from sdks
type ApiResponseInfo interface {
	HttpNotFound() bool
	LogInfo()
}

// ResourceReadyFunc polls api to see if resource exists based on id
type ResourceReadyFunc func(ctx context.Context, d *schema.ResourceData) (bool, error)

// WaitForResourceToBeReady - keeps retrying until resource is ready(true is returned), or until err is thrown, or ctx is cancelled
func WaitForResourceToBeReady(ctx context.Context, d *schema.ResourceData, fn ResourceReadyFunc) error {
	if d.Id() == "" {
		return fmt.Errorf("resource with id %s not ready, still trying ", d.Id())
	}
	err := resource.RetryContext(ctx, DefaultTimeout, func() *resource.RetryError {
		isReady, err := fn(ctx, d)
		if isReady == true {
			return nil
		}
		if err != nil {
			resource.NonRetryableError(err)
		}
		log.Printf("[DEBUG] resource with id %s not ready, still trying ", d.Id())
		return resource.RetryableError(fmt.Errorf("resource with id %s not ready, still trying ", d.Id()))
	})
	return err
}

// IsResourceDeletedFunc polls api to see if resource exists based on id
type IsResourceDeletedFunc func(ctx context.Context, d *schema.ResourceData) (bool, error)

// WaitForResourceToBeDeleted - keeps retrying until resource is not found(404), or until ctx is cancelled
func WaitForResourceToBeDeleted(ctx context.Context, d *schema.ResourceData, fn IsResourceDeletedFunc) error {

	err := resource.RetryContext(ctx, DefaultTimeout, func() *resource.RetryError {
		isDeleted, err := fn(ctx, d)
		if isDeleted {
			return nil
		}
		if err != nil {
			resource.NonRetryableError(err)
		}
		log.Printf("[DEBUG] resource with id %s still has not been deleted", d.Id())
		return resource.RetryableError(fmt.Errorf("resource with id %s found, still trying ", d.Id()))
	})
	return err
}

// DecodeInterfaceToStruct can decode from interface{}, or from []interface
// will turn "" into nil values
// takes snake_case fields and decodes them into camelcase fields of struct
// used to decode values from TypeList and TypeSet of schema(`d`) directly into sdk structs
func DecodeInterfaceToStruct(input, output interface{}) error {
	config := mapstructure.DecoderConfig{
		DecodeHook:       PointerEmptyToNil(),
		ErrorUnused:      false,
		ErrorUnset:       false,
		ZeroFields:       false,
		WeaklyTypedInput: true,
		MatchName:        IsSnakeEqualToCamelCase,
		Result:           &output,
	}
	customDecoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] rawdata to decode %s \n", input)
	err = customDecoder.Decode(input)
	if err != nil {
		return err
	}
	return err
}

func IsSnakeEqualToCamelCase(a, b string) bool {
	return strings.EqualFold(strcase.ToCamel(a), b)
}

func PointerEmptyToNil() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() == reflect.String && data == "" {
			return nil, nil
		}
		return data, nil
	}
}

// checks if value['1'] of key[`id`] is present inside a slice of maps[string]interface{}
func IsValueInSliceOfMap[T comparable](sliceOfMaps []interface{}, key string, value T) bool {
	for _, mmap := range sliceOfMaps {
		//do not delete if the id in the old rule is present in the new rules to be updated
		if value == mmap.(map[string]interface{})[key] {
			return true
		}
	}
	return false
}
