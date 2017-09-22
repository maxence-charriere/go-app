#ifndef driver_h
#define driver_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>

@interface Driver : NSObject <NSApplicationDelegate>
@property OBJCBridge *objc;
@property NSMutableDictionary<NSString *, id> *elements;
@property NSMenu *dock;

+ (instancetype)current;
- (instancetype)init;
- (bridge_result)run:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)resources:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)support:(NSURLComponents *)url payload:(NSString *)payload;
@end

#endif /* driver_h */