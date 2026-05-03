package app

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestPage(t *testing.T) {
	testPage(t, &requestPage{
		width:      42,
		height:     21,
		url:        &url.URL{Path: "/"},
		resolveURL: func(v string) string { return v },
	})
}

func TestBrowserPage(t *testing.T) {
	testSkipNonWasm(t)

	// client := NewClientTester(Div())
	// defer client.Close()
	// testPage(t, browserPage{
	// 	resolveStaticResource: func(v string) string { return v },
	// })
}

func testPage(t *testing.T, p Page) {
	p.SetTitle("go-app")
	require.Equal(t, "go-app", p.Title())
	p.SetTitlef("%s %d", "go-app", 10)
	require.Equal(t, "go-app 10", p.Title())

	p.SetLang("fr")
	require.Equal(t, "fr", p.Lang())

	p.SetDescription("test")
	require.Equal(t, "test", p.Description())
	p.SetDescriptionf("%s %d", "test", 10)
	require.Equal(t, "test 10", p.Description())

	p.SetAuthor("Maxence")
	require.Equal(t, "Maxence", p.Author())
	p.SetAuthorf("%s %d", "Maxence", 10)
	require.Equal(t, "Maxence 10", p.Author())

	p.SetKeywords("go", "app")
	require.Equal(t, "go, app", p.Keywords())

	p.SetLoadingLabel("loading test")
	p.SetLoadingLabelf("loading %d", 10)

	p.SetImage("image")
	require.Equal(t, "image", p.Image())
	p.SetImagef("image-%d", 10)
	require.Equal(t, "image-10", p.Image())

	u := p.URL()
	u.Path = "/test"
	p.ReplaceURL(u)
	require.Equal(t, u.String(), p.URL().String())

	w, h := p.Size()
	require.NotZero(t, w)
	require.NotZero(t, h)

	p.SetTwitterCard(TwitterCard{Card: "summary"})
	p.SetCanonicalLink("/canon")
	p.SetCanonicalLinkf("/%s", "canon")
}
