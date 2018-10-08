package dom

type node struct {
	ID          string
	ParentID    string
	CompoID     string
	Type        string
	Namespace   string
	Text        string
	Attrs       map[string]string
	ChildrenIDs []string
	Dom         *Engine
}

type change struct {
	Action     changeAction
	NodeID     string
	Type       string `json:",omitempty"`
	Namespace  string `json:",omitempty"`
	Key        string `json:",omitempty"`
	Value      string `json:",omitempty"`
	ChildID    string `json:",omitempty"`
	NewChildID string `json:",omitempty"`
}

type changeAction int

const (
	newNode changeAction = iota
	delNode
	setAttr
	delAttr
	setText
	appendChild
	removeChild
	replaceChild
)
