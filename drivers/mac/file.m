
#include "file.h"
#include "driver.h"
#include "json.h"

@implementation FilePanel
+ (bridge_result)newFilePanel:(NSURLComponents *)url
                      payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];
  NSString *returnID = [url queryValue:@"return-id"];

  NSDictionary *config = [JSONDecoder decodeObject:payload];
  BOOL multipleSelection = [config[@"multiple-selection"] boolValue];
  BOOL ignoreDirectories = [config[@"ignore-directories"] boolValue];
  BOOL ignoreFiles = [config[@"ignore-files"] boolValue];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    NSString *err = nil;

    @try {
      NSOpenPanel *panel = [NSOpenPanel openPanel];
      [panel setAllowsMultipleSelection:multipleSelection];
      [panel setCanChooseDirectories:!ignoreDirectories];
      [panel setCanChooseFiles:!ignoreFiles];

      id onComplete = ^(NSInteger result) {
        NSMutableArray<NSString *> *filenames = [[NSMutableArray alloc] init];

        if (result == NSModalResponseOK) {
          for (NSURL *url in panel.URLs) {
            [filenames addObject:url.path];
          }
        }

        [driver.golang
            request:[NSString stringWithFormat:@"/file/panel/select?id=%@", ID]
            payload:[JSONEncoder encodeObject:filenames]];
      };

      NSWindow *win = NSApp.keyWindow;
      if (win == nil) {
        [panel beginWithCompletionHandler:onComplete];
      } else {
        [panel beginSheetModalForWindow:win completionHandler:onComplete];
      }
    } @catch (NSException *exception) {
      err = exception.reason;
    }

    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
  });
  return make_bridge_result(nil, nil);
}
@end