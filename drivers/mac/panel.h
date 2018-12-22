#ifndef file_h
#define file_h

#import <Cocoa/Cocoa.h>

@interface FilePanel : NSObject
+ (void)newFilePanel:(NSDictionary *)in return:(NSString *)returnID;
+ (void)newSaveFilePanel:(NSDictionary *)in return:(NSString *)returnID;
@end

#endif /* file_h */
