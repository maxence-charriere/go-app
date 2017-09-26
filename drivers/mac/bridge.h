#ifndef bridge_h
#define bridge_h

#import <Cocoa/Cocoa.h>

// bridge_result is the struct that represent an Objective C request result.
typedef struct bridge_result {
  char *payload;
  char *err;
} bridge_result;

// OBJCHandler decribes the func that will handle requests to Objective C.
typedef bridge_result (^OBJCHandler)(NSURLComponents *, NSString *);

// make_bridge_result create a bridge result. If not nil, payload and err are a
// copy of the  NSString bytes. They should be free after use.
bridge_result make_bridge_result(NSString *payload, NSString *err);

// copyNSString copy the bytes of a NSString. Copyed bytes should be free after
// use.
char *copyNSString(NSString *str);

// macosRequest is the function to be called to handle a MacOS request.
bridge_result macosRequest(char *rawurl, char *payload);

// OBJCBridge is an objective-c bridge implementation.
@interface OBJCBridge : NSObject
@property NSMutableDictionary<NSString *, OBJCHandler> *handlers;

- (instancetype)init;
- (void)handle:(NSString *)path handler:(OBJCHandler)handler;
- (void)returnFor:(NSString *)id result:(bridge_result)res;
@end

// GoBridge is a golang bridge implementation.
@interface GoBridge : NSObject
- (void)request:(NSString *)path payload:(NSString *)payload;
- (NSString *)requestWithResult:(NSString *)path payload:(NSString *)payload;
@end

@interface NSURLComponents (Queryable)
- (NSString *)queryValue:(NSString *)name;
@end

#endif /* bridge_h */