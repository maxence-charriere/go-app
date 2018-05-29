#ifndef status_h
#define status_h

#import "bridge.h"
#import "menu.h"
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface StatusMenu : Menu
@property NSStatusItem *item;
+ (void) new:(NSDictionary *)in return:(NSString *)returnID;
+ (void)setText:(NSDictionary *)in return:(NSString *)returnID;
+ (void)setIcon:(NSDictionary *)in return:(NSString *)returnID;
@end
#endif /* status_h */