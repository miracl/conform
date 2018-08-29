package conform

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonpointer"
	"regexp"
	"strconv"
	"text/template"
)

var (
	jsonPtr  = gojsonpointer.NewJsonPointer
	keyRegex = regexp.MustCompile(`{{(\s*key\s.*?)}}`)
)

const (
	// note: we use these special delimiters to cater for the case where a key's value needs to contain a golang template string
	openDelim  = `{{<<`
	closeDelim = `>>}}`
)

// GetKey returns the value of the given key in the given dataset
func GetKey(data interface{}, key string) (interface{}, error) {
	p, err := jsonPtr(key)
	if err != nil {
		return nil, err
	}
	val, _, err := p.Get(data)
	if err != nil {
		return nil, err
	}
	return val, err
}

// GetKeyAsString returns the value of the given key in the given dataset
func GetKeyAsString(data interface{}, key string) (string, error) {
	val, err := GetKey(data, key)
	if err != nil {
		return "", nil
	}
	s, ok := val.(string)
	if !ok {
		return "", errors.Errorf("Value is not a string '%v'", val)
	}
	return s, nil
}

// SetKey sets the value of the given key in the given dataset
// If 'val' is a string it can optionally contain '{{ key /my/path }}' placeholders.
func SetKey(data interface{}, key string, val interface{}) error {
	p, err := jsonPtr(key)
	if err != nil {
		return err
	}
	if s, ok := val.(string); ok {
		s, err = renderTemplate(data, s)
		if err != nil {
			return err
		}
		val = s
	}
	_, err = p.Set(data, val)
	return err
}

// DeleteKey removes the given key from the given dataset
func DeleteKey(data interface{}, key string) error {
	p, err := jsonPtr(key)
	if err != nil {
		return err
	}
	_, err = p.Delete(data)
	return err
}

// MoveKey moves the key 'from' to the key 'to' in the given dataset
func MoveKey(data interface{}, from string, to string) error {
	err := CopyKey(data, from, to)
	if err != nil {
		return err
	}
	return DeleteKey(data, from)
}

// TransformKey allows the caller to specify a custom operation to the value of a given key in the dataset
func TransformKey(data interface{}, key string, f func(val interface{}) (interface{}, error)) error {
	val, err := GetKey(data, key)
	if err != nil {
		return err
	}
	val, err = f(val)
	if err != nil {
		return err
	}
	return SetKey(data, key, val)
}

// renderTemplate replaces any '{{ key /my/path }}' placeholders in the tmpl string, with the value of data["my"]["path"]
func renderTemplate(data interface{}, tmpl string) (string, error) {
	s := keyRegex.ReplaceAllString(tmpl, openDelim+"$1"+closeDelim)
	if s == tmpl {
		return s, nil
	}
	funcs := map[string]interface{}{
		"key": func(key string) (interface{}, error) {
			return GetKey(data, key)
		},
	}
	t, err := template.New("").Delims(openDelim, closeDelim).Funcs(funcs).Parse(s)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = t.Execute(&b, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// RegexKey performs a regular expression replacement operation to the value of the given key in the given dataset.
// The 'repl' string can optionally contain '{{ key /my/path }}' placeholders.
func RegexKey(data interface{}, key string, re *regexp.Regexp, repl string) error {
	return TransformKey(data, key,
		func(val interface{}) (interface{}, error) {
			s, ok := val.(string)
			if !ok {
				return nil, errors.Errorf("Key is not a string value '%v'", key)
			}
			r, err := renderTemplate(data, repl)
			if err != nil {
				return nil, err
			}
			return re.ReplaceAllString(s, r), nil
		})
}

// RegexMatchKey performs a regular expression replacement operation to the value of the given key in the given dataset
// The 'match' or 'repl' strings can optionally contain '{{ key /my/path }}' placeholders.
func RegexMatchKey(data interface{}, key string, match string, repl string) error {
	m, err := renderTemplate(data, match)
	if err != nil {
		return err
	}
	re, err := regexp.Compile(m)
	if err != nil {
		return err
	}
	return RegexKey(data, key, re, repl)
}

// CopyKey copies the key 'from' to the key 'to' in the given dataset
func CopyKey(data interface{}, from string, to string) error {
	val, err := GetKey(data, from)
	if err != nil {
		return err
	}
	return SetKey(data, to, val)
}

// WalkKey calls the given KeyUpdater for each child key of the given key in the dataset. It does not recurse.
func WalkKey(data interface{}, key string, f KeyUpdater) error {
	val, err := GetKey(data, key)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	} else if m, ok := val.(map[string]interface{}); ok {
		for k := range m {
			err = f(JoinKey(key, k)).Do(data)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if s, ok := val.([]interface{}); ok {
		for i := range s {
			k := strconv.Itoa(i)
			err = f(JoinKey(key, k)).Do(data)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return errors.Errorf("Cannot walk key '%v'", key)
}
