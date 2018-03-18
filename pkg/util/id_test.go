package util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNames(t *testing.T) {
	t.Run("[]string", func(t *testing.T) {
		var value interface{}
		json.Unmarshal([]byte(`["test1", "test2"]`), &value)

		names := GetNames(value)
		assert.Equal(t, []string{"test1", "test2"}, names)
	})
}

func TestGetName(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		var value interface{}
		json.Unmarshal([]byte(`"test"`), &value)

		name := GetName(value)
		assert.Equal(t, "test", name)
	})

	t.Run("object name", func(t *testing.T) {
		var value interface{}
		json.Unmarshal([]byte(`{"name":"test"}`), &value)

		id := GetName(value)
		assert.Equal(t, "test", id)
	})
}

func TestGetID(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		var value interface{}
		json.Unmarshal([]byte("123"), &value)

		id := GetID(value)
		assert.Equal(t, 123, id)
	})

	t.Run("null", func(t *testing.T) {
		var value interface{}
		json.Unmarshal([]byte("null"), &value)

		id := GetID(value)
		assert.Equal(t, 0, id)
	})

	t.Run("int", func(t *testing.T) {
		id := GetID(interface{}(123))
		assert.Equal(t, 123, id)
	})

	t.Run("object id", func(t *testing.T) {
		var value interface{}
		json.Unmarshal([]byte(`{"id":123, "name":"test"}`), &value)

		id := GetID(value)
		assert.Equal(t, 123, id)
	})
}
