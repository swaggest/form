// Package internal contains test fixtures.
package internal

import (
	"encoding"
	"io"
)

// DeeperEmbedded is a test structure in external package.
type DeeperEmbedded struct {
	embeddedUnexported
	*embeddedUnexportedWithExportedField
}

// SetDeeplyEmbedded sets value to an unexported pointer.
func (d *DeeperEmbedded) SetDeeplyEmbedded(s string) {
	if d.embeddedUnexportedWithExportedField == nil {
		d.embeddedUnexportedWithExportedField = new(embeddedUnexportedWithExportedField)
	}

	d.embeddedUnexported.nothingToSeeHereToo = true
	d.nothingToSeeHere = true
	d.DeeplyEmbedded = s
}

type embeddedUnexportedWithExportedField struct {
	nothingToSeeHere bool
	DeeplyEmbedded   string `form:"deeply-embedded"`
}

type embeddedUnexported struct {
	nothingToSeeHereToo bool
}

type writerWithExported struct {
	WriterExp string `form:"writer-exported"`
}

func (e writerWithExported) Write(_ []byte) (n int, err error) {
	return 0, nil
}

var _ encoding.TextMarshaler = writerWithExported{}

func (e writerWithExported) MarshalText() (text []byte, err error) {
	return []byte("hello!"), nil
}

// MakeWriterWithExported creates an instance of unexported type.
func MakeWriterWithExported() io.Writer {
	return writerWithExported{
		WriterExp: "bar",
	}
}
