#ifndef driver_h
#define driver_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>

#define defer(code)                                                            \
  dispatch_async(dispatch_get_main_queue(), ^{                                 \
                     code})

@interface Driver : NSObject <NSApplicationDelegate>
@property OBJCBridge *objc;
@property NSMenu *dock;

- (instancetype)init;
+ (instancetype)current;
- (bridge_result)run:(NSURL *)url payload:(NSString *)payload;
- (bridge_result)resources:(NSURL *)url payload:(NSString *)payload;
@end

#endif /* driver_h */