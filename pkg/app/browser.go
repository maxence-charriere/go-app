package app

type browser struct {
	anchorClick Func
	popState    Func
}

func (b *browser) HandleEvents(ctx Context) {
	b.handleAnchorClick(ctx)
	b.handlePopState(ctx)
	b.handleNavigationFromJS(ctx)
	b.handleAppUpdate(ctx)
	b.handleAppInstallChange(ctx)
	b.handleAppResize(ctx)
	b.handleAppOrientationChange(ctx)
}

func (b *browser) handleAnchorClick(ctx Context) {
	b.anchorClick = FuncOf(func(this Value, args []Value) any {
		ctx.Dispatch(func(ctx Context) {
			event := Event{Value: args[0]}
			if meta := event.Get("metaKey"); meta.Truthy() && meta.Bool() {
				return
			}
			if ctrl := event.Get("ctrlKey"); ctrl.Truthy() && ctrl.Bool() {
				return
			}

			for target := event.Get("target"); target.Truthy(); target = target.Get("parentElement") {
				switch target.Get("tagName").String() {
				case "A":
					if download := target.Call("getAttribute", "download"); !download.IsNull() {
						return
					}

					event.PreventDefault()
					if href := target.Get("href"); href.Truthy() {
						ctx.Navigate(target.Get("href").String())
					}
					return

				case "BODY":
					return
				}
			}
		})

		return nil
	})
	Window().Set("onclick", b.anchorClick)
}

func (b *browser) handlePopState(ctx Context) {}

func (b *browser) handleNavigationFromJS(ctx Context) {}

func (b *browser) handleAppUpdate(ctx Context) {}

func (b *browser) handleAppInstallChange(ctx Context) {}

func (b *browser) handleAppResize(ctx Context) {}

func (b *browser) handleAppOrientationChange(ctx Context) {}
