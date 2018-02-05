
#include "notification.h"
#include "driver.h"
#include "json.h"

@implementation Notification
+ (bridge_result)newNotification:(NSURLComponents *)url
                         payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];
  NSString *returnID = [url queryValue:@"return-id"];

  NSDictionary *config = [JSONDecoder decodeObject:payload];
  NSString *title = config[@"title"];
  NSString *subtitle = config[@"subtitle"];
  NSString *text = config[@"text"];
  NSString *imageName = config[@"image-name"];
  BOOL reply = [config[@"reply"] boolValue];
  BOOL sound = [config[@"sound"] boolValue];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    NSString *err = nil;

    @try {
      NSUserNotification *notification = [[NSUserNotification alloc] init];
      notification.identifier = ID;
      notification.title = title;
      notification.subtitle = subtitle;
      notification.informativeText = text;
      notification.hasReplyButton = reply;
      notification.responsePlaceholder = @"murlok";

      if (imageName.length != 0) {
        notification.contentImage =
            [[NSImage alloc] initByReferencingFile:imageName];
      }

      if (sound) {
        notification.soundName = NSUserNotificationDefaultSoundName;
      }

      [[NSUserNotificationCenter defaultUserNotificationCenter]
          deliverNotification:notification];
    } @catch (NSException *exception) {
      err = exception.reason;
    }

    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
  });
  return make_bridge_result(nil, nil);
}
@end