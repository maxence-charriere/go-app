#ifndef bridge_h
#define bridge_h

#import <Cocoa/Cocoa.h>

// bridge_result is the struct that represent an Objective C request result.
typedef struct bridge_result {
  char *payload;
  char *err;
} bridge_result;

// OBJCHandler decribes the func that will handle requests to Objective C.
typedef bridge_result (^OBJCHandler)(NSURL *, NSString *);

bridge_result make_bridge_result();
char *new_bridge_result_string(NSString *str);
bridge_result macosRequest(char *rawurl, char *payload);

@interface OBJCBridge : NSObject
@property NSMutableDictionary<NSString *, OBJCHandler> *handlers;

- (instancetype)init;
- (void)handle:(NSString *)path handler:(OBJCHandler)handler;
@end

#endif /* bridge_h */