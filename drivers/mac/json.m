#import "json.h"

@implementation JSONEncoder
+ (NSString *)encodeObject:(id)object {
  NSError *err = nil;
  NSData *jsonData =
      [NSJSONSerialization dataWithJSONObject:object
                                      options:NSJSONWritingPrettyPrinted
                                        error:&err];
  if (err != nil) {
    @throw [NSException exceptionWithName:@"encoding to JSON failed"
                                   reason:err.localizedDescription
                                 userInfo:nil];
  }

  NSString *jsonString =
      [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];

  return jsonString;
}

+ (NSString *)encodeString:(NSString *)s {
  return [NSString stringWithFormat:@"\"%@\"", s];
}

+ (NSString *)encodeBool:(BOOL)b {
  return b ? @"true" : @"false";
}
@end

@implementation JSONDecoder
+ (BOOL)decodeBool:(const char *)b {
  NSString *boolean = [NSString stringWithUTF8String:b];
  return [boolean isEqualToString:@"true"];
}
@end