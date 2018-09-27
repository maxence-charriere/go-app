#include "webview.h"
#include "driver.h"
#include "json.h"

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101300

@implementation AppWebView
- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender {
  NSPasteboard *pasteboard = [sender draggingPasteboard];
  NSMutableArray<NSString *> *filenames = [[NSMutableArray alloc] init];
  NSArray<NSURL *> *url = [pasteboard propertyListForType:NSPasteboardTypeURL];

  for (NSPasteboardItem *item in pasteboard.pasteboardItems) {
    NSData *data = [item dataForType:NSPasteboardTypeURL];
    if (data == nil) {
      continue;
    }

    NSString *uString =
        [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    NSURL *u = [[NSURL alloc] initWithString:uString];

    u = u.filePathURL;
    if (u == nil) {
      continue;
    }

    u = u.absoluteURL;
    [filenames addObject:u.path];
  }

  if (filenames.count != 0) {
    NSDictionary *in = @{
      @"Filenames" : filenames,
    };

    Driver *driver = [Driver current];
    [driver.goRPC call:@"driver.OnFileDrop" withInput:in onUI:YES];
  }

  return [super prepareForDragOperation:sender];
}
@end

#else

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

#endif