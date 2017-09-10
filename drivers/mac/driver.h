#ifndef driver_h
#define driver_h

#import <Cocoa/Cocoa.h>

#define defer(code)                                                            \
  dispatch_async(dispatch_get_main_queue(), ^{                                 \
                     code})

void driver_run();

@interface DriverDelegate : NSObject <NSApplicationDelegate>
@property NSMenu *dock;

- (instancetype)init;
@end

#endif /* driver_h */