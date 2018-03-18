// +build darwin,amd64

package mac

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa
#cgo LDFLAGS: -framework WebKit
#cgo LDFLAGS: -framework CoreImage
#cgo LDFLAGS: -framework Security
#include "bridge.h"
*/
import "C"
import (
	"net/url"
	"unsafe"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/pkg/errors"
)

//export goRequest
func goRequest(url *C.char, payload *C.char) {
	driver.golang.Request(
		C.GoString(url),
		bridge.PayloadFromString(C.GoString(payload)),
	)
}

//export goRequestWithResult
// res should be free after each call of goRequestWithResult.
func goRequestWithResult(url *C.char, payload *C.char) (res *C.char) {
	pret := driver.golang.RequestWithResponse(
		C.GoString(url),
		bridge.PayloadFromString(C.GoString(payload)),
	)

	if pret != nil {
		res = C.CString(pret.String())
	}
	return res
}

func windowHandler(h func(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating window handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			return nil
		}

		win, ok := elem.(app.Window)
		if !ok {
			panic(errors.Errorf("creating window handler failed: element with id %v is not a window", id))
		}
		return h(win.Base().(*Window), u, p)
	}
}

func menuHandler(h func(m *Menu, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating menu handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			panic(errors.Wrap(err, "creating menu handler failed"))
		}

		menu, ok := elem.(app.Menu)
		if !ok {
			panic(errors.Errorf("creating menu handler failed: element with id %v is not a menu", id))
		}
		return h(menu.Base().(*Menu), u, p)
	}
}

func filePanelHandler(h func(panel *FilePanel, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating file panel handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			panic(errors.Wrap(err, "creating file panel handler failed"))
		}

		panel, ok := elem.(*FilePanel)
		if !ok {
			panic(errors.Errorf("creating file panel handler failed: element with id %v is not a file panel", id))
		}

		return h(panel, u, p)
	}
}

func saveFilePanelHandler(h func(panel *SaveFilePanel, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating save file panel handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			panic(errors.Wrap(err, "creating save file panel handler failed"))
		}

		panel, ok := elem.(*SaveFilePanel)
		if !ok {
			panic(errors.Errorf("creating save file panel handler failed: element with id %v is not a file panel", id))
		}

		return h(panel, u, p)
	}
}

func notificationHandler(h func(n *Notification, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating notification handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			return nil
			// panic(errors.Wrap(err, "creating notification handler failed"))
		}

		notification, ok := elem.(*Notification)
		if !ok {
			panic(errors.Errorf("creating notification handler failed: element with id %v is not notification", id))
		}

		return h(notification, u, p)
	}
}

func macCall(call string) error {
	ccall := C.CString(call)
	defer C.free(unsafe.Pointer(ccall))
	C.macCall(ccall)
	return nil
}

//export macCallReturn
func macCallReturn(retID, ret, err *C.char) {
	driver.macRPC.Return(
		C.GoString(retID),
		C.GoString(ret),
		C.GoString(err),
	)
}

//export goCall
func goCall(ccall *C.char, ui C.BOOL) (cout *C.char) {
	call := C.GoString(ccall)

	if ui == 1 {
		driver.CallOnUIGoroutine(func() {
			if _, err := driver.goRPC.Call(call); err != nil {
				panic(errors.Wrap(err, "go call"))
			}
		})
		return nil
	}

	ret, err := driver.goRPC.Call(call)
	if err != nil {
		panic(errors.Wrap(err, "go call"))
	}

	// Returned string must be free in objc code.
	return C.CString(ret)
}

func stringSlice(v interface{}) []string {
	src := v.([]interface{})
	s := make([]string, 0, len(src))

	for _, item := range src {
		s = append(s, item.(string))
	}
	return s
}
