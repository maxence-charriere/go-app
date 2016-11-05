package app

import (
	"fmt"
)

const (
	jsFmt = `
function Call(msg) {
    %v
}
    `
)

func MurlokJS() string {
	return fmt.Sprintf(jsFmt, driver.JavascriptBridge())
}
