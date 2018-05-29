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

    StatusMenu *menu = [[StatusMenu alloc] init];

    NSStatusBar *bar = [NSStatusBar systemStatusBar];
    menu.item = [bar statusItemWithLength:NSVariableStatusItemLength];
    menu.item.button.title = text;

    if (icon != nil) {
      CGFloat menuBarHeight = [[NSApp mainMenu] menuBarHeight];
      NSImage *img = [[NSImage alloc] initByReferencingFile:icon];
      menu.item.button.image =
          [NSImage resizedImage:img
              toPixelDimensions:NSMakeSize(menuBarHeight, menuBarHeight)];
      menu.item.button.imagePosition =
          text.length == 0 ? NSImageOnly : NSImageRight;
    }

    menu.ID = ID;
    driver.elements[ID] = menu;

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)setText:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)setIcon:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}
@end