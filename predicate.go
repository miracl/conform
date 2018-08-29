package conform

import (
	"reflect"
	"regexp"
)

// Predicate is a logical test on a given dataset
type Predicate func(data interface{}) bool

// KeyPredicate returns a predicate based on a key in a given dataset.
type KeyPredicate func(key string) Predicate

// KeyExists returns a predicate that returns true if a key exists in the given data
func KeyExists(key string) Predicate {
	return func(data interface{}) bool {
		_, err := GetKey(data, key)
		return err == nil
	}
}

// ValEmpty returns a predicate that returns true if a key contains an 'empty' value.
// An empty value is either nil, the zero value for the type, a map/slice with zero length.
func ValEmpty(key string) Predicate {
	return func(data interface{}) bool {
		val, err := GetKey(data, key)
		if err != nil {
			return false
		}
		isNil := val == nil
		isZero := reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface())
		m, ok := val.(map[string]interface{})
		isEmptyMap := ok && len(m) == 0
		s, ok := val.([]interface{})
		isEmptySlice := ok && len(s) == 0
		return isNil || isZero || isEmptyMap || isEmptySlice
	}
}

// ValEqual returns a predicate that returns true if the value of the given key in the given data, equals the given value
func ValEqual(key string, val interface{}) Predicate {
	return func(data interface{}) bool {
		v, err := GetKey(data, key)
		return err == nil && reflect.DeepEqual(v, val)
	}
}

// ValRegex returns a predicate that returns true if the value of the given key in the given data, matches the given regular expression
func ValRegex(key string, re *regexp.Regexp) Predicate {
	return func(data interface{}) bool {
		s, err := GetKeyAsString(data, key)
		return err == nil && re.MatchString(s)
	}
}

// Not returns a predicate that is the negation of the given predicate
func Not(pred Predicate) Predicate {
	return func(data interface{}) bool {
		return !pred(data)
	}
}
