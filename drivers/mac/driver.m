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

  // Drivers handlers.
  [self.objc handle:@"/driver/run"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self run:url payload:payload];
            }];
  [self.objc handle:@"/driver/appname"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self appName:url payload:payload];
            }];
  [self.objc handle:@"/driver/resources"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self resources:url payload:payload];
            }];
  [self.objc handle:@"/driver/support"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self support:url payload:payload];
            }];
  [self.objc handle:@"/driver/contextmenu/set"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self setContextMenu:url payload:payload];
            }];
  [self.objc handle:@"/driver/menubar/set"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self setMenuBar:url payload:payload];
            }];
  [self.objc handle:@"/driver/dock/set"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self setDock:url payload:payload];
            }];
  [self.objc handle:@"/driver/dock/icon"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self setDockIcon:url payload:payload];
            }];
  [self.objc handle:@"/driver/dock/badge"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self setDockBadge:url payload:payload];
            }];
  [self.objc handle:@"/driver/share"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self share:url payload:payload];
            }];
  [self.objc handle:@"/driver/quit"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self quit:url payload:payload];
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

- (bridge_result)run:(NSURLComponents *)url payload:(NSString *)payload {
  [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
  [NSApp activateIgnoringOtherApps:YES];
  [NSApp run];
  return make_bridge_result(nil, nil);
}

- (bridge_result)appName:(NSURLComponents *)url payload:(NSString *)payload {
  NSBundle *mainBundle = [NSBundle mainBundle];
  NSString *res =
      [JSONEncoder encodeString:mainBundle.infoDictionary[@"CFBundleName"]];
  return make_bridge_result(res, nil);
}

- (bridge_result)resources:(NSURLComponents *)url payload:(NSString *)payload {
  NSBundle *mainBundle = [NSBundle mainBundle];
  NSString *res = [JSONEncoder encodeString:mainBundle.resourcePath];
  return make_bridge_result(res, nil);
}

- (bridge_result)support:(NSURLComponents *)url payload:(NSString *)payload {
  NSBundle *mainBundle = [NSBundle mainBundle];
  NSString *dirname = nil;

  if ([mainBundle isSandboxed]) {
    dirname = [JSONEncoder encodeString:NSHomeDirectory()];
    return make_bridge_result(dirname, nil);
  }

  NSArray *paths = NSSearchPathForDirectoriesInDomains(
      NSApplicationSupportDirectory, NSUserDomainMask, YES);
  NSString *applicationSupportDirectory = [paths firstObject];

  if (mainBundle.bundleIdentifier.length == 0) {
    dirname = [NSString
        stringWithFormat:@"%@/goapp/{appname}", applicationSupportDirectory];
  } else {
    dirname = [NSString stringWithFormat:@"%@/%@", applicationSupportDirectory,
                                         mainBundle.bundleIdentifier];
  }
  dirname = [JSONEncoder encodeString:dirname];
  return make_bridge_result(dirname, nil);
}

- (bridge_result)setContextMenu:(NSURLComponents *)url
                        payload:(NSString *)payload {
  NSString *menuID = [url queryValue:@"menu-id"];
  NSString *returnID = [url queryValue:@"return-id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Menu *menu = self.elements[menuID];
    NSWindow *win = NSApp.keyWindow;

    if (win == nil) {
      [self.objc asyncReturn:returnID
                      result:make_bridge_result(
                                 nil, @"no window to host the context menu")];
      return;
    }

    [menu.root popUpMenuPositioningItem:menu.root.itemArray[0]
                             atLocation:[win mouseLocationOutsideOfEventStream]
                                 inView:win.contentView];
    [self.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
  });
  return make_bridge_result(nil, nil);
}

- (bridge_result)setMenuBar:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *menuID = [url queryValue:@"menu-id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Menu *menu = self.elements[menuID];
    NSApp.mainMenu = menu.root;
  });
  return make_bridge_result(nil, nil);
}

- (bridge_result)setDock:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *menuID = [url queryValue:@"menu-id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Menu *menu = self.elements[menuID];
    self.dock = menu.root;
  });
  return make_bridge_result(nil, nil);
}

- (bridge_result)setDockIcon:(NSURLComponents *)url
                     payload:(NSString *)payload {
  NSString *returnID = [url queryValue:@"return-id"];
  NSDictionary *icon = [JSONDecoder decodeObject:payload];
  NSString *path = icon[@"path"];

  dispatch_async(dispatch_get_main_queue(), ^{
    NSString *err = nil;

    @try {
      if (path.length != 0) {
        NSApp.applicationIconImage =
            [[NSImage alloc] initByReferencingFile:path];
      } else {
        NSApp.applicationIconImage = nil;
      }
      [self.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
    } @catch (NSException *exception) {
      err = exception.reason;
      [self.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
    }

  });
  return make_bridge_result(nil, nil);
}

- (bridge_result)setDockBadge:(NSURLComponents *)url
                      payload:(NSString *)payload {
  NSDictionary *badge = [JSONDecoder decodeObject:payload];
  NSString *msg = badge[@"message"];

  dispatch_async(dispatch_get_main_queue(), ^{
    [NSApp.dockTile setBadgeLabel:msg];
  });
  return make_bridge_result(nil, nil);
}

- (bridge_result)share:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *returnID = [url queryValue:@"return-id"];

  NSDictionary *share = [JSONDecoder decodeObject:payload];
  NSString *value = share[@"value"];
  NSString *type = share[@"type"];

  dispatch_async(dispatch_get_main_queue(), ^{
    NSString *err = nil;

    @try {
      NSWindow *rawWindow = NSApp.keyWindow;
      if (rawWindow == nil) {
        @throw
            [NSException exceptionWithName:@"NoKeyWindowExeption"
                                    reason:@"no window to host the share menu"
                                  userInfo:nil];
      }

      id valueToShare = value;
      if ([type isEqual:@"url"]) {
        valueToShare = [NSURL URLWithString:value];
      }

      Window *win = (Window *)rawWindow.windowController;
      NSPoint pos = [win.window mouseLocationOutsideOfEventStream];
      pos = [win.webview convertPoint:pos fromView:rawWindow.contentView];
      NSRect rect = NSMakeRect(pos.x, pos.y, 1, 1);

      NSSharingServicePicker *picker =
          [[NSSharingServicePicker alloc] initWithItems:@[ valueToShare ]];
      [picker showRelativeToRect:rect
                          ofView:win.webview
                   preferredEdge:NSMinYEdge];
    } @catch (NSException *exception) {
      err = exception.reason;
    }

    [self.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
  });
  return make_bridge_result(nil, nil);
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

- (bridge_result)quit:(NSURLComponents *)url payload:(NSString *)payload {
  dispatch_async(dispatch_get_main_queue(), ^{
    [NSApp terminate:self];
  });
  return make_bridge_result(nil, nil);
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
