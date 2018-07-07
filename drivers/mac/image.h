#ifndef image_h
#define image_h

#import <Cocoa/Cocoa.h>
#if MAC_OS_X_VERSION_MAX_ALLOWED < 101200
static const NSCompositingOperation NSCompositingOperationCopy = NSCompositeCopy;
#endif
 
@interface NSImage (ImageCategory)
+ (NSImage *)resizedImage:(NSImage *)sourceImage
        toPixelDimensions:(NSSize)newSize;
@end

#endif /* image_h */