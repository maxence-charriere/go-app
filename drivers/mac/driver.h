#ifndef driver_h
#define driver_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>

@interface Driver
    : NSObject <NSApplicationDelegate, NSUserNotificationCenterDelegate>
@property MacRPC *macRPC;
@property GoRPC *goRPC;
@property NSMutableDictionary<NSString *, id> *elements;
@property NSMenu *dock;

+ (instancetype)current;
- (instancetype)init;
- (void)run:(id)in return:(NSString *)returnID;
- (void)bundle:(id)in return:(NSString *)returnID;
- (NSString *)support;
- (void)setContextMenu:(NSString *)menuID return:(NSString *)returnID;
- (void)setMenubar:(NSString *)menuID return:(NSString *)returnID;
- (void)share:(NSDictionary *)in return:(NSString *)returnID;
- (void)quit:(id)in return:(NSString *)returnID;
@end

#endif /* driver_h */
