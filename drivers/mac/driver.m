#include "driver.h"
#include "_cgo_export.h"
#include "json.h"

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
  self.objc = [[OBJCBridge alloc] init];

  [self.objc handle:@"/driver/run"
            handler:^(NSURL *url, NSString *payload) {
              return [self run:url payload:payload];
            }];
  [self.objc handle:@"/driver/resources"
            handler:^(NSURL *url, NSString *payload) {
              return [self resources:url payload:payload];
            }];

  self.dock = [[NSMenu alloc] initWithTitle:@""];
  return self;
}

- (bridge_result)run:(NSURL *)url payload:(NSString *)payload {
  [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
  [NSApp run];
  return make_bridge_result();
}

- (bridge_result)resources:(NSURL *)url payload:(NSString *)payload {
  NSBundle *mainBundle = [NSBundle mainBundle];
  NSString *resp = [JSONEncoder encodeString:mainBundle.resourcePath];

  bridge_result res = make_bridge_result();
  res.payload = new_bridge_result_string(resp);
  return res;
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
  goRequest("/driver/run", nil);
}

- (void)applicationDidBecomeActive:(NSNotification *)aNotification {
  goRequest("/driver/focus", nil);
}

- (void)applicationDidResignActive:(NSNotification *)aNotification {
  goRequest("/driver/blur", nil);
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender
                    hasVisibleWindows:(BOOL)flag {
  NSString *payload = flag ? @"true" : @"false";
  goRequest("/driver/reopen", (char *)payload.UTF8String);
  return YES;
}

- (void)application:(NSApplication *)sender
          openFiles:(NSArray<NSString *> *)filenames {
  NSString *payload = [JSONEncoder encodeObject:filenames];
  goRequest("/driver/filesopen", (char *)payload.UTF8String);
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
  goRequest("/driver/urlopen", (char *)payload.UTF8String);
}

- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication *)sender {
  char *res = goRequestWithResult("/driver/quit", nil);
  BOOL shouldTerminate = [JSONDecoder decodeBool:res];
  free(res);
  return shouldTerminate;
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
  char *res = goRequestWithResult("/driver/exit", nil);
  free(res);
}

- (NSMenu *)applicationDockMenu:(NSApplication *)sender {
  return self.dock;
}
@end