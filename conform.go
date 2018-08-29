package conform

import (
	"encoding/json"
	"github.com/miracl/conflate"
)

// Conformer provides operations to convert data conforming to one schema, into data conforming to another schema.
// See Conform() for details.
type Conformer struct {
	Schema  *conflate.Schema
	Updater Updater
	Next    *Conformer
}

// Conform attempts to ensure the given data conforms to 'Schema'. If it does conform, then nil is returned.
// If the data does not conform to 'Schema', then the data is conformed to another schema using the 'Next' Conformer, if one is defined.
// Following this, the 'Updater' is used to ensure that the data conforms to 'Schema'.
func (c Conformer) Conform(data interface{}) error {
	err := c.validate(data)
	if err == nil {
		return nil
	}
	if c.Next == nil {
		return err
	}
	cerr := c.Next.Conform(data)
	if cerr != nil {
		// failed to conform to other schema, so we return the original validation error
		return err
	}
	uerr := c.Updater.Do(data)
	if uerr != nil {
		// failed to conform to other schema, so we return the original validation error
		return err
	}
	return c.validate(data)
}

func deepCopy(in interface{}, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return nil
	}
	return json.Unmarshal(b, out)
}

func (c Conformer) validate(data interface{}) error {
	var def interface{}
	err := deepCopy(data, &def)
	if err != nil {
		return err
	}
	err = c.Schema.ApplyDefaults(&def)
	if err != nil {
		return err
	}
	return c.Schema.Validate(def)
}
