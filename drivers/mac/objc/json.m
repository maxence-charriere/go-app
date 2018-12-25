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
@end

@implementation JSONDecoder
+ (id)decode:(NSString *)json {
  NSData *data = [json dataUsingEncoding:NSUTF8StringEncoding];
  return [NSJSONSerialization JSONObjectWithData:data
                                         options:NSJSONReadingMutableContainers
                                           error:nil];
}
@end
