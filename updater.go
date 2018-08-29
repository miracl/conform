package conform

import (
	"regexp"
)

// Updater defines an operation performed on a given dataset
type Updater func(data interface{}) error

// KeyUpdater defines an Updater that is based on the given key in the dataset
type KeyUpdater func(key string) Updater

// NewObject creates a new empty object that represents an 'object' type in the underlying datasource (i.e. JSON/YAML/TOML)
func NewObject() map[string]interface{} {
	return map[string]interface{}{}
}

// NewArray creates a new empty object that represents ian 'array' type in the underlying datasource (i.e. JSON/YAML/TOML)
func NewArray(data ...interface{}) []interface{} {
	return append([]interface{}{}, data...)
}

// Set sets the value of the given key in the given dataset
func Set(key string, val interface{}) Updater {
	return func(data interface{}) error {
		return SetKey(data, key, val)
	}
}

// Delete removes the given key from the given dataset
func Delete(key string) Updater {
	return func(data interface{}) error {
		return DeleteKey(data, key)
	}
}

// Move moves the key 'from' to the key 'to' in the given dataset
func Move(from string, to string) Updater {
	return func(data interface{}) error {
		return MoveKey(data, from, to)
	}
}

// Copy copies the key 'from' to the key 'to' in the given dataset
func Copy(from string, to string) Updater {
	return func(data interface{}) error {
		return CopyKey(data, from, to)
	}
}

// Transform allows the caller to specify a custom operation to the value of a given key in the dataset
func Transform(key string, f func(val interface{}) (interface{}, error)) Updater {
	return func(data interface{}) error {
		return TransformKey(data, key, f)
	}
}

// Regex performs a regular expression replacement operation to the value of the given key in the given dataset
func Regex(key string, re *regexp.Regexp, repl string) Updater {
	return func(data interface{}) error {
		return RegexKey(data, key, re, repl)
	}
}

// RegexMatch performs a regular expression replacement operation to the value of the given key in the given dataset
func RegexMatch(key string, match string, repl string) Updater {
	return func(data interface{}) error {
		return RegexMatchKey(data, key, match, repl)
	}
}

// Walk calls the given KeyUpdater for each child key of the given key in the dataset. It does not recurse.
func Walk(key string, u KeyUpdater) Updater {
	return func(data interface{}) error {
		return WalkKey(data, key, u)
	}
}

// Do is a convenience function to apply the updater to the given dataset. It returns nil if the updater is nil.
func (u Updater) Do(data interface{}) error {
	if u != nil {
		return u(data)
	}
	return nil
}

// Do is a convenience function to return an updater for the given key. If the keyupdater is nil then a nil updater is returned.
func (u KeyUpdater) Do(key string) Updater {
	if u != nil {
		return u(key)
	}
	return nil
}

// Then executes the next updater if the current updater succeeds. If either updater fails an error is returned.
func (u Updater) Then(next Updater) Updater {
	return func(data interface{}) error {
		err := u.Do(data)
		if err != nil {
			return err
		}
		return next.Do(data)
	}
}
