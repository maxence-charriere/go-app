#include "menu.h"
#include "driver.h"
#include "json.h"

@implementation Menu
+ (bridge_result)newMenu:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Menu *menu = [[Menu alloc] init];

    // Registering menu.
    Driver *driver = [Driver current];
    driver.elements[ID] = menu;
  });
  return make_bridge_result(nil, nil);
}

+ (bridge_result)load:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];
  NSString *returnID = [url queryValue:@"return-id"];
  NSDictionary *content = [JSONDecoder decodeObject:payload];

  NSLog(@"load payload: %@", content);

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Menu *menu = driver.elements[ID];
    NSString *err = nil;

    @try {
      menu.root = [menu newContainer:content];
    } @catch (NSException *exception) {
      err = exception.reason;
    }

    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
  });
  return make_bridge_result(nil, nil);
}

- (MenuContainer *)newContainer:(NSDictionary *)map {
  NSString *name = map[@"name"];
  NSString *label = nil;
  NSDictionary *attributes = map[@"attributes"];
  NSArray *children = map[@"children"];

  if (attributes != (id)[NSNull null]) {
    label = attributes[@"label"];
  }

  if (![name isEqual:@"menu"]) {
    @throw [NSException
        exceptionWithName:@"ErrMenuContainer"
                   reason:[NSString
                              stringWithFormat:
                                  @"cannot create a NSMenu from a %@", name]
                 userInfo:nil];
  }

  label = label == nil ? @"" : label;
  MenuContainer *container = [[MenuContainer alloc] initWithTitle:label];

  if (children != (id)[NSNull null]) {
    for (NSDictionary *child in children) {
      NSString *childName = child[@"name"];
      MenuItem *item = nil;

      if ([childName isEqual:@"menu"]) {
        item = [[MenuItem alloc] init];
        item.submenu = [self newContainer:child];
      } else {
        item = [self newItem:child];
      }

      [container addItem:item];
    }
  }
  return container;
}

- (MenuItem *)newItem:(NSDictionary *)map {
  MenuItem *item = [[MenuItem alloc] init];
  item.title = @"an item";
  return item;
}
@end

@implementation MenuContainer
@end

@implementation MenuItem
@end