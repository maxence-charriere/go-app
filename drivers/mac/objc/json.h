#ifndef json_h
#define json_h

#import <Cocoa/Cocoa.h>

@interface JSONEncoder : NSObject
+ (NSString *)encode:(id)object;
@end

@interface JSONDecoder : NSObject
+ (id)decode:(NSString *)json;
@end

#endif /* json_h */
