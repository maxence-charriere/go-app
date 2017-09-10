#include "driver.h"
#include "_cgo_export.h"

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
  goRequest("/driver/run", "caca");
}

- (void)applicationDidBecomeActive:(NSNotification *)aNotification {
}

- (void)applicationDidResignActive:(NSNotification *)aNotification {
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender
                    hasVisibleWindows:(BOOL)flag {
  return YES;
}

- (void)application:(NSApplication *)sender
          openFiles:(NSArray<NSString *> *)filenames {
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
  NSURL *url =
      [NSURL URLWithString:[[event paramDescriptorForKeyword:keyDirectObject]
                               stringValue]];
}

- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication *)sender {
  return YES;
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
}

- (NSMenu *)applicationDockMenu:(NSApplication *)sender {
  return self.dock;
}
@end