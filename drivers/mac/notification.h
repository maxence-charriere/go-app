#ifndef notification_h
#define notification_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>

@interface Notification : NSObject
+ (bridge_result)newNotification:(NSURLComponents *)url
                         payload:(NSString *)payload;
@end

#endif /* notification_h */
