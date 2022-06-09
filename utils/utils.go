package utils

import (
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"
)

//CreateTransport - creates customizable transport for http clients
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
	return fmt.Errorf("an error occured while setting %s property for %s, %s", field, resource, err)
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
