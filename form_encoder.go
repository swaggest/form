package form

import (
	"bytes"
	"errors"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

// DEPRECATED
// Use EncodeFunc
// EncodeCustomTypeFunc allows for registering/overriding types to be parsed.
type EncodeCustomTypeFunc func(x interface{}) ([]string, error)

// EncodeFunc allows for registering/overriding types to be parsed.
type EncodeFunc func(x interface{}) (string, error)

// EncodeErrors is a map of errors encountered during form encoding
type EncodeErrors map[string]error

func (e EncodeErrors) Error() string {
	buff := bytes.NewBufferString(blank)

	for k, err := range e {
		buff.WriteString(fieldNS)
		buff.WriteString(k)
		buff.WriteString(errorText)
		buff.WriteString(err.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

// An InvalidEncodeError describes an invalid argument passed to Encode.
type InvalidEncodeError struct {
	Type reflect.Type
}

func (e *InvalidEncodeError) Error() string {

	if e.Type == nil {
		return "form: Encode(nil)"
	}

	return "form: Encode(nil " + e.Type.String() + ")"
}

// Encoder is the main encode instance
type Encoder struct {
	tagName         string
	structCache     *structCacheMap
	customTypeFuncs map[reflect.Type]EncodeFunc
	dataPool        *sync.Pool
	mode            Mode
	embedAnonymous  bool
}

// NewEncoder creates a new encoder instance with sane defaults
func NewEncoder() *Encoder {

	e := &Encoder{
		tagName:        "form",
		mode:           ModeImplicit,
		structCache:    newStructCacheMap(),
		embedAnonymous: true,
	}

	e.dataPool = &sync.Pool{New: func() interface{} {
		return &encoder{
			e:         e,
			namespace: make([]byte, 0, 64),
		}
	}}

	return e
}

// SetTagName sets the given tag name to be used by the encoder.
// Default is "form"
func (e *Encoder) SetTagName(tagName string) {
	e.tagName = tagName
}

// SetMode sets the mode the encoder should run
// Default is ModeImplicit
func (e *Encoder) SetMode(mode Mode) {
	e.mode = mode
}

// SetAnonymousMode sets the mode the encoder should run
// Default is AnonymousEmbed
func (e *Encoder) SetAnonymousMode(mode AnonymousMode) {
	e.embedAnonymous = mode == AnonymousEmbed
}

// RegisterTagNameFunc registers a custom tag name parser function
// NOTE: This method is not thread-safe it is intended that these all be registered prior to any parsing
//
// ADDITIONAL: once a custom function has been registered the default, or custom set, tag name is ignored
// and relies 100% on the function for the name data. The return value WILL BE CACHED and so return value
// must be consistent.
func (e *Encoder) RegisterTagNameFunc(fn TagNameFunc) {
	e.structCache.tagFn = fn
}

// RegisterFunc registers a EncodeFunc against a number of types
// NOTE: this method is not thread-safe it is intended that these all be registered prior to any parsing
func (e *Encoder) RegisterFunc(fn EncodeFunc, types ...interface{}) {

	if e.customTypeFuncs == nil {
		e.customTypeFuncs = map[reflect.Type]EncodeFunc{}
	}

	for _, t := range types {
		e.customTypeFuncs[reflect.TypeOf(t)] = fn
	}
}

// DEPRECATED
// Use RegisterFunc
// RegisterCustomTypeFunc registers a CustomTypeFunc against a number of types
// NOTE: this method is not thread-safe it is intended that these all be registered prior to any parsing
func (e *Encoder) RegisterCustomTypeFunc(fn EncodeCustomTypeFunc, types ...interface{}) {

	if e.customTypeFuncs == nil {
		e.customTypeFuncs = map[reflect.Type]EncodeFunc{}
	}

	for _, t := range types {
		e.customTypeFuncs[reflect.TypeOf(t)] = func(x interface{}) (string, error) {
			res, err := fn(x)
			if err != nil {
				return "", err
			}
			if len(res) > 0 {
				return res[0], err
			}
			return "", errors.New("empty result")
		}
	}
}

// Encode encodes the given values and sets the corresponding struct values
func (e *Encoder) Encode(v interface{}) (values url.Values, err error) {

	val, kind := ExtractType(reflect.ValueOf(v))

	if kind == reflect.Ptr || kind == reflect.Interface || kind == reflect.Invalid {
		return nil, &InvalidEncodeError{reflect.TypeOf(v)}
	}

	enc := e.dataPool.Get().(*encoder)
	enc.values = make(url.Values)

	if kind == reflect.Struct && val.Type() != timeType {
		enc.traverseStruct(val, enc.namespace[0:0], -1)
	} else {
		enc.setFieldByType(val, enc.namespace[0:0], -1, false)
	}

	if len(enc.errs) > 0 {
		err = enc.errs
		enc.errs = nil
	}

	values = enc.values

	e.dataPool.Put(enc)

	return
}
