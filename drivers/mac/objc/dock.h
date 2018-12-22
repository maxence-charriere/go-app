#ifndef dock_h
#define dock_h

#import "menu.h"
#import <Cocoa/Cocoa.h>

@interface Dock : Menu
+ (void)setMenu:(NSDictionary *)in return:(NSString *)returnID;
+ (void)setIcon:(NSDictionary *)in return:(NSString *)returnID;
+ (void)setBadge:(NSDictionary *)in return:(NSString *)returnID;
@end
#endif /* dock_h */