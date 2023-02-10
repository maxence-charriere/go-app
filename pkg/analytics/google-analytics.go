package analytics

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// GoogleAnalyticsHeader returns the header to use in the app.Handler.RawHeader
// field to initialize Google Analytics.
func GoogleAnalyticsHeader(propertyID string) string {
	return fmt.Sprintf(`<!-- Global site tag (gtag.js) - Google Analytics -->
	<script async src="https://www.googletagmanager.com/gtag/js?id=%s"></script>
	<script>
	  window.dataLayer = window.dataLayer || [];
	  function gtag(){dataLayer.push(arguments);}
	  gtag('js', new Date());
	
	  gtag('config', '%s', {'send_page_view': false});
	</script>`, propertyID, propertyID)
}

func NewGoogleAnalytics() Backend {
	return googleAnalytics{}
}

type googleAnalytics struct {
}

func (a googleAnalytics) Identify(userID string, traits map[string]interface{}) {
	a.gtag("set", map[string]interface{}{
		"user_id": userID,
	})
}

func (a googleAnalytics) Track(event string, properties map[string]interface{}) {
	a.gtag("event", event, properties)
}

func (a googleAnalytics) Page(name string, properties map[string]interface{}) {
	a.gtag("event", "page_view", map[string]interface{}{
		"page_title":    properties["title"],
		"page_location": properties["url"],
		"page_path":     properties["path"],
	})
}

func (a googleAnalytics) gtag(args ...interface{}) {
	gtag := app.Window().Get("gtag")
	if !gtag.Truthy() {
		return
	}
	app.Window().Call("gtag", args...)
}
