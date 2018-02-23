package mac

import (
	"fmt"
	"net/url"

	"github.com/murlokswarm/app"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/bridge"
)

// Notification implements the app.Element interface.
type Notification struct {
	id      uuid.UUID
	onReply func(reply string)
}

func newNotification(config app.NotificationConfig) (n *Notification, err error) {
	n = &Notification{
		id:      uuid.New(),
		onReply: config.OnReply,
	}

	if n.onReply != nil {
		driver.elements.Add(n)
	}

	p := struct {
		Title     string `json:"title"`
		Subtitle  string `json:"subtitle"`
		Text      string `json:"text"`
		ImageName string `json:"image-name"`
		Sound     bool   `json:"sound"`
		Reply     bool   `json:"reply"`
	}{
		Title:     config.Title,
		Subtitle:  config.Subtitle,
		Text:      config.Text,
		ImageName: config.ImageName,
		Sound:     config.Sound,
		Reply:     config.OnReply != nil,
	}

	_, err = driver.macos.RequestWithAsyncResponse(
		fmt.Sprintf("/notification/new?id=%s", n.id),
		bridge.NewPayload(p),
	)
	return n, err
}

// ID satisfies the app.Element interface.
func (n *Notification) ID() uuid.UUID {
	return n.id
}

func onNotificationReply(n *Notification, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	var reply string
	p.Unmarshal(&reply)

	if reply != "(null)" {
		n.onReply(reply)
	}

	driver.elements.Remove(n)
	return nil
}
