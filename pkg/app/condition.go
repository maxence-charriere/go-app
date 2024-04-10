package app

// Condition represents a control structure for conditionally displaying UI
// elements. It extends the UI interface to include methods for handling
// conditional logic.
type Condition interface {
	UI

	// ElseIf sets a UI element to be displayed when the given boolean
	// expression is true and all previous conditions have been false.
	//
	// expr: Boolean expression to evaluate.
	// elem: Function that returns the UI element to display.
	ElseIf(expr bool, elem func() UI) Condition

	// ElseIfSlice sets multiple UI elements to be displayed when the given
	// boolean expression is true and all previous conditions have been false.
	//
	// expr: Boolean expression to evaluate.
	// elems: Function that returns a slice of UI elements to display.
	ElseIfSlice(expr bool, elems func() []UI) Condition

	// Else sets a UI element to be displayed as a fallback when all previous
	// conditions have been false.
	//
	// elem: Function that returns the UI element to display.
	Else(elem func() UI) Condition

	// ElseSlice sets multiple UI elements to be displayed as a fallback when
	// all previous conditions have been false.
	//
	// expr: Boolean expression to evaluate.
	// elems: Function that returns a slice of UI elements to display.
	ElseSlice(elems func() []UI) Condition

	body() []UI
}

// If returns a Condition that will display the given UI element based on the
// evaluation of the provided boolean expression.
func If(expr bool, elem func() UI) Condition {
	return IfSlice(expr, func() []UI {
		return []UI{elem()}
	})
}

// IfSlice returns a Condition that will display the given slice of UI elements
// based on the evaluation of the provided boolean expression.
func IfSlice(expr bool, elems func() []UI) Condition {
	if !expr {
		return condition{}
	}
	return condition{
		children: FilterUIElems(elems()...),
		matched:  true,
	}
}

type condition struct {
	children []UI
	matched  bool
}

func (c condition) ElseIf(expr bool, elem func() UI) Condition {
	return c.ElseIfSlice(expr, func() []UI {
		return []UI{elem()}
	})
}

func (c condition) ElseIfSlice(expr bool, elems func() []UI) Condition {
	if c.matched || !expr {
		return c
	}

	c.children = FilterUIElems(elems()...)
	c.matched = true
	return c
}

func (c condition) Else(elem func() UI) Condition {
	return c.ElseSlice(func() []UI {
		return []UI{elem()}
	})
}

func (c condition) ElseSlice(elems func() []UI) Condition {
	return c.ElseIfSlice(true, elems)
}

func (c condition) JSValue() Value {
	return nil
}

func (c condition) Mounted() bool {
	return false
}

func (c condition) body() []UI {
	return c.children
}

func (c condition) parent() UI {
	return nil
}

func (c condition) setParent(UI) UI {
	return nil
}
