
#include "driver.h"
#include "json.h"
#include "panel.h"

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
      if (result != NSModalResponseOK) {
        return;
      }

      NSMutableArray<NSString *> *filenames = [[NSMutableArray alloc] init];
      for (NSURL *url in panel.URLs) {
        [filenames addObject:url.path];
      }

      NSDictionary *in = @{
        @"ID" : ID,
        @"Filenames" : filenames,
      };

      [driver.goRPC call:@"filePanels.OnSelect" withInput:in onUI:YES];
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
      if (result != NSModalResponseOK) {
        return;
      }

      NSDictionary *in = @{
        @"ID" : ID,
        @"Filename" : panel.URL.absoluteString,
      };

      [driver.goRPC call:@"saveFilePanels.OnSelect" withInput:in onUI:YES];
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