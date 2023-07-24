package form_test

import (
	"encoding/json"
	"fmt"
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

func TestDecoder_Decode_nested(t *testing.T) {
	type S struct {
		Foo    string `form:"foo"`
		Deeper struct {
			Bar    int `form:"bar"`
			Deeper struct {
				Baz bool `form:"baz"`
			} `form:"deeper"`
		} `form:"deeper"`
	}

	vals := url.Values{
		"foo":                 []string{"abc"},
		"deeper[bar]":         []string{"123"},
		"deeper[deeper][baz]": []string{"true"},
	}
	s := S{}
	collect := make(map[string]interface{})
	dec := form.NewDecoder()
	dec.SetNamespacePrefix("[")
	dec.SetNamespaceSuffix("]")
	dec.SetMode(form.ModeExplicit)

	require.NoError(t, dec.Decode(&s, vals, collect))
	assert.Equal(t, "abc", s.Foo)
	assert.Equal(t, 123, s.Deeper.Bar)
	assert.Equal(t, true, s.Deeper.Deeper.Baz)

	fmt.Printf("%#v\n", collect)

	expected := map[string]interface{}{
		"deeper": s.Deeper,
		"foo":    "abc",
	}

	assert.Equal(t, expected, collect)
}

func TestDecoder_Decode_queryForm(t *testing.T) {
	type jsonFilter struct {
		Foo string `json:"foo" maxLength:"5"`
	}

	type deepObjectFilter struct {
		Bar string `query:"bar" minLength:"3"`
	}

	type inputQueryObject struct {
		Query            map[int]float64  `query:"in_query" description:"Object value in query."`
		JSONFilter       jsonFilter       `query:"json_filter" description:"JSON object value in query."`
		DeepObjectFilter deepObjectFilter `query:"deep_object_filter" description:"Deep object value in query params."`
	}

	vals := url.Values{
		"in_query[1]":             []string{"1"},
		"in_query[2]":             []string{"2"},
		"in_query[3]":             []string{"3"},
		"json_filter":             []string{`{"foo":"strin"}`},
		"deep_object_filter[bar]": []string{"sd"},
	}

	s := inputQueryObject{}
	collect := make(map[string]interface{})

	dec := form.NewDecoder()
	dec.SetTagName("query")
	dec.SetNamespacePrefix("[")
	dec.SetNamespaceSuffix("]")
	dec.SetMode(form.ModeExplicit)
	dec.RegisterFunc(func(s string) (interface{}, error) {
		var j jsonFilter
		err := json.Unmarshal([]byte(s), &j)

		return j, err
	}, jsonFilter{})

	require.NoError(t, dec.Decode(&s, vals, collect))

	expected := map[string]interface{}{
		"deep_object_filter": deepObjectFilter{Bar: "sd"},
		"in_query":           map[int]float64{1: 1, 2: 2, 3: 3},
		"json_filter":        jsonFilter{Foo: "strin"},
	}

	assert.Equal(t, expected, collect)
}
