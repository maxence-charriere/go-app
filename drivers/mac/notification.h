#ifndef notification_h
#define notification_h

#import <Cocoa/Cocoa.h>

@interface Notification : NSObject
+ (void) new:(NSDictionary *)in return:(NSString *)returnID;
@end

#endif /* notification_h */
