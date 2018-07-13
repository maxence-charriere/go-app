// +build darwin,amd64

package mac

import (
	"github.com/murlokswarm/app"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/bridge"
	"github.com/murlokswarm/app/internal/core"
)

// Notification implements the app.Element interface.
type Notification struct {
	core.Elem

	id      uuid.UUID
	onReply func(reply string)
}

func newNotification(c app.NotificationConfig) error {
	n := &Notification{
		id:      uuid.New(),
		onReply: c.OnReply,
	}

	if n.onReply != nil {
		driver.elems.Put(n)
	}

	return driver.macRPC.Call("notifications.New", nil, struct {
		ID        string
		Title     string
		Subtitle  string
		Text      string
		ImageName string
		Sound     bool
		Reply     bool
	}{
		ID:        n.ID().String(),
		Title:     c.Title,
		Subtitle:  c.Subtitle,
		Text:      c.Text,
		ImageName: c.ImageName,
		Sound:     c.Sound,
		Reply:     c.OnReply != nil,
	})
}

// ID satisfies the app.Element interface.
func (n *Notification) ID() uuid.UUID {
	return n.id
}

func onNotificationReply(n *Notification, in map[string]interface{}) interface{} {
	if reply := in["Reply"].(string); n.onReply != nil && len(reply) != 0 {
		n.onReply(reply)
	}

	driver.elems.Delete(n)
	return nil
}

func handleNotification(h func(n *Notification, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := uuid.Parse(in["ID"].(string))

		e := driver.elems.GetByID(id)
		if e.IsNotSet() {
			return nil
		}

		return h(e.(*Notification), in)
	}
}
