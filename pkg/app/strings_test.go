package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringToString(t *testing.T) {
	var s string
	stringTo("hello", &s)
	require.Equal(t, "hello", s)
}

func TestStringToInt(t *testing.T) {
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64

	err := stringTo("42", &i)
	require.NoError(t, err)
	require.Equal(t, 42, i)

	err = stringTo("42", &i8)
	require.NoError(t, err)
	require.Equal(t, int8(42), i8)

	err = stringTo("42", &i16)
	require.NoError(t, err)
	require.Equal(t, int16(42), i16)

	err = stringTo("-42", &i32)
	require.NoError(t, err)
	require.Equal(t, int32(-42), i32)

	err = stringTo("42", &i64)
	require.NoError(t, err)
	require.Equal(t, int64(42), i64)

	err = stringTo("", &i)
	require.NoError(t, err)
	require.Equal(t, 0, i)
}

func TestStringToUInt(t *testing.T) {
	var i uint
	var i8 uint8
	var i16 uint16
	var i32 uint32
	var i64 uint64

	err := stringTo("42", &i)
	require.NoError(t, err)
	require.Equal(t, uint(42), i)

	err = stringTo("42", &i8)
	require.NoError(t, err)
	require.Equal(t, uint8(42), i8)

	err = stringTo("42", &i16)
	require.NoError(t, err)
	require.Equal(t, uint16(42), i16)

	err = stringTo("42", &i32)
	require.NoError(t, err)
	require.Equal(t, uint32(42), i32)

	err = stringTo("42", &i64)
	require.NoError(t, err)
	require.Equal(t, uint64(42), i64)

	err = stringTo("", &i)
	require.NoError(t, err)
	require.Equal(t, uint(0), i)

	err = stringTo("-42", &i)
	require.NoError(t, err)
	require.Equal(t, uint(0), i)
}

func TestStringToFloat(t *testing.T) {
	var f64 float64
	var f32 float32

	err := stringTo("42.1", &f64)
	require.NoError(t, err)
	require.Equal(t, float64(42.1), f64)

	err = stringTo("21.2", &f32)
	require.NoError(t, err)
	require.Equal(t, float32(21.2), f32)

	err = stringTo("", &f64)
	require.NoError(t, err)
	require.Equal(t, float64(0), f64)
}
