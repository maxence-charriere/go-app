#ifndef json_h
#define json_h

#import <Cocoa/Cocoa.h>

@interface JSONEncoder : NSObject
+ (NSString *)encode:(id)object;
+ (NSString *)encodeObject:(id)object;
+ (NSString *)encodeString:(NSString *)s;
+ (NSString *)encodeBool:(BOOL)b;
@end

@interface JSONDecoder : NSObject
+ (id)decode:(NSString *)json;
+ (NSDictionary *)decodeObject:(NSString *)json;
+ (BOOL)decodeBool:(NSString *)json;
@end

#endif /* json_h */
