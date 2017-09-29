#include "window.h"
#include "driver.h"
#include "json.h"

@implementation Window
+ (bridge_result)newWindow:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];
  NSString *returnID = [url queryValue:@"return-id"];

  NSDictionary *config = [JSONDecoder decodeObject:payload];
  NSString *title = config[@"title"];
  NSNumber *x = config[@"x"];
  NSNumber *y = config[@"y"];
  NSNumber *width = config[@"width"];
  NSNumber *minWidth = config[@"min-width"];
  NSNumber *maxWidth = config[@"max-width"];
  NSNumber *height = config[@"height"];
  NSNumber *minHeight = config[@"min-height"];
  NSNumber *maxHeight = config[@"max-height"];

  dispatch_async(dispatch_get_main_queue(), ^{
    NSRect rect = NSMakeRect(x.floatValue, y.floatValue, width.floatValue,
                             height.floatValue);
    NSUInteger styleMask =
        NSWindowStyleMaskTitled | NSWindowStyleMaskFullSizeContentView |
        NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable |
        NSWindowStyleMaskResizable;
    NSWindow *rawWindow =
        [[NSWindow alloc] initWithContentRect:rect
                                    styleMask:styleMask
                                      backing:NSBackingStoreBuffered
                                        defer:NO];

    Window *win = [[Window alloc] initWithWindow:rawWindow];
    win.ID = ID;
    win.windowFrameAutosaveName = title;
    win.window.delegate = win;

    Driver *driver = [Driver current];
    driver.elements[ID] = win;

    [win showWindow:nil];
    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];

  });
  return make_bridge_result(nil, nil);
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
  NSString *returnID = [url queryValue:@"return-id"];

  NSDictionary *pos = [JSONDecoder decodeObject:payload];
  NSNumber *x = pos[@"x"];
  NSNumber *y = pos[@"y"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    [win.window setFrameOrigin:NSMakePoint(x.doubleValue, y.doubleValue)];
    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
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
  NSString *returnID = [url queryValue:@"return-id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    [win.window center];
    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
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
  NSString *returnID = [url queryValue:@"return-id"];

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

    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
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
  NSString *returnID = [url queryValue:@"return-id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    [win showWindow:nil];

    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
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