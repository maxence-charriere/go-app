#import "json.h"

@implementation JSONEncoder
+ (NSString *)encode:(id)object {
  NSError *err = nil;

  NSData *json =
      [NSJSONSerialization dataWithJSONObject:object
                                      options:NSJSONWritingPrettyPrinted
                                        error:&err];

  if (err != nil) {
    [NSException raise:@"ErrJSONEncode" format:@"%@", err.localizedDescription];
  }

  return [[NSString alloc] initWithData:json encoding:NSUTF8StringEncoding];
}

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
+ (id)decode:(NSString *)json {
  NSData *data = [json dataUsingEncoding:NSUTF8StringEncoding];
  return [NSJSONSerialization JSONObjectWithData:data
                                         options:NSJSONReadingMutableContainers
                                           error:nil];
}

+ (NSDictionary *)decodeObject:(NSString *)json {
  NSData *data = [json dataUsingEncoding:NSUTF8StringEncoding];
  return [NSJSONSerialization JSONObjectWithData:data
                                         options:NSJSONReadingMutableContainers
                                           error:nil];
}

+ (BOOL)decodeBool:(NSString *)json {
  return [json isEqualToString:@"true"];
}
@end
