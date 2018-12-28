#ifndef color_h
#define color_h

#import <QuartzCore/QuartzCore.h>

@interface CIColor (MBCategory)
+ (CIColor *)colorWithHex:(UInt32)color;
+ (CIColor *)colorWithHexString:(NSString *)color;
@end

#endif /* color_h */
