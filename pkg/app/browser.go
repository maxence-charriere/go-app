package app

type browser struct {
	AppUpdatable bool

	anchorClick      Func
	popState         Func
	navigationFromJS Func
	appUpdate        Func
}

func (b *browser) HandleEvents(ctx nodeContext, notifyComponentEvent func(any)) {
	b.handleAnchorClick(ctx)
	b.handlePopState(ctx)
	b.handleNavigationFromJS(ctx)
	b.handleAppUpdate(ctx, notifyComponentEvent)
	b.handleAppInstallChange(ctx)
	b.handleAppResize(ctx)
	b.handleAppOrientationChange(ctx)
}

func (b *browser) handleAnchorClick(ctx nodeContext) {
	b.anchorClick = FuncOf(func(this Value, args []Value) any {
		ctx.dispatch(func() {
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

func (b *browser) handlePopState(ctx nodeContext) {
	b.popState = FuncOf(func(this Value, args []Value) any {
		ctx.dispatch(func() {
			ctx.navigate(Window().URL(), false)
		})
		return nil
	})
	Window().Set("onpopstate", b.popState)
}

func (b *browser) handleNavigationFromJS(ctx nodeContext) {
	b.navigationFromJS = FuncOf(func(this Value, args []Value) any {
		ctx.dispatch(func() {
			ctx.Navigate(args[0].String())
		})
		return nil
	})
	Window().Set("goappNav", b.navigationFromJS)
}

func (b *browser) handleAppUpdate(ctx nodeContext, notifyComponentEvent func(any)) {
	b.appUpdate = FuncOf(func(this Value, args []Value) any {
		ctx.dispatch(func() {
			b.AppUpdatable = true
			notifyComponentEvent(appUpdate{})
		})
		return nil
	})
	Window().Set("goappOnUpdate", b.appUpdate)
}

func (b *browser) handleAppInstallChange(ctx Context) {}

func (b *browser) handleAppResize(ctx Context) {}

func (b *browser) handleAppOrientationChange(ctx Context) {}
