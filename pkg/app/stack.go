package app

// UIStack is the interface that describes a container that displays its items
// as stacked panels.
//
// EXPERIMENTAL WIDGET.
type UIStack interface {
	UI

	// ID sets the stack root HTML element id property.
	ID(string) UIStack

	// Class adds a CSS class to the stack root HTML element class property.
	Class(string) UIStack

	// Left aligns the content on the left.
	Left() UIStack

	// Center aligns the content on the horizontal center.
	Center() UIStack

	// Right aligns the content on the right.
	Right() UIStack

	// Top aligns the content on the top.
	Top() UIStack

	// Middle aligns the content on the vertical center.
	Middle() UIStack

	// Bottom aligns the content on the bottom.
	Bottom() UIStack

	// Stretch stretches the content vertically.
	Stretch() UIStack

	// Content sets the content with the given UI elements.
	Content(elems ...UI) UIStack
}

// Stack creates a container that displays its items as stacked panels.
//
// EXPERIMENTAL WIDGET.
func Stack() UIStack {
	return &stack{
		IhorizontalAlign: "flex-start",
		IverticalAlign:   "flex-start",
	}
}

type stack struct {
	Compo

	Iid              string
	Iclass           string
	IhorizontalAlign string
	IverticalAlign   string
	Icontent         []UI
}

func (s *stack) ID(v string) UIStack {
	s.Iid = v
	return s
}

func (s *stack) Left() UIStack {
	s.IhorizontalAlign = "flex-start"
	return s
}

func (s *stack) Center() UIStack {
	s.IhorizontalAlign = "center"
	return s
}

func (s *stack) Right() UIStack {
	s.IhorizontalAlign = "flex-end"
	return s
}

func (s *stack) Top() UIStack {
	s.IverticalAlign = "flex-start"
	return s
}

func (s *stack) Middle() UIStack {
	s.IverticalAlign = "center"
	return s
}

func (s *stack) Bottom() UIStack {
	s.IverticalAlign = "flex-end"
	return s
}

func (s *stack) Stretch() UIStack {
	s.IverticalAlign = "stretch"
	return s
}

func (s *stack) Class(v string) UIStack {
	if v == "" {
		return s
	}
	if s.Iclass != "" {
		s.Iclass += " "
	}
	s.Iclass += v
	return s
}

func (s *stack) Content(elems ...UI) UIStack {
	s.Icontent = FilterUIElems(elems...)
	return s
}

func (s *stack) Render() UI {
	return Div().
		DataSet("goapp", "Stack").
		ID(s.Iid).
		Class(s.Iclass).
		Style("display", "flex").
		Style("justify-content", s.IhorizontalAlign).
		Style("align-items", s.IverticalAlign).
		Body(s.Icontent...)
}
