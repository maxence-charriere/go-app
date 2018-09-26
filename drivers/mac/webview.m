#include "webview.h"
#include "driver.h"
#include "json.h"

@implementation AppWebView
- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender {
  NSPasteboard *pasteboard = [sender draggingPasteboard];
  NSMutableArray<NSString *> *filenames = [[NSMutableArray alloc] init];
  NSArray<NSURL *> *url =
      [pasteboard propertyListForType:NSPasteboardTypeURL];

  for (NSPasteboardItem *item in pasteboard.pasteboardItems) {
    NSData *data = [item dataForType:NSPasteboardTypeFileURL];
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