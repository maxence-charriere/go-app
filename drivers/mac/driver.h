#ifndef driver_h
#define driver_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>

@interface Driver
    : NSObject <NSApplicationDelegate, NSUserNotificationCenterDelegate>
@property OBJCBridge *objc;
@property GoBridge *golang;
@property MacRPC *macRPC;
@property NSMutableDictionary<NSString *, id> *elements;
@property NSMenu *dock;

+ (instancetype)current;
- (instancetype)init;
- (void)run:(NSDictionary *)in return:(NSString *)returnID;
- (void)bundle:(NSDictionary *)in return:(NSString *)returnID;
- (NSString *)support;
- (void)setContextMenu:(NSDictionary *)in return:(NSString *)returnID;
- (void)setMenubar:(NSDictionary *)in return:(NSString *)returnID;
- (void)setDock:(NSDictionary *)in return:(NSString *)returnID;
- (void)setDockIcon:(NSDictionary *)in return:(NSString *)returnID;
- (void)setDockBadge:(NSDictionary *)in return:(NSString *)returnID;
- (void)share:(NSDictionary *)in return:(NSString *)returnID;
- (void)quit:(NSDictionary *)in return:(NSString *)returnID;
@end

#endif /* driver_h */
