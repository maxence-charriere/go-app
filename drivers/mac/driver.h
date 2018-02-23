#ifndef driver_h
#define driver_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>

@interface Driver
    : NSObject <NSApplicationDelegate, NSUserNotificationCenterDelegate>
@property OBJCBridge *objc;
@property GoBridge *golang;
@property NSMutableDictionary<NSString *, id> *elements;
@property NSMenu *dock;

+ (instancetype)current;
- (instancetype)init;
- (bridge_result)run:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)appName:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)resources:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)support:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)setContextMenu:(NSURLComponents *)url
                        payload:(NSString *)payload;
- (bridge_result)setMenuBar:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)setDock:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)setDockIcon:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)setDockBadge:(NSURLComponents *)url
                      payload:(NSString *)payload;
- (bridge_result)share:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)quit:(NSURLComponents *)url payload:(NSString *)payload;
@end

#endif /* driver_h */
