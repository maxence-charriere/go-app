#include "driver.h"
#include "json.h"
#include "sandbox.h"
#include "window.h"

@implementation Driver
+ (instancetype)current {
  NSApplication *app = [NSApplication sharedApplication];

  if (app.delegate != nil) {
    return app.delegate;
  }

  Driver *driver = [[Driver alloc] init];
  app.delegate = driver;
  return driver;
}

- (instancetype)init {
  self.elements = [NSMutableDictionary dictionaryWithCapacity:256];
  self.objc = [[OBJCBridge alloc] init];

  [self.objc handle:@"/driver/run"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self run:url payload:payload];
            }];
  [self.objc handle:@"/driver/resources"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self resources:url payload:payload];
            }];
  [self.objc handle:@"/driver/support"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [self support:url payload:payload];
            }];

  [self.objc handle:@"/window/new"
            handler:^(NSURLComponents *url, NSString *payload) {
              return [Window newWindow:url payload:payload];
            }];

  self.dock = [[NSMenu alloc] initWithTitle:@""];
  return self;
}

- (bridge_result)run:(NSURLComponents *)url payload:(NSString *)payload {
  [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
  [NSApp run];
  return make_bridge_result();
}

- (bridge_result)resources:(NSURLComponents *)url payload:(NSString *)payload {
  NSBundle *mainBundle = [NSBundle mainBundle];
  NSString *resp = [JSONEncoder encodeString:mainBundle.resourcePath];

  bridge_result res = make_bridge_result();
  res.payload = new_bridge_result_string(resp);
  return res;
}

- (bridge_result)support:(NSURLComponents *)url payload:(NSString *)payload {
  bridge_result res = make_bridge_result();
  NSBundle *mainBundle = [NSBundle mainBundle];
  NSString *storagename = nil;

  if ([mainBundle isSandboxed]) {
    storagename = [JSONEncoder encodeString:NSHomeDirectory()];
    res.payload = new_bridge_result_string(storagename);
    return res;
  }

  NSArray *paths = NSSearchPathForDirectoriesInDomains(
      NSApplicationSupportDirectory, NSUserDomainMask, YES);
  NSString *applicationSupportDirectory = [paths firstObject];

  if (mainBundle.bundleIdentifier.length == 0) {
    storagename = [NSString
        stringWithFormat:@"%@/goapp/{appname}", applicationSupportDirectory];
  } else {
    storagename =
        [NSString stringWithFormat:@"%@/%@", applicationSupportDirectory,
                                   mainBundle.bundleIdentifier];
  }
  storagename = [JSONEncoder encodeString:storagename];
  res.payload = new_bridge_result_string(storagename);
  return res;
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
  [GoBridge request:@"/driver/run" payload:nil];
}

- (void)applicationDidBecomeActive:(NSNotification *)aNotification {
  [GoBridge request:@"/driver/focus" payload:nil];
}

- (void)applicationDidResignActive:(NSNotification *)aNotification {
  [GoBridge request:@"/driver/blur" payload:nil];
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender
                    hasVisibleWindows:(BOOL)flag {
  NSString *payload = flag ? @"true" : @"false";
  [GoBridge request:@"/driver/reopen" payload:payload];
  return YES;
}

- (void)application:(NSApplication *)sender
          openFiles:(NSArray<NSString *> *)filenames {
  NSString *payload = [JSONEncoder encodeObject:filenames];
  [GoBridge request:@"/driver/filesopen" payload:payload];
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
  [GoBridge request:@"/driver/urlopen" payload:payload];
}

- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication *)sender {
  NSString *res = [GoBridge requestWithResult:@"/driver/quit" payload:nil];
  return [JSONDecoder decodeBool:res];
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
  [GoBridge requestWithResult:@"/driver/exit" payload:nil];
}

- (NSMenu *)applicationDockMenu:(NSApplication *)sender {
  return self.dock;
}
@end