package conform

import (
	"encoding/json"
	"github.com/miracl/conflate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testData(t *testing.T, data string) interface{} {
	var obj interface{}
	err := json.Unmarshal([]byte(data), &obj)
	assert.Nil(t, err)
	return obj
}

func testSchema(t *testing.T, data string) *conflate.Schema {
	s, err := conflate.NewSchemaData([]byte(data))
	assert.Nil(t, err)
	return s
}

func TestConform_MoveProperty(t *testing.T) {
	schema1 := `
{
	"title": "schema1",
	"type": "object",
	"properties": {
		"key1": {
			"type": "string",
			"default": "def1"
		}
	},
	"additionalProperties": false
}`
	schema2 := `
{
	"title": "schema2",
	"type": "object",
	"properties": {
			"key2": {
				"type": "string"
			}
	},
	"additionalProperties": false
}`

	data := testData(t, `{ "key2": "val2"}`)
	expData := testData(t, `{ "key1": "val2"}`)

	c := Conformer{
		Schema:  testSchema(t, schema1),
		Updater: Move("/key2", "/key1"),
	}
	c.Next = &Conformer{
		Schema: testSchema(t, schema2),
	}

	err := c.Conform(data)
	assert.Nil(t, err)
	assert.Equal(t, expData, data)
}

func TestConform_SetProperty(t *testing.T) {
	schema1 := `
{
	"title": "schema2",
	"type": "object",
	"properties": {
			"key": {
				"type": "string"
			}
	},
	"additionalProperties": false
}`
	schema2 := `
{
	"title": "schema1",
	"type": "object",
	"properties": {
		"key1": {
			"type": "string"
		},
		"key2": {
			"type": "string"
		}
	},
	"additionalProperties": false
}`

	data := testData(t, `{ "key1": "val1", "key2": "val2"}`)
	expData := testData(t, `{ "key": "val1val2"}`)

	c := Conformer{
		Schema: testSchema(t, schema1),
		Updater: Compose(
			Set("/key", `{{key "/key1"}}{{key "/key2"}}`),
			Delete("/key1"),
			Delete("/key2"),
		),
	}
	c.Next = &Conformer{
		Schema: testSchema(t, schema2),
	}

	err := c.Conform(data)
	assert.Nil(t, err)
	assert.Equal(t, expData, data)
}

func TestConform_RegexProperty(t *testing.T) {
	schema1 := `
{
	"title": "schema2",
	"type": "object",
	"properties": {
			"key": {
				"type": "string"
			}
	},
	"additionalProperties": false
}`
	schema2 := `
{
	"title": "schema1",
	"type": "object",
	"properties": {
		"key1": {
			"type": "string"
		},
		"key2": {
			"type": "string"
		}
	},
	"additionalProperties": false
}`

	data := testData(t, `{ "key1": "val1", "key2": "val2"}`)
	expData := testData(t, `{ "key": "val1val2"}`)

	c := Conformer{
		Schema: testSchema(t, schema1),
		Updater: Compose(
			Set("/key", ""),
			RegexMatch("/key", `.*`, `{{key "/key1"}}{{key "/key2"}}`),
			Delete("/key1"),
			Delete("/key2"),
		),
	}
	c.Next = &Conformer{
		Schema: testSchema(t, schema2),
	}

	err := c.Conform(data)
	assert.Nil(t, err)
	assert.Equal(t, expData, data)
}

func TestConform_ModifyArrayItems(t *testing.T) {
	schema1 := `
{
	"title": "schema1",
	"type": "object",
	"properties": {
		"array1": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"key1": {
						"type": "string"
					}
				},
				"additionalProperties": false
			}
		}
	}
}`
	schema2 := `
{
	"title": "schema2",
	"type": "object",
	"properties": {
		"array1": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"key2": {
						"type": "string"
					}
				},
				"additionalProperties": false
			}
		}
	}
}`

	data := testData(t, `{ "array1": [ {"key2": "value1"}, {"key2": "value2"} ] }`)
	expData := testData(t, `{ "array1": [ {"key1": "val1"}, {"key1": "val2"} ] }`)

	c := Conformer{
		Schema: testSchema(t, schema1),
		Updater: Walk("/array1",
			func(key string) Updater {
				return Move(JoinKey(key, "key2"), JoinKey(key, "key1")).Then(
					RegexMatch(JoinKey(key, "/key1"), `^.*(\d+)$`, "val$1"))
			}),
	}
	c.Next = &Conformer{
		Schema: testSchema(t, schema2),
	}

	err := c.Conform(data)
	assert.Nil(t, err)
	assert.Equal(t, expData, data)
}
