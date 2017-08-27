package url

import "testing"

func TestComponentURL(t *testing.T) {
	u, err := Parse("/app.mycomponent?foo=bar")
	if err != nil {
		t.Fatal(err)
	}

	name, ok := u.Component()
	if !ok {
		t.Fatal("url should target a component")
	}
	if name != "app.mycomponent" {
		t.Fatal("component name should be app.mycomponent:", name)
	}
}

func TestNonComponentURL(t *testing.T) {
	u, err := Parse("http://google.com")
	if err != nil {
		t.Fatal(err)
	}

	_, ok := u.Component()
	if ok {
		t.Fatal("url should not target a component")
	}
}

func TestParseError(t *testing.T) {
	_, err := Parse("http://goo@$#$%@gle.com")
	if err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}
