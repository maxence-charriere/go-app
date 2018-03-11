#include "window.h"
#include "color.h"
#include "driver.h"
#include "json.h"

@implementation Window
+ (void) new:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    NSString *ID = in[@"ID"];
    NSString *title = in[@"Title"];
    NSNumber *x = in[@"X"];
    NSNumber *y = in[@"Y"];
    NSNumber *width = in[@"Width"];
    NSNumber *minWidth = in[@"MinWidth"];
    NSNumber *maxWidth = in[@"MaxWidth"];
    NSNumber *height = in[@"Height"];
    NSNumber *minHeight = in[@"MinHeight"];
    NSNumber *maxHeight = in[@"MaxHeight"];
    NSString *backgroundColor = in[@"BackgroundColor"];
    BOOL fixedSize = [in[@"FixedSize"] boolValue];
    BOOL closeHidden = [in[@"CloseHidden"] boolValue];
    BOOL minimizeHidden = [in[@"MinimizeHidden"] boolValue];
    BOOL titlebarHidden = [in[@"TitlebarHidden"] boolValue];
    NSNumber *backgroundVibrancy = in[@"BackgroundVibrancy"];

    NSRect rect = NSMakeRect(x.floatValue, y.floatValue, width.floatValue,
                             height.floatValue);

    NSUInteger styleMask =
        NSWindowStyleMaskTitled | NSWindowStyleMaskFullSizeContentView |
        NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable |
        NSWindowStyleMaskResizable;

    if (fixedSize) {
      styleMask = styleMask & ~NSWindowStyleMaskResizable;
    }

    if (closeHidden) {
      styleMask = styleMask & ~NSWindowStyleMaskClosable;
    }

    if (minimizeHidden) {
      styleMask = styleMask & ~NSWindowStyleMaskMiniaturizable;
    }

    NSWindow *rawWindow =
        [[NSWindow alloc] initWithContentRect:rect
                                    styleMask:styleMask
                                      backing:NSBackingStoreBuffered
                                        defer:NO];

    Window *win = [[Window alloc] initWithWindow:rawWindow];
    win.ID = ID;
    win.windowFrameAutosaveName = title;
    win.window.delegate = win;

    win.window.minSize =
        NSMakeSize(minWidth.doubleValue, minHeight.doubleValue);
    win.window.maxSize =
        NSMakeSize(maxWidth.doubleValue, maxHeight.doubleValue);

    [win configBackgroundColor:backgroundColor
                      vibrancy:backgroundVibrancy.integerValue];
    [win configWebview];
    [win configTitlebar:title hidden:titlebarHidden];

    driver.elements[ID] = win;
    [win showWindow:nil];

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)configBackgroundColor:(NSString *)color
                     vibrancy:(NSVisualEffectMaterial)vibrancy {
  if (vibrancy != NSVisualEffectMaterialAppearanceBased) {
    NSVisualEffectView *visualEffectView =
        [[NSVisualEffectView alloc] initWithFrame:self.window.frame];
    visualEffectView.material = vibrancy;
    visualEffectView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
    visualEffectView.state = NSVisualEffectStateActive;

    self.window.contentView = visualEffectView;
    return;
  }

  if (color.length == 0) {
    return;
  }
  self.window.backgroundColor =
      [NSColor colorWithCIColor:[CIColor colorWithHexString:color]];
}

- (void)configWebview {
  WKUserContentController *userContentController =
      [[WKUserContentController alloc] init];
  [userContentController addScriptMessageHandler:self name:@"golangRequest"];

  WKWebViewConfiguration *conf = [[WKWebViewConfiguration alloc] init];
  conf.userContentController = userContentController;

  AppWebView *webview = [[AppWebView alloc] initWithFrame:NSMakeRect(0, 0, 0, 0)
                                            configuration:conf];
  webview.translatesAutoresizingMaskIntoConstraints = NO;
  webview.navigationDelegate = self;
  webview.UIDelegate = self;

  // Make background transparent.
  [webview setValue:@(NO) forKey:@"drawsBackground"];

  [self.window.contentView addSubview:webview];
  [self.window.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"|[webview]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              webview)]];
  [self.window.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"V:|[webview]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              webview)]];
  self.webview = webview;
}

- (void)userContentController:(WKUserContentController *)userContentController
      didReceiveScriptMessage:(WKScriptMessage *)message {
  if (![message.name isEqual:@"golangRequest"]) {
    return;
  }

  Driver *driver = [Driver current];
  [driver.golang
      request:[NSString stringWithFormat:@"/window/callback?id=%@", self.ID]
      payload:message.body];
}

- (void)configTitlebar:(NSString *)title hidden:(BOOL)isHidden {
  self.window.title = title;

  if (!isHidden) {
    return;
  }

  self.window.titleVisibility = NSWindowTitleHidden;
  self.window.titlebarAppearsTransparent = isHidden;

  WindowTitleBar *titlebar = [[WindowTitleBar alloc] init];
  titlebar.translatesAutoresizingMaskIntoConstraints = NO;

  [self.window.contentView addSubview:titlebar];
  [self.window.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"|[titlebar]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              titlebar)]];
  [self.window.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"V:|[titlebar(==22)]"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              titlebar)]];
}

+ (void)load:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    Window *win = driver.elements[in[@"ID"]];
    win.loadReturnID = returnID;
    win.loadURL = [NSURL URLWithString:in[@"LoadURL"]];
    win.baseURL = [NSURL fileURLWithPath:in[@"BaseURL"]];

    [win.webview loadHTMLString:in[@"Page"] baseURL:win.baseURL];
  });
}

- (void)webView:(WKWebView *)webView
    didFinishNavigation:(WKNavigation *)navigation {
  Driver *driver = [Driver current];

  NSString *returnID = self.loadReturnID;
  if (returnID == nil || returnID.length == 0) {
    return;
  }
  [driver.macRPC return:returnID withOutput:nil andError:nil];

  self.loadReturnID = nil;
}

- (void)webView:(WKWebView *)webView
    decidePolicyForNavigationAction:(WKNavigationAction *)navigationAction
                    decisionHandler:
                        (void (^)(WKNavigationActionPolicy))decisionHandler {
  NSURL *url = navigationAction.request.URL;

  switch (navigationAction.navigationType) {
  case WKNavigationTypeOther:
    // Allow the loadHTMLString to not be blocked.
    if ([url isEqual:self.baseURL]) {
      decisionHandler(WKNavigationActionPolicyAllow);
      return;
    }
    break;

  case WKNavigationTypeReload:
    url = self.loadURL;
    break;

  case WKNavigationTypeLinkActivated:
  case WKNavigationTypeFormSubmitted:
  case WKNavigationTypeBackForward:
  case WKNavigationTypeFormResubmitted:
  default:
    break;
  }

  Driver *driver = [Driver current];
  [driver.golang
      request:[NSString stringWithFormat:@"/window/navigate?id=%@", self.ID]
      payload:[JSONEncoder encodeString:url.absoluteString]];
  decisionHandler(WKNavigationActionPolicyCancel);
}

+ (void)render:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    Window *win = driver.elements[in[@"ID"]];
    NSString *js = [NSString stringWithFormat:@"render(%@)", in[@"Render"]];
    [win.webview evaluateJavaScript:js completionHandler:nil];

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)renderAttributes:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    Window *win = driver.elements[in[@"ID"]];
    NSString *js =
        [NSString stringWithFormat:@"renderAttributes(%@)", in[@"Render"]];
    [win.webview evaluateJavaScript:js completionHandler:nil];

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (bridge_result)position:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];
  NSString *returnID = [url queryValue:@"return-id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    NSMutableDictionary<NSString *, id> *pos =
        [[NSMutableDictionary alloc] init];
    pos[@"x"] = [NSNumber numberWithDouble:win.window.frame.origin.x];
    pos[@"y"] = [NSNumber numberWithDouble:win.window.frame.origin.y];

    NSString *payload = [JSONEncoder encodeObject:pos];
    [driver.objc asyncReturn:returnID result:make_bridge_result(payload, nil)];
  });
  return make_bridge_result(nil, nil);
}

+ (bridge_result)move:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];

  NSDictionary *pos = [JSONDecoder decodeObject:payload];
  NSNumber *x = pos[@"x"];
  NSNumber *y = pos[@"y"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    [win.window setFrameOrigin:NSMakePoint(x.doubleValue, y.doubleValue)];
  });
  return make_bridge_result(nil, nil);
}

- (void)windowDidMove:(NSNotification *)notification {
  Driver *driver = [Driver current];

  NSMutableDictionary<NSString *, id> *pos = [[NSMutableDictionary alloc] init];
  pos[@"x"] = [NSNumber numberWithDouble:self.window.frame.origin.x];
  pos[@"y"] = [NSNumber numberWithDouble:self.window.frame.origin.y];

  [driver.golang
      request:[NSString stringWithFormat:@"/window/move?id=%@", self.ID]
      payload:[JSONEncoder encodeObject:pos]];
}

+ (bridge_result)center:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    [win.window center];
  });
  return make_bridge_result(nil, nil);
}

+ (bridge_result)size:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];
  NSString *returnID = [url queryValue:@"return-id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    NSMutableDictionary<NSString *, id> *size =
        [[NSMutableDictionary alloc] init];
    size[@"width"] = [NSNumber numberWithDouble:win.window.frame.size.width];
    size[@"height"] = [NSNumber numberWithDouble:win.window.frame.size.height];

    NSString *payload = [JSONEncoder encodeObject:size];
    [driver.objc asyncReturn:returnID result:make_bridge_result(payload, nil)];
  });
  return make_bridge_result(nil, nil);
}

+ (bridge_result)resize:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];

  NSDictionary *size = [JSONDecoder decodeObject:payload];
  NSNumber *width = size[@"width"];
  NSNumber *height = size[@"height"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    CGRect frame = win.window.frame;
    frame.size.width = width.doubleValue;
    frame.size.height = height.doubleValue;

    [win.window setFrame:frame display:YES];
  });
  return make_bridge_result(nil, nil);
}

- (void)windowDidResize:(NSNotification *)notification {
  Driver *driver = [Driver current];

  NSMutableDictionary<NSString *, id> *size =
      [[NSMutableDictionary alloc] init];
  size[@"width"] = [NSNumber numberWithDouble:self.window.frame.size.width];
  size[@"height"] = [NSNumber numberWithDouble:self.window.frame.size.height];

  [driver.golang
      request:[NSString stringWithFormat:@"/window/resize?id=%@", self.ID]
      payload:[JSONEncoder encodeObject:size]];
}

+ (bridge_result)focus:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    [win showWindow:nil];
  });
  return make_bridge_result(nil, nil);
}

- (void)windowDidBecomeKey:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver.golang
      request:[NSString stringWithFormat:@"/window/focus?id=%@", self.ID]
      payload:nil];
}

- (void)windowDidResignKey:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver.golang
      request:[NSString stringWithFormat:@"/window/blur?id=%@", self.ID]
      payload:nil];
}

+ (bridge_result)toggleFullScreen:(NSURLComponents *)url
                          payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    [win.window toggleFullScreen:nil];
  });
  return make_bridge_result(nil, nil);
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver.golang
      request:[NSString stringWithFormat:@"/window/fullscreen?id=%@", self.ID]
      payload:nil];
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver.golang
      request:[NSString
                  stringWithFormat:@"/window/fullscreen/exit?id=%@", self.ID]
      payload:nil];
}

+ (bridge_result)toggleMinimize:(NSURLComponents *)url
                        payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    if (!win.window.miniaturized) {
      [win.window miniaturize:nil];
    } else {
      [win.window deminiaturize:nil];
    }
  });
  return make_bridge_result(nil, nil);
}

- (void)windowDidMiniaturize:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver.golang
      request:[NSString stringWithFormat:@"/window/minimize?id=%@", self.ID]
      payload:nil];
}

- (void)windowDidDeminiaturize:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver.golang
      request:[NSString stringWithFormat:@"/window/deminimize?id=%@", self.ID]
      payload:nil];
}

+ (bridge_result)close:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    [win.window performClose:nil];
  });
  return make_bridge_result(nil, nil);
}

- (BOOL)windowShouldClose:(NSWindow *)sender {
  Driver *driver = [Driver current];

  NSString *res = [driver.golang
      requestWithResult:[NSString
                            stringWithFormat:@"/window/close?id=%@", self.ID]
                payload:nil];
  return [JSONDecoder decodeBool:res];
}

- (void)windowWillClose:(NSNotification *)notification {
  self.window = nil;

  Driver *driver = [Driver current];
  [driver.elements removeObjectForKey:self.ID];
}
@end

@implementation WindowTitleBar
- (void)mouseDragged:(nonnull NSEvent *)theEvent {
  [self.window performWindowDragWithEvent:theEvent];
}

- (void)mouseUp:(NSEvent *)event {
  Window *win = (Window *)self.window.windowController;
  [win.webview mouseUp:event];

  if (event.clickCount == 2) {
    [win.window zoom:nil];
  }
}
@end
