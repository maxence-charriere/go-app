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

    Driver *driver = [Driver current];
    driver.elements[ID] = win;

    [win showWindow:nil];

    [driver.objc returnFor:returnID result:make_bridge_result(nil, nil)];

  });
  return make_bridge_result(nil, nil);
}

+ (bridge_result)position:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];
  NSString *returnID = [url queryValue:@"return-id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Window *win = driver.elements[ID];

    NSMutableDictionary<NSString *, id> *res =
        [[NSMutableDictionary alloc] init];
    res[@"x"] = [NSNumber numberWithDouble:win.window.frame.origin.x];
    res[@"y"] = [NSNumber numberWithDouble:win.window.frame.origin.y];

    NSString *payload = [JSONEncoder encodeObject:res];
    [driver.objc returnFor:returnID result:make_bridge_result(payload, nil)];
  });
  return make_bridge_result(nil, nil);
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