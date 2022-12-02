package form_test

import (
	"io"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/form/v5"
)

func TestEncoder_Encode_deep_embed(t *testing.T) {
	type S struct {
		io.Writer

		Header string `form:"header"`
	}

	enc := form.NewEncoder()

	enc.SetMode(form.ModeExplicit)
	enc.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Name
	})

	v := S{}
	v.Writer = form.MakeEmbeddedUnexported()
	v.Header = "foo"

	e, err := enc.Encode(v)
	require.NoError(t, err)

	assert.Equal(t, url.Values{"": []string{"hello!"}, "Header": []string{"foo"}}, e)
}
