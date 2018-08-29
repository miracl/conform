package conform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdater_SetExisting(t *testing.T) {
	u := Set("/key", "value")
	m := map[string]interface{}{"key": "old"}
	err := u(m)
	assert.Nil(t, err)
	assert.Equal(t, "value", m["key"])
}

func TestUpdater_SetMissing(t *testing.T) {
	u := Set("/key", "value")
	m := map[string]interface{}{}
	err := u(m)
	assert.Nil(t, err)
	assert.Equal(t, "value", m["key"])
}

func TestUpdater_SetNil(t *testing.T) {
	u := Set("/key", nil)
	m := map[string]interface{}{}
	err := u(m)
	assert.Nil(t, err)
	assert.Contains(t, m, "key")
	assert.Nil(t, m["key"])
}

func TestUpdater_KeyEmpty(t *testing.T) {
	b := ValEmpty("/key")(
		map[string]interface{}{
			"key": NewObject(),
		})
	assert.True(t, b)
}
