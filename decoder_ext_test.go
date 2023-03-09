package form_test

import (
	"io"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/form/v5"
	"github.com/swaggest/form/v5/internal"
)

func TestDecoder_Decode_deep_embed(t *testing.T) {
	type S struct {
		io.Writer
		internal.DeeperEmbedded

		Header string `form:"header"`
	}

	dec := form.NewDecoder()

	dec.SetMode(form.ModeExplicit)

	s := S{}
	collect := make(map[string]interface{})
	vals := url.Values{"deeply-embedded": []string{"baz"}, "header": []string{"foo"}, "writer-exported": []string{"bar"}}

	require.NoError(t, dec.Decode(&s, vals, collect))
	assert.Equal(t, "foo", s.Header)

	assert.Equal(t, map[string]interface{}{"header": "foo"}, collect)
}
