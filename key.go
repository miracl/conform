package conform

import (
	"path"
)

// JoinKey creates a key from its component parts
func JoinKey(key ...string) string {
	return path.Join(key...)
}

// SplitKey returns the parent key and the leaf name of the given key
func SplitKey(key string) (string, string) {
	return path.Split(key)
}

// KeyName returns the leaf name of the given key
func KeyName(key string) string {
	_, name := path.Split(key)
	return name
}
