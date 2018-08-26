#include "status.h"
#include "driver.h"
#include "image.h"
#include "json.h"

@implementation StatusMenu
+ (void) new:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];
    NSString *text = in[@"Text"];
    NSString *icon = in[@"Icon"];

    StatusMenu *menu = [[StatusMenu alloc] initWithID:ID];

    menu.item = [NSStatusBar.systemStatusBar
        statusItemWithLength:NSVariableStatusItemLength];
    menu.item.button.title = text;

    if (icon != nil) {
      CGFloat menuBarHeight = [[NSApp mainMenu] menuBarHeight];
      NSImage *img = [[NSImage alloc] initByReferencingFile:icon];
      menu.item.button.image =
          [NSImage resizeImage:img
              toPixelDimensions:NSMakeSize(menuBarHeight, menuBarHeight)];
      menu.item.button.imagePosition =
          text.length == 0 ? NSImageOnly : NSImageRight;
    }

    driver.elements[ID] = menu;
    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)setMenu:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];

    StatusMenu *menu = driver.elements[ID];
    menu.item.menu = menu.root;

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)setText:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];
    NSString *text = in[@"Text"];

    StatusMenu *menu = driver.elements[ID];
    menu.item.button.title = text;
    menu.item.button.imagePosition =
        text.length == 0 ? NSImageOnly : NSImageRight;

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)setIcon:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];
    NSString *icon = in[@"Icon"];
    StatusMenu *menu = driver.elements[ID];

    if (icon.length == 0) {
      menu.item.button.image = nil;
      menu.item.button.imagePosition = NSNoImage;
    } else {
      CGFloat menuBarHeight = [[NSApp mainMenu] menuBarHeight];
      NSImage *img = [[NSImage alloc] initByReferencingFile:icon];
      menu.item.button.image =
          [NSImage resizeImage:img
              toPixelDimensions:NSMakeSize(menuBarHeight, menuBarHeight)];

      NSString *text = menu.item.button.title;
      menu.item.button.imagePosition =
          text.length == 0 ? NSImageOnly : NSImageRight;
    }

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)close:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];
    StatusMenu *menu = driver.elements[ID];

    menu.item.menu = nil;
    [driver.elements removeObjectForKey:ID];

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}
@end