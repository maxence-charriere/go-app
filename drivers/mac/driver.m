#include "driver.h"
#include "_cgo_export.h"
#include "json.h"

void driver_run() {
  [NSApplication sharedApplication];
  [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

  DriverDelegate *delegate = [[DriverDelegate alloc] init];
  NSApp.delegate = delegate;

  [NSApp run];
}

@implementation DriverDelegate
- (instancetype)init {
  self.dock = [[NSMenu alloc] initWithTitle:@""];
  return self;
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