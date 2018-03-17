
#include "file.h"
#include "driver.h"
#include "json.h"

@implementation FilePanel
+ (void)newFilePanel:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    NSString *ID = in[@"ID"];
    BOOL multipleSelection = [in[@"MultipleSelection"] boolValue];
    BOOL ignoreDirectories = [in[@"IgnoreDirectories"] boolValue];
    BOOL ignoreFiles = [in[@"IgnoreFiles"] boolValue];
    BOOL showHiddenFiles = [in[@"ShowHiddenFiles"] boolValue];
    NSArray<NSString *> *fileTypes = in[@"FileTypes"];

    NSOpenPanel *panel = [NSOpenPanel openPanel];
    panel.allowsMultipleSelection = multipleSelection;
    panel.canChooseDirectories = !ignoreDirectories;
    panel.canChooseFiles = !ignoreFiles;
    panel.showsHiddenFiles = showHiddenFiles;

    if (fileTypes != nil && fileTypes.count != 0) {
      panel.allowedFileTypes = fileTypes;
    }

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

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)newSaveFilePanel:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    NSString *ID = in[@"ID"];
    BOOL showHiddenFiles = [in[@"ShowHiddenFiles"] boolValue];
    NSArray<NSString *> *fileTypes = in[@"FileTypes"];

    NSSavePanel *panel = [NSSavePanel savePanel];
    panel.canCreateDirectories = YES;
    panel.showsHiddenFiles = showHiddenFiles;

    if (fileTypes != nil && fileTypes.count != 0) {
      panel.allowedFileTypes = fileTypes;
    }

    id onComplete = ^(NSInteger result) {
      NSString *filename = @"";

      if (result == NSModalResponseOK) {
        filename = panel.URL.absoluteString;
      }

      [driver.golang
          request:[NSString
                      stringWithFormat:@"/file/savepanel/select?id=%@", ID]
          payload:[JSONEncoder encodeString:filename]];
    };

    NSWindow *win = NSApp.keyWindow;
    if (win == nil) {
      [panel beginWithCompletionHandler:onComplete];
    } else {
      [panel beginSheetModalForWindow:win completionHandler:onComplete];
    }

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}
@end