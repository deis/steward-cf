package lib

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/pborman/uuid"
)

func TestJSONObjectString(t *testing.T) {
	t.Skip("TODO")
}

func TestJSONObjectToFromStringRoundTrip(t *testing.T) {
	t.Run("full jsonObject", func(t *testing.T) {
		jso := jsonObject(map[string]interface{}{
			"key1":     "val1",
			"key2":     "val2",
			"key3":     "val3",
			uuid.New(): uuid.New(),
		})
		jsoStr := jso.EncodeToString()
		jsoDecoded, err := jsonObjectFromString(jsoStr)
		assert.NoErr(t, err)
		assert.Equal(t, len(jsoDecoded), len(jso), "decoded json object length")
	})
	t.Run("empty jsonObject", func(t *testing.T) {
		jso := jsonObject(map[string]interface{}{})
		jsoStr := jso.EncodeToString()
		jsoDecoded, err := jsonObjectFromString(jsoStr)
		assert.NoErr(t, err)
		assert.Equal(t, len(jsoDecoded), len(jso), "decoded json object length")
	})
}
