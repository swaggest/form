package form

import (
	"testing"

	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//
//
// go test -cpuprofile cpu.out
// ./validator.test -test.bench=. -test.cpuprofile=cpu.prof
// go tool pprof validator.test cpu.prof
//
//
// go test -memprofile mem.out

func TestEncoderInt(t *testing.T) {

	type TestInt struct {
		Int              int
		Int8             int8
		Int16            int16
		Int32            int32
		Int64            int64
		IntPtr           *int
		Int8Ptr          *int8
		Int16Ptr         *int16
		Int32Ptr         *int32
		Int64Ptr         *int64
		IntArray         []int
		IntPtrArray      []*int
		IntArrayArray    [][]int
		IntPtrArrayArray [][]*int
		IntMap           map[int]int
		IntPtrMap        map[*int]*int
		NoValue          int
		NoPtrValue       *int
	}

	i := int(3)
	i8 := int8(3)
	i16 := int16(3)
	i32 := int32(3)
	i64 := int64(3)

	zero := int(0)
	one := int(1)
	two := int(2)
	three := int(3)

	test := TestInt{
		Int:              i,
		Int8:             i8,
		Int16:            i16,
		Int32:            i32,
		Int64:            i64,
		IntPtr:           &i,
		Int8Ptr:          &i8,
		Int16Ptr:         &i16,
		Int32Ptr:         &i32,
		Int64Ptr:         &i64,
		IntArray:         []int{one, two, three},
		IntPtrArray:      []*int{&one, &two, &three},
		IntArrayArray:    [][]int{{one, zero, three}},
		IntPtrArrayArray: [][]*int{{&one, &zero, &three}},
		IntMap:           map[int]int{one: three, zero: two},
		IntPtrMap:        map[*int]*int{&one: &three, &zero: &two},
	}

	encoder := NewEncoder()
	values, errs := encoder.Encode(test)

	Equal(t, errs, nil)
	Equal(t, len(values), 25)

	val, ok := values["Int8"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int8"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int16"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int32"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int64"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int8"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int8Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int16Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int32Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int64Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntArray"]
	Equal(t, ok, true)
	Equal(t, len(val), 3)
	Equal(t, val[0], "1")
	Equal(t, val[1], "2")
	Equal(t, val[2], "3")

	val, ok = values["IntPtrArray[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["IntPtrArray[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["IntPtrArray[2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["IntArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	val, ok = values["IntArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntPtrArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["IntPtrArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	val, ok = values["IntPtrArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntMap[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["IntMap[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntPtrMap[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["IntPtrMap[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["NoValue"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	val, ok = values["NoPtrValue"]
	Equal(t, ok, false)
}
