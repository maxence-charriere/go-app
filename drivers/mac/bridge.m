#include "bridge.h"
#include "_cgo_export.h"
#include "driver.h"

bridge_result make_bridge_result() {
  bridge_result res;
  res.payload = nil;
  res.err = nil;
  return res;
}

char *new_bridge_result_string(NSString *str) {
  int len = strlen(str.UTF8String) + 1;
  char *ret = calloc(len, sizeof(char));
  strcpy(ret, str.UTF8String);
  return ret;
}

bridge_result macosRequest(char *rawurl, char *cpayload) {
  NSString *urlstr = [NSString stringWithUTF8String:rawurl];
  free(rawurl);

  NSString *payload = nil;
  if (cpayload != nil) {
    payload = [NSString stringWithUTF8String:cpayload];
    free(cpayload);
  }

  NSURLComponents *url = [NSURLComponents componentsWithString:urlstr];
  bridge_result res = make_bridge_result();
  Driver *driver = [Driver current];

  OBJCHandler handler = driver.objc.handlers[url.path];
  if (handler == nil) {
    NSString *err = [NSString stringWithFormat:@"%@ is not handled", url.path];
    res.err = new_bridge_result_string(err);
    return res;
  }

  return handler(url, payload);
}

@implementation OBJCBridge
- (instancetype)init {
  self.handlers = [NSMutableDictionary dictionaryWithCapacity:128];
  return self;
}

- (void)handle:(NSString *)path handler:(OBJCHandler)handler {
  self.handlers[path] = handler;
}
@end

@implementation GoBridge
+ (void)request:(NSString *)path payload:(NSString *)payload {
  char *p = nil;
  if (payload != nil) {
    p = (char *)payload.UTF8String;
  }

  goRequest((char *)path.UTF8String, p);
}

+ (NSString *)requestWithResult:(NSString *)path payload:(NSString *)payload {
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

@implementation NSURLComponents (Queryable)
- (NSString *)queryValue:(NSString *)name {
  NSPredicate *predicate = [NSPredicate predicateWithFormat:@"name=%@", name];
  NSURLQueryItem *queryItem =
      [[self.queryItems filteredArrayUsingPredicate:predicate] firstObject];
  return queryItem.value;
}
@end
