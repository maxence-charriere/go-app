package main

import "github.com/murlokswarm/app"

// Notification is a component that show how notifications work.
type Notification struct {
	Actions []windowAction
}

// OnMount is the func called when the component is mounted. It initializes the
// available actions.
func (n *Notification) OnMount() {
	n.Actions = []windowAction{
		{
			Name:        "Simple notification",
			Description: "Pop up a notification that can be clicked.",
			Action: func() {
			},
		},
		{
			Name:        "Notification with reply",
			Description: "Pop up a notification that contains a text input which allows to send a reply.",
			Action: func() {
			},
		},
	}

	app.Render(n)
}

// Render returns a html string that describes the component.
func (n *Notification) Render() string {
	return `
<div class="Layout">
	<navpane current="notification">
	<div class="Window-Tracking">
		<h1>Notification</h1>
		<table>
			<tr>
				<th>Clicked at</th>
				<th>Reply</th>
			</tr>
		</table>
	</div>
	<div class="Window-Actions">
		<h1 class="TopTitle">Actions</h1>
		<div class="Window-ActionList">
			{{range $idx, $v := .Actions}}
			<!-- 
				bind format a field or method binding. The below call will give
				"Actions.0.Action"
				"Actions.1.Action"
				...
			-->
			<div class="Window-Action" onclick="{{bind "Actions" $idx "Action"}}">
				<h2>{{.Name}}</h2>
				<p>{{.Description}}</p>

				{{if .Err}} 
				<p class="Error">{{.Err.Error}}</p>
				{{end}}
			</div>
			{{else}}
				<h2>Not supported</h2>
			{{end}}
		</div>
	</div>
</div>
	`
}
