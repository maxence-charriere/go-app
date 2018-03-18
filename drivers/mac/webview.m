#include "webview.h"
#include "driver.h"
#include "json.h"

@implementation AppWebView
- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender {
  NSPasteboard *pboard = [sender draggingPasteboard];

  if ([[pboard types] containsObject:NSFilenamesPboardType]) {
    Driver *driver = [Driver current];

    NSDictionary *in = @{
      @"Filenames" : [pboard propertyListForType:NSFilenamesPboardType],
    };

    [driver.goRPC call:@"driver.OnFileDrop" withInput:in onUI:YES];
  }
  return [super prepareForDragOperation:sender];
}
@end