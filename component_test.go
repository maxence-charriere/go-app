package app

import "testing"

type Go struct {
	Placeholder bool
}

func (g *Go) Render() string {
	return `
<div>
    Go !
    <Bar />
</div>
    `
}

type Ku struct {
	Placeholder bool
}

func (k *Ku) Render() string {
	return `<div>Ku !</div>`
}

func TestRegisterComponent(t *testing.T) {
	RegisterComponent(&Go{})
	RegisterComponent(&Ku{})
}

func TestRegisterComponentWithConstructor(t *testing.T) {
	RegisterComponentWithConstructor(func() Componer { return &Go{} })
}
