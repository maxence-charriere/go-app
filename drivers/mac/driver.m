#include "driver.h"
#include "file.h"
#include "json.h"
#include "menu.h"
#include "notification.h"
#include "sandbox.h"
#include "window.h"

@implementation Driver
+ (instancetype)current {
  static Driver *driver = nil;

  @synchronized(self) {
    if (driver == nil) {
      driver = [[Driver alloc] init];
      NSApplication *app = [NSApplication sharedApplication];
      app.delegate = driver;
    }
  }
  return driver;
}

- (instancetype)init {
  self = [super init];

  self.elements = [NSMutableDictionary dictionaryWithCapacity:256];
  self.objc = [[OBJCBridge alloc] init];
  self.golang = [[GoBridge alloc] init];
  self.macRPC = [[MacRPC alloc] init];

  // Driver handlers.
  [self.macRPC handle:@"driver.Run"
          withHandler:^(NSDictionary *in, NSString *returnID) {
            return [self run:in return:returnID];
          }];
  [self.macRPC handle:@"driver.Bundle"
          withHandler:^(NSDictionary *in, NSString *returnID) {
            return [self bundle:in return:returnID];
          }];
  [self.macRPC handle:@"driver.SetContextMenu"
          withHandler:^(NSDictionary *in, NSString *returnID) {
            return [self setContextMenu:in return:returnID];
          }];
  [self.macRPC handle:@"driver.SetMenubar"
          withHandler:^(NSDictionary *in, NSString *returnID) {
            return [self setMenubar:in return:returnID];
          }];
  [self.macRPC handle:@"driver.SetDock"
          withHandler:^(NSDictionary *in, NSString *returnID) {
            return [self setDock:in return:returnID];
          }];
  [self.macRPC handle:@"driver.SetDockIcon"
          withHandler:^(NSDictionary *in, NSString *returnID) {
            return [self setDockIcon:in return:returnID];
          }];
  [self.macRPC handle:@"driver.SetDockBadge"
          withHandler:^(NSDictionary *in, NSString *returnID) {
            return [self setDockBadge:in return:returnID];
          }];
  [self.macRPC handle:@"driver.Share"
          withHandler:^(NSDictionary *in, NSString *returnID) {
            return [self share:in return:returnID];
          }];
  [self.macRPC handle:@"driver.Quit"
          withHandler:^(NSDictionary *in, NSString *returnID) {
            return [self quit:in return:returnID];
          }];

  // Window handlers.
  [self.objc handle:@"/window/new"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window newWindow:url payload:payload];
            }];
  [self.objc handle:@"/window/load"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window load:url payload:payload];
            }];
  [self.objc handle:@"/window/render"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window render:url payload:payload];
            }];
  [self.objc handle:@"/window/render/attributes"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window renderAttributes:url payload:payload];
            }];
  [self.objc handle:@"/window/position"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window position:url payload:payload];
            }];
  [self.objc handle:@"/window/move"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window move:url payload:payload];
            }];
  [self.objc handle:@"/window/center"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window center:url payload:payload];
            }];
  [self.objc handle:@"/window/size"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window size:url payload:payload];
            }];
  [self.objc handle:@"/window/resize"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window resize:url payload:payload];
            }];
  [self.objc handle:@"/window/focus"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window focus:url payload:payload];
            }];
  [self.objc handle:@"/window/togglefullscreen"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window toggleFullScreen:url payload:payload];
            }];
  [self.objc handle:@"/window/toggleminimize"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window toggleMinimize:url payload:payload];
            }];
  [self.objc handle:@"/window/close"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window close:url payload:payload];
            }];

  // Menu handlers.
  [self.objc handle:@"/menu/new"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Menu newMenu:url payload:payload];
            }];
  [self.objc handle:@"/menu/load"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Menu load:url payload:payload];
            }];
  [self.objc handle:@"/menu/render"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Menu render:url payload:payload];
            }];
  [self.objc handle:@"/menu/render/attributes"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Menu renderAttributes:url payload:payload];
            }];
  [self.objc handle:@"/menu/delete"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Menu delete:url payload:payload];
            }];

  // File panel handlers.
  [self.objc handle:@"/file/panel/new"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [FilePanel newFilePanel:url payload:payload];
            }];
  [self.objc handle:@"/file/savepanel/new"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [FilePanel newSaveFilePanel:url payload:payload];
            }];

  // Notification handlers.
  [self.objc handle:@"/notification/new"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Notification newNotification:url payload:payload];
            }];

  // Notifications.
  NSUserNotificationCenter *userNotificationCenter =
      [NSUserNotificationCenter defaultUserNotificationCenter];
  userNotificationCenter.delegate = self;

  return self;
}

- (void)run:(NSDictionary *)in return:(NSString *)returnID {
  [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
  [NSApp activateIgnoringOtherApps:YES];
  [NSApp run];

  [self.macRPC return:returnID withOutput:nil andError:nil];
}

- (void)bundle:(NSDictionary *)in return:(NSString *)returnID {
  NSBundle *mainBundle = [NSBundle mainBundle];

  NSMutableDictionary *out = [[NSMutableDictionary alloc] init];
  out[@"AppName"] = mainBundle.infoDictionary[@"CFBundleName"];
  out[@"Resources"] = mainBundle.resourcePath;
  out[@"Support"] = [self support];

  [self.macRPC return:returnID withOutput:out andError:nil];
}

- (NSString *)support {
  NSBundle *mainBundle = [NSBundle mainBundle];

  if ([mainBundle isSandboxed]) {
    return NSHomeDirectory();
  }

  NSArray *paths = NSSearchPathForDirectoriesInDomains(
      NSApplicationSupportDirectory, NSUserDomainMask, YES);
  NSString *applicationSupportDirectory = [paths firstObject];

  if (mainBundle.bundleIdentifier.length == 0) {
    return [NSString
        stringWithFormat:@"%@/goapp/{appname}", applicationSupportDirectory];
  }
  return [NSString stringWithFormat:@"%@/%@", applicationSupportDirectory,
                                    mainBundle.bundleIdentifier];
}

- (void)setContextMenu:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Menu *menu = self.elements[in[@"MenuID"]];
    NSWindow *win = NSApp.keyWindow;

    if (win == nil) {
      [self.macRPC return:returnID withOutput:nil
          andError:@"no window to host the context menu"];
      return;
    }

    [menu.root popUpMenuPositioningItem:menu.root.itemArray[0]
                             atLocation:[win mouseLocationOutsideOfEventStream]
                                 inView:win.contentView];
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)setMenubar:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Menu *menu = self.elements[in[@"MenuID"]];
    NSApp.mainMenu = menu.root;
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)setDock:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Menu *menu = self.elements[in[@"MenuID"]];
    self.dock = menu.root;
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)setDockIcon:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    NSString *icon = in[@"Icon"];

    if (icon.length != 0) {
      NSApp.applicationIconImage = [[NSImage alloc] initByReferencingFile:icon];
    } else {
      NSApp.applicationIconImage = nil;
    }

    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)setDockBadge:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    [NSApp.dockTile setBadgeLabel:in[@"Badge"]];
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)share:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    NSWindow *rawWindow = NSApp.keyWindow;
    if (rawWindow == nil) {
      [NSException raise:@"NoKeyWindowExeption"
                  format:@"no window to host the share menu"];
    }
    Window *win = (Window *)rawWindow.windowController;

    id share = in[@"Share"];
    if ([in[@"Type"] isEqual:@"url"]) {
      [NSURL URLWithString:share];
    }

    NSPoint pos = [win.window mouseLocationOutsideOfEventStream];
    pos = [win.webview convertPoint:pos fromView:rawWindow.contentView];
    NSRect rect = NSMakeRect(pos.x, pos.y, 1, 1);

    NSSharingServicePicker *picker =
        [[NSSharingServicePicker alloc] initWithItems:@[ share ]];
    [picker showRelativeToRect:rect
                        ofView:win.webview
                 preferredEdge:NSMinYEdge];

    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)quit:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    [NSApp terminate:self];
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
  [self.golang request:@"/driver/run" payload:nil];
}

- (void)applicationDidBecomeActive:(NSNotification *)aNotification {
  [self.golang request:@"/driver/focus" payload:nil];
}

- (void)applicationDidResignActive:(NSNotification *)aNotification {
  [self.golang request:@"/driver/blur" payload:nil];
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender
                    hasVisibleWindows:(BOOL)flag {
  NSString *payload = flag ? @"true" : @"false";
  [self.golang request:@"/driver/reopen" payload:payload];
  return YES;
}

- (void)application:(NSApplication *)sender
          openFiles:(NSArray<NSString *> *)filenames {
  NSString *payload = [JSONEncoder encodeObject:filenames];
  [self.golang request:@"/driver/filesopen" payload:payload];
}

- (void)applicationWillFinishLaunching:(NSNotification *)aNotification {
  NSAppleEventManager *appleEventManager =
      [NSAppleEventManager sharedAppleEventManager];
  [appleEventManager
      setEventHandler:self
          andSelector:@selector(handleGetURLEvent:withReplyEvent:)
        forEventClass:kInternetEventClass
           andEventID:kAEGetURL];
}

- (void)handleGetURLEvent:(NSAppleEventDescriptor *)event
           withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
  NSString *rawurl =
      [event paramDescriptorForKeyword:keyDirectObject].stringValue;
  NSString *payload = [JSONEncoder encodeString:rawurl];
  [self.golang request:@"/driver/urlopen" payload:payload];
}

- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication *)sender {
  NSString *res = [self.golang requestWithResult:@"/driver/quit" payload:nil];
  return [JSONDecoder decodeBool:res];
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
  [self.golang requestWithResult:@"/driver/exit" payload:nil];
}

- (NSMenu *)applicationDockMenu:(NSApplication *)sender {
  return self.dock;
}

- (BOOL)userNotificationCenter:(NSUserNotificationCenter *)center
     shouldPresentNotification:(NSUserNotification *)notification {
  return YES;
}

- (void)userNotificationCenter:(NSUserNotificationCenter *)center
       didActivateNotification:(NSUserNotification *)notification {
  [self.golang request:[NSString stringWithFormat:@"/notification/reply?id=%@",
                                                  notification.identifier]
               payload:[JSONEncoder encodeString:notification.response.string]];
}
@end
