package app

import (
	"net/url"
	"testing"

	"github.com/murlokswarm/log"
)

type ShareTest struct{}

func (s *ShareTest) Text(v string) {
	log.Info("sharing text:", v)
}

func (s *ShareTest) URL(v *url.URL) {
	log.Info("sharing URL:", v)
}

func TestShare(t *testing.T) {
	Share().Text("Hello Murlok!")
}
