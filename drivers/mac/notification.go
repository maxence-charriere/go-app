// +build darwin,amd64

package mac

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Notification implements the app.Element interface.
type Notification struct {
	core.Elem

	id string

	onReply func(reply string)
}

func newNotification(c app.NotificationConfig) *Notification {
	n := &Notification{
		id: uuid.New().String(),

		onReply: c.OnReply,
	}

	if n.onReply != nil {
		driver.Elems.Put(n)
	}

	err := driver.Platform.Call("notifications.New", nil, struct {
		ID        string
		Title     string
		Subtitle  string
		Text      string
		ImageName string
		Sound     bool
		Reply     bool
	}{
		ID:        n.ID(),
		Title:     c.Title,
		Subtitle:  c.Subtitle,
		Text:      c.Text,
		ImageName: c.ImageName,
		Sound:     c.Sound,
		Reply:     c.OnReply != nil,
	})

	n.SetErr(err)
	return n
}

// ID satisfies the app.Element interface.
func (n *Notification) ID() string {
	return n.id
}

func onNotificationReply(n *Notification, in map[string]interface{}) {
	if reply := in["Reply"].(string); n.onReply != nil && len(reply) != 0 {
		n.onReply(reply)
	}

	driver.Elems.Delete(n)
}

func handleNotification(h func(n *Notification, in map[string]interface{})) core.GoHandler {
	return func(in map[string]interface{}) {
		id, _ := in["ID"].(string)

		e := driver.Elems.GetByID(id)
		if e.Err() == app.ErrElemNotSet {
			return
		}

		h(e.(*Notification), in)
	}
}
