#include "bridge.h"
#include "_cgo_export.h"
#include "driver.h"
#include "json.h"

@implementation GoBridge
- (void)request:(NSString *)path payload:(NSString *)payload {
  char *p = nil;
  if (payload != nil) {
    p = (char *)payload.UTF8String;
  }

  goRequest((char *)path.UTF8String, p);
}

- (NSString *)requestWithResult:(NSString *)path payload:(NSString *)payload {
  char *p = nil;
  if (payload != nil) {
    p = (char *)payload.UTF8String;
  }

  char *cres = goRequestWithResult((char *)path.UTF8String, p);
  NSString *res = nil;
  if (cres != nil) {
    res = [NSString stringWithUTF8String:cres];
    free(cres);
  }
  return res;
}
@end

void macCall(char *rawCall) {
  NSDictionary *call =
      [JSONDecoder decode:[NSString stringWithUTF8String:rawCall]];

  NSString *method = call[@"Method"];
  id in = call[@"Input"];
  NSString *returnID = call[@"ReturnID"];

  @try {
    Driver *driver = [Driver current];
    MacRPCHandler handler = driver.macRPC.handlers[method];

    if (handler == nil) {
      [NSException raise:@"rpcNotHandled" format:@"%@ is not handled", method];
    }

    handler(in, returnID);
  } @catch (NSException *exception) {
    NSString *err = exception.reason;
    macCallReturn((char *)returnID.UTF8String, nil, (char *)err.UTF8String);
  }
}

void defer(NSString *returnID, dispatch_block_t block) {
  dispatch_async(dispatch_get_main_queue(), ^{
    @try {
      block();
    } @catch (NSException *exception) {
      NSString *err = exception.reason;
      macCallReturn((char *)returnID.UTF8String, nil, (char *)err.UTF8String);
    }
  });
}

@implementation MacRPC
- (instancetype)init {
  self = [super init];
  self.handlers = [NSMutableDictionary dictionaryWithCapacity:64];
  return self;
}

- (void)handle:(NSString *)method withHandler:(MacRPCHandler)handler {
  self.handlers[method] = handler;
}

- (void) return:(NSString *)returnID
     withOutput:(id)out
       andError:(NSString *)err {

  char *creturnID = returnID != nil ? (char *)returnID.UTF8String : nil;
  char *cout = out != nil ? (char *)[JSONEncoder encode:out].UTF8String : nil;
  char *cerr = err != nil ? (char *)err.UTF8String : nil;

  macCallReturn(creturnID, cout, cerr);
}
@end

@implementation GoRPC
- (id)call:(NSString *)method withInput:(id)in {
  NSDictionary *call = @{
    @"Method" : method,
    @"Input" : in,
  };
  NSString *callString = [JSONEncoder encode:call];

  char *cout = goCall((char *)callString.UTF8String);
  NSString *out = [NSString stringWithUTF8String:cout];
  free(cout);

  return [JSONDecoder decodeObject:out];
}
@end