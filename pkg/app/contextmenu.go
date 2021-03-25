package app

// // MenuItemNode is the interface that describes a menu node.
// type MenuItemNode interface {
// 	UI

// 	// Disabled specifies whether the menu item is disabled.
// 	Disabled(bool) MenuItemNode

// 	// Keys set the menu item keys.
// 	Keys(string) MenuItemNode

// 	// Icon set the menu item icon.
// 	// "cmdorctrl" is replaced by the platform corresponding key.
// 	Icon(string) MenuItemNode

// 	// Label set the menu item label.
// 	Label(string) MenuItemNode

// 	// OnClick calls the given handler when there is a mouse click on the element.
// 	OnClick(EventHandler) MenuItemNode

// 	// Separator specifies that the menu item is a separator.
// 	Separator() MenuItemNode

// 	// Title specifies extra information about the item.
// 	Title(string) MenuItemNode
// }

// // MenuItem returns a menu item.
// func MenuItem() MenuItemNode {
// 	return &menuItem{}
// }

// type menuItem struct {
// 	Compo

// 	Props struct {
// 		disabled  bool
// 		keys      string
// 		icon      string
// 		label     string
// 		onClick   EventHandler
// 		separator bool
// 		title     string
// 	}
// }

// func (m *menuItem) Disabled(v bool) MenuItemNode {
// 	m.Props.disabled = v
// 	return m
// }

// func (m *menuItem) Keys(k string) MenuItemNode {
// 	k = strings.ToLower(k)

// 	switch Window().Get("navigator").Get("platform").String() {
// 	case "Macintosh", "MacIntel", "MacPPC", "Mac68K":
// 		k = strings.Replace(k, "cmdorctrl", "⌘", -1)
// 		k = strings.Replace(k, "cmd", "⌘", -1)
// 		k = strings.Replace(k, "command", "⌘", -1)
// 		k = strings.Replace(k, "ctrl", "⌃", -1)
// 		k = strings.Replace(k, "control", "⌃", -1)
// 		k = strings.Replace(k, "alt", "⌥", -1)
// 		k = strings.Replace(k, "option", "⌥", -1)
// 		k = strings.Replace(k, "shift", "⇧", -1)
// 		k = strings.Replace(k, "capslock", "⇪", -1)
// 		k = strings.Replace(k, "del", "⌫", -1)
// 		k = strings.Replace(k, "delete", "⌫", -1)
// 		k = strings.Replace(k, "+", "", -1)

// 	default:
// 		k = strings.Replace(k, "cmdorctrl", "ctrl", -1)
// 		k = strings.Replace(k, "cmd", "win", -1)
// 		k = strings.Replace(k, "command", "win", -1)
// 		k = strings.Replace(k, "control", "ctrl", -1)
// 	}

// 	m.Props.keys = k
// 	return m
// }

// func (m *menuItem) Icon(src string) MenuItemNode {
// 	m.Props.icon = src
// 	return m
// }

// func (m *menuItem) Label(l string) MenuItemNode {
// 	m.Props.label = l
// 	return m
// }

// func (m *menuItem) OnClick(h EventHandler) MenuItemNode {
// 	m.Props.onClick = h
// 	return m
// }

// func (m *menuItem) Separator() MenuItemNode {
// 	m.Props.separator = true
// 	return m
// }

// func (m *menuItem) Title(t string) MenuItemNode {
// 	m.Props.title = t
// 	return m
// }

// func (m *menuItem) Render() UI {
// 	if m.Props.separator {
// 		return Div().Class("goapp-menuitem-separator")
// 	}

// 	item := Button().
// 		Class("goapp-menuitem").
// 		Disabled(m.Props.disabled).
// 		Body(
// 			If(m.Props.icon != "",
// 				Img().
// 					Class("goapp-menuitem-icon").
// 					Src(m.Props.icon)),
// 			Div().
// 				Class("goapp-menuitem-label").
// 				Body(
// 					Text(m.Props.label),
// 				),
// 			Div().
// 				Class("goapp-menuitem-keys").
// 				Body(
// 					Text(m.Props.keys),
// 				),
// 		)

// 	if m.Props.onClick != nil {
// 		item = item.OnClick(m.Props.onClick)
// 	} else {
// 		item = item.Disabled(true)
// 	}

// 	if m.Props.title != "" {
// 		item = item.Title(m.Props.title)
// 	}

// 	return item
// }

// type contextMenuLayout struct {
// 	Compo
// 	visible bool
// 	items   []MenuItemNode
// }

// func (l *contextMenuLayout) Render() UI {
// 	class := "goapp-contextmenu-hidden"
// 	if l.visible {
// 		class = "goapp-contextmenu-visible"
// 	}

// 	return Div().
// 		ID("app-contextmenu-background").
// 		Class(class).
// 		OnClick(l.onHide).
// 		Body(
// 			Div().
// 				ID("app-contextmenu").
// 				Class("goapp-contextmenu").
// 				Body(
// 					Range(l.items).
// 						Slice(func(i int) UI {
// 							return l.items[i].(UI)
// 						}),
// 				),
// 		)
// }

// func (l *contextMenuLayout) hide() {
// 	l.onHide(makeContext(l, browserPage{}), Event{Value: Null()})
// }

// func (l *contextMenuLayout) onHide(ctx Context, e Event) {
// 	l.visible = false
// 	l.items = nil
// 	l.Update()
// }

// func (l *contextMenuLayout) show(items ...MenuItemNode) {
// 	l.items = items
// 	l.visible = true
// 	l.Update()

// 	menu := Window().
// 		Get("document").
// 		Call("getElementById", "app-contextmenu")
// 	menuWidth := menu.Get("offsetWidth").Int()
// 	menuHeight := menu.Get("offsetHeight").Int()

// 	winWidth, winHeight := Window().Size()
// 	cursorX, cursorY := Window().CursorPosition()

// 	x := cursorX
// 	if x+menuWidth > winWidth {
// 		x = winWidth - menuWidth
// 	}

// 	y := cursorY
// 	if y+menuHeight > winHeight {
// 		y = winHeight - menuHeight
// 	}

// 	menu.Get("style").Set("left", strconv.Itoa(x)+"px")
// 	menu.Get("style").Set("top", strconv.Itoa(y)+"px")
// }
