#include "window.h"

@implementation Window
+ (bridge_result)newWindow:(NSURL *)url payload:(NSString *)payload {
  NSLog(@"Should create at cocoa window here");
  return make_bridge_result();
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