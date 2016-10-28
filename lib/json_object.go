package lib

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errMissing      = errors.New("key is missing")
	emptyjsonObject = jsonObject(map[string]interface{}{})
)

type errMalformedKV struct {
	kv []string
}

type errNotAString struct {
	value interface{}
}

func (e errNotAString) Error() string {
	return fmt.Sprintf("value %v is not a string", e.value)
}

func (e errMalformedKV) Error() string {
	return fmt.Sprintf("malformed jsonObject-encoded key/value pair %s", e.kv)
}

// jsonObject is a convenience wrapper around a Go type that represents a JSON object
type jsonObject map[string]interface{}

// EmptyjsonObject returns an empty jsonObject
func EmptyjsonObject() jsonObject {
	return emptyjsonObject
}

func (j jsonObject) String(key string) (string, error) {
	i, ok := j[key]
	if !ok {
		return "", errMissing
	}
	s, ok := i.(string)
	if !ok {
		return "", errNotAString{value: i}
	}
	return s, nil
}

// MarshalText is the encoding.TextMarshaler implementation
func (j jsonObject) EncodeToString() string {
	slc := make([]string, len(j))
	i := 0
	for key, val := range j {
		slc[i] = fmt.Sprintf("%s=%s", key, val)
		i++
	}
	return strings.Join(slc, ",")
}

// jsonObjectFromString decodes a string into a jsonObject. Returns a non-nil error if the string was not a valid jsonObject
func jsonObjectFromString(str string) (jsonObject, error) {
	if len(str) == 0 {
		return jsonObject(map[string]interface{}{}), nil
	}
	mp := map[string]interface{}{}
	spl := strings.Split(str, ",")
	for _, s := range spl {
		kv := strings.Split(s, "=")
		if len(kv) != 2 {
			return jsonObject(map[string]interface{}{}), errMalformedKV{kv: kv}
		}
		key, val := kv[0], kv[1]
		mp[key] = val
	}
	return jsonObject(mp), nil
}
