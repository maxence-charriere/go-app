#include "dock.h"
#include "driver.h"
#include "image.h"
#include "json.h"

@implementation Dock
+ (void)setMenu:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];

    Menu *menu = driver.elements[ID];
    driver.dock = menu.root;

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)setIcon:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *icon = in[@"Icon"];

    if (icon.length != 0) {
      NSApp.applicationIconImage = [[NSImage alloc] initByReferencingFile:icon];
    } else {
      NSApp.applicationIconImage = nil;
    }

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)setBadge:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *badge = in[@"Badge"];

    [NSApp.dockTile setBadgeLabel:badge];
    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}
@end