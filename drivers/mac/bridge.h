#ifndef bridge_h
#define bridge_h

#define defer(code)                                                            \
  dispatch_async(dispatch_get_main_queue(), ^{                                 \
                     code})

#import <Cocoa/Cocoa.h>

// bridge_result is the struct that represent an Objective C request result.
typedef struct bridge_result {
  char *payload;
  char *err;
} bridge_result;

// OBJCHandler decribes the func that will handle requests to Objective C.
typedef bridge_result (^OBJCHandler)(NSURLComponents *, NSString *);

bridge_result make_bridge_result();
char *new_bridge_result_string(NSString *str);
bridge_result macosRequest(char *rawurl, char *payload);

@interface OBJCBridge : NSObject
@property NSMutableDictionary<NSString *, OBJCHandler> *handlers;

- (instancetype)init;
- (void)handle:(NSString *)path handler:(OBJCHandler)handler;
@end

@interface GoBridge : NSObject
+ (void)request:(NSString *)path payload:(NSString *)payload;
+ (NSString *)requestWithResult:(NSString *)path payload:(NSString *)payload;
@end

@interface NSURLComponents (Queryable)
- (NSString *)queryValue:(NSString *)name;
@end

#endif /* bridge_h */