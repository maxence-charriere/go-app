package app_test

import (
	"fmt"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestEmit(t *testing.T) {
	app.Emit("hello")
}

func TestImport(t *testing.T) {
	app.Import(&tests.Foo{})

	defer func() { recover() }()
	app.Import(tests.NoPointerCompo{})
}

func TestLog(t *testing.T) {
	log := ""

	app.Logger = func(format string, a ...interface{}) {
		log = fmt.Sprintf(format, a...)
	}

	app.Log("hello", "world")
	assert.Equal(t, "hello world", log)

	app.Logf("%s %s", "bye", "world")
	assert.Equal(t, "bye world", log)
}

func TestPanic(t *testing.T) {
	log := ""

	app.Logger = func(format string, a ...interface{}) {
		log = fmt.Sprintf(format, a...)
	}

	defer func() {
		err := recover()
		assert.Equal(t, "hello world", log)
		assert.Equal(t, "hello world", err)
	}()

	app.Panic("hello", "world")
	assert.Fail(t, "no panic")
}

func TestPanicf(t *testing.T) {
	log := ""

	app.Logger = func(format string, a ...interface{}) {
		log = fmt.Sprintf(format, a...)
	}

	defer func() {
		err := recover()
		assert.Equal(t, "bye world", log)
		assert.Equal(t, "bye world", err)
	}()

	app.Panicf("%s %s", "bye", "world")
	assert.Fail(t, "no panic")
}

func TestPretty(t *testing.T) {
	t.Log(app.Pretty(struct {
		Hello string
		World string
	}{}))
}
