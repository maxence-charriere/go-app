#include "image.h"

@implementation NSImage (ImageCategory)
+ (NSImage *)resizeImage:(NSImage *)sourceImage
       toPixelDimensions:(NSSize)newSize {
  if (!sourceImage.isValid)
    return nil;

  NSBitmapImageRep *rep = [[NSBitmapImageRep alloc]
      initWithBitmapDataPlanes:NULL
                    pixelsWide:newSize.width
                    pixelsHigh:newSize.height
                 bitsPerSample:8
               samplesPerPixel:4
                      hasAlpha:YES
                      isPlanar:NO
                colorSpaceName:NSCalibratedRGBColorSpace
                   bytesPerRow:0
                  bitsPerPixel:0];
  rep.size = newSize;

  [NSGraphicsContext saveGraphicsState];
  [NSGraphicsContext
      setCurrentContext:[NSGraphicsContext
                            graphicsContextWithBitmapImageRep:rep]];
  [sourceImage drawInRect:NSMakeRect(0, 0, newSize.width, newSize.height)
                 fromRect:NSZeroRect
                operation:NSCompositingOperationCopy
                 fraction:1.0];
  [NSGraphicsContext restoreGraphicsState];

  NSImage *newImage = [[NSImage alloc] initWithSize:newSize];
  [newImage addRepresentation:rep];
  return newImage;
}
@end
