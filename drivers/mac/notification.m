
#include "notification.h"
#include "driver.h"
#include "json.h"

@implementation Notification
+ (void) new:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    NSString *ID = in[@"ID"];
    NSString *title = in[@"Title"];
    NSString *subtitle = in[@"Subtitle"];
    NSString *text = in[@"Text"];
    NSString *imageName = in[@"ImageName"];
    BOOL sound = [in[@"Sound"] boolValue];
    BOOL reply = [in[@"Reply"] boolValue];

    NSUserNotification *notification = [[NSUserNotification alloc] init];
    notification.identifier = ID;
    notification.title = title;
    notification.subtitle = subtitle;
    notification.informativeText = text;
    notification.hasReplyButton = reply;

    if (imageName.length != 0) {
      notification.contentImage =
          [[NSImage alloc] initByReferencingFile:imageName];
    }

    if (sound) {
      notification.soundName = NSUserNotificationDefaultSoundName;
    }

    [[NSUserNotificationCenter defaultUserNotificationCenter]
        deliverNotification:notification];

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}
@end