package dbaas

import (
	"fmt"
	"reflect"
)

// TODO: create a separate utils package to be used after both autoscaling and dbaas are merged on master

func generateSetError(resource, field string, err error) error {
	return fmt.Errorf("an error occured while setting %s property for %s, %s", field, resource, err)
}

func setPropWithNilCheck(m map[string]interface{}, prop string, v interface{}) {

	rVal := reflect.ValueOf(v)
	if rVal.Kind() == reflect.Ptr {
		if !rVal.IsNil() {
			m[prop] = rVal.Elem().Interface()
		}
	} else {
		m[prop] = v
	}
}
