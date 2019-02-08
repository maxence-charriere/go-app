package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmit(t *testing.T) {
	Emit("hello")
}

func TestImport(t *testing.T) {
	Import(&Foo{})

	defer func() { recover() }()
	Import(NoPointerCompo{})
}

func TestLog(t *testing.T) {
	log := ""

	Logger = func(format string, a ...interface{}) {
		log = fmt.Sprintf(format, a...)
	}

	Log("hello", "world")
	assert.Equal(t, "hello world", log)

	Logf("%s %s", "bye", "world")
	assert.Equal(t, "bye world", log)
}

func TestPanic(t *testing.T) {
	log := ""

	Logger = func(format string, a ...interface{}) {
		log = fmt.Sprintf(format, a...)
	}

	defer func() {
		err := recover()
		assert.Equal(t, "hello world", log)
		assert.Equal(t, "hello world", err)
	}()

	Panic("hello", "world")
	assert.Fail(t, "no panic")
}

func TestPanicf(t *testing.T) {
	log := ""

	Logger = func(format string, a ...interface{}) {
		log = fmt.Sprintf(format, a...)
	}

	defer func() {
		err := recover()
		assert.Equal(t, "bye world", log)
		assert.Equal(t, "bye world", err)
	}()

	Panicf("%s %s", "bye", "world")
	assert.Fail(t, "no panic")
}
