package autoscaling

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"reflect"
)

// TODO: create a separate utils package to be used
func readPublicKey(path string) (key string, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		return "", err
	}
	return string(ssh.MarshalAuthorizedKey(pubKey)[:]), nil
}

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
