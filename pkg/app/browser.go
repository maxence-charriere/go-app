package app

import (
	"time"
)

type browser struct {
	AppUpdatable bool

	anchorClick      Func
	popState         Func
	navigationFromJS Func
	appUpdate        Func
	appInstallChange Func
	appResize        Func
	resizeTimer      *time.Timer
}

func (b *browser) HandleEvents(ctx Context, notifyComponentEvent func(any)) {
	b.handleAnchorClick(ctx)
	b.handlePopState(ctx)
	b.handleNavigationFromJS(ctx)
	b.handleAppUpdate(ctx, notifyComponentEvent)
	b.handleAppInstallChange(ctx, notifyComponentEvent)
	b.handleAppResize(ctx, notifyComponentEvent)
}

func (b *browser) handleAnchorClick(ctx Context) {
	b.anchorClick = FuncOf(func(this Value, args []Value) any {
		ctx.dispatch(func() {
			event := Event{Value: args[0]}

			for target := event.Get("target"); target.Truthy(); target = target.Get("parentElement") {
				switch target.Get("tagName").String() {
				case "A":
					if meta := event.Get("metaKey"); meta.Truthy() && meta.Bool() {
						return
					}

					if ctrl := event.Get("ctrlKey"); ctrl.Truthy() && ctrl.Bool() {
						return
					}

					if download := target.Call("getAttribute", "download"); !download.IsNull() {
						return
					}

					switch target.Get("target").String() {
					case "_blank":
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

func (b *browser) handlePopState(ctx Context) {
	b.popState = FuncOf(func(this Value, args []Value) any {
		ctx.dispatch(func() {
			ctx.navigate(Window().URL(), false)
		})
		return nil
	})
	Window().Set("onpopstate", b.popState)
}

func (b *browser) handleNavigationFromJS(ctx Context) {
	b.navigationFromJS = FuncOf(func(this Value, args []Value) any {
		ctx.dispatch(func() {
			ctx.Navigate(args[0].String())
		})
		return nil
	})
	Window().Set("goappNav", b.navigationFromJS)
}

func (b *browser) handleAppUpdate(ctx Context, notifyComponentEvent func(any)) {
	appUpdate := func() {
		ctx.dispatch(func() {
			b.AppUpdatable = true
			notifyComponentEvent(appUpdate{})
		})
		ctx.defere(func() {
			Log(Window().URL().Hostname() + " has been updated, reload to see changes")
		})
	}

	b.appUpdate = FuncOf(func(this Value, args []Value) any {
		appUpdate()
		return nil
	})
	Window().Set("goappOnUpdate", b.appUpdate)

	if Window().Get("goappUpdatedBeforeWasmLoaded").Truthy() {
		appUpdate()
	}
}

func (b *browser) handleAppInstallChange(ctx Context, notifyComponentEvent func(any)) {
	appInstallChange := func() {
		ctx.dispatch(func() {
			notifyComponentEvent(appInstallChange{})
		})
	}

	b.appInstallChange = FuncOf(func(this Value, args []Value) any {
		appInstallChange()
		return nil
	})
	Window().Set("goappOnAppInstallChange", b.appInstallChange)

	if Window().Get("goappAppInstallChangedBeforeWasmLoaded").Truthy() {
		appInstallChange()
	}
}

func (b *browser) handleAppResize(ctx Context, notifyComponentEvent func(any)) {
	const resizeCooldown = time.Millisecond * 250

	b.appResize = FuncOf(func(this Value, args []Value) any {
		ctx.dispatch(func() {
			if b.resizeTimer != nil {
				b.resizeTimer.Stop()
				b.resizeTimer.Reset(resizeCooldown)
				return
			}

			b.resizeTimer = time.AfterFunc(resizeCooldown, func() {
				ctx.dispatch(func() {
					notifyComponentEvent(resize{})
				})
			})
		})
		return nil
	})
	Window().Set("onresize", b.appResize)
}
