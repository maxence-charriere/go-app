#include "webview.h"
#include "driver.h"
#include "json.h"

@implementation AppWebView
- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender {

  NSPasteboard *pboard = [sender draggingPasteboard];

  if ([[pboard types] containsObject:NSFilenamesPboardType]) {
    Driver *driver = [Driver current];
    NSArray *files = [pboard propertyListForType:NSFilenamesPboardType];

    [driver.golang request:@"/driver/filedrop"
                   payload:[JSONEncoder encodeObject:files]];
  }
  return [super prepareForDragOperation:sender];
}
@end