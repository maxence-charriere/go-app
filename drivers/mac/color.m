#include "color.h"

@implementation CIColor (MBCategory)
+ (CIColor *)colorWithHexString:(NSString *)color {
  const char *str = [color cStringUsingEncoding:NSASCIIStringEncoding];
  return [CIColor colorWithHex:strtol(str + 1, NULL, 16)];
}

+ (CIColor *)colorWithHex:(UInt32)color {
  unsigned char b = color & 0xFF;
  unsigned char g = (color >> 8) & 0xFF;
  unsigned char r = (color >> 16) & 0xFF;

  return [CIColor colorWithRed:(float)r / 255.0f
                         green:(float)g / 255.0f
                          blue:(float)b / 255.0f
                         alpha:1];
}
@end
