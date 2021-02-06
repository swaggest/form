package form

import (
	"reflect"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestDecoderMultipleSimultaniousParseStructRequests(t *testing.T) {
	t.Parallel()

	sc := newStructCacheMap()

	type Struct struct {
		Array []int
	}

	proceed := make(chan struct{})

	var test Struct

	sv := reflect.ValueOf(test)
	typ := sv.Type()

	for i := 0; i < 200; i++ {
		go func() {
			<-proceed

			s := sc.parseStruct(ModeImplicit, sv, typ, "form")

			NotEqual(t, s, nil)
		}()
	}

	close(proceed)
}
