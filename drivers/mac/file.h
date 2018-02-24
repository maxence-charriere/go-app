#ifndef file_h
#define file_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>

@interface FilePanel : NSObject
+ (bridge_result)newFilePanel:(NSURLComponents *)url
                      payload:(NSString *)payload;
+ (bridge_result)newSaveFilePanel:(NSURLComponents *)url
                          payload:(NSString *)payload;
@end

#endif /* file_h */
