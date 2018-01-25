#include "menu.h"
#include "driver.h"
#include "json.h"

@implementation MenuItem
- (instancetype)initFromTag:(NSDictionary *)tag {
  self = [super init];

  NSString *name = tag[@"name"];
  NSString *ID = tag[@"id"];
  NSString *compoID = tag[@"compo-id"];

  NSDictionary *attributes = tag[@"attributes"];
  NSString *label = @"";
  BOOL disabled = NO;
  BOOL separator = NO;
  NSString *selector = nil;
  NSString *onClick = nil;
  NSString *keys = nil;

  if (attributes != (id)[NSNull null]) {
    label = attributes[@"label"];
    label = label == nil ? @"" : label;

    disabled = attributes[@"disabled"] != nil ? YES : NO;
    separator = attributes[@"separator"] != nil ? YES : NO;
    selector = attributes[@"selector"];
    onClick = attributes[@"onclick"];
    keys = attributes[@"keys"];
  }

  self.ID = ID;
  self.compoID = compoID;
  self.title = label;
  self.enabled = !disabled;

  if (self.submenu != nil) {
    self.submenu.title = label;
  }

  if (separator && self.separator == nil) {
    self.separator = [NSMenuItem separatorItem];
    return self;
  }

  self.onClick = onClick;
  [self setupOnClick:selector];

  [self setupKeys:keys];
  return self;
}

- (BOOL)isSeparator {
  return self.separator != nil;
}

- (void)setupOnClick:(NSString *)selector {
  if (!self.enabled) {
    self.action = nil;
    return;
  }

  if (self.hasSubmenu) {
    self.action = @selector(submenuAction:);
    return;
  }

  if (selector != nil && selector.length > 0) {
    self.action = NSSelectorFromString(selector);
    return;
  }

  if (self.onClick == nil || self.onClick.length == 0) {
    return;
  }

  self.target = self;
  self.action = @selector(clicked:);
}

- (void)clicked:(id)sender {
  NSMutableDictionary<NSString *, id> *mapping =
      [[NSMutableDictionary alloc] init];
  mapping[@"compo-id"] = self.compoID;
  mapping[@"target"] = self.onClick;
  mapping[@"json-value"] = @"{}";

  Driver *driver = [Driver current];
  [driver.golang
      request:[NSString stringWithFormat:@"/menu/callback?id=%@", self.elemID]
      payload:[JSONEncoder encodeObject:mapping]];
}

- (void)setupKeys:(NSString *)keys {
  if (keys == nil || keys.length == 0) {
    return;
  }

  keys = [keys lowercaseString];

  NSArray *keyArray = [keys componentsSeparatedByString:@"+"];
  self.keyEquivalentModifierMask = 0;

  for (NSString *key in keyArray) {
    if ([key isEqual:@"cmd"] || [key isEqual:@"cmdorctrl"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagCommand;
    } else if ([key isEqual:@"ctrl"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagControl;
    } else if ([key isEqual:@"alt"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagOption;
    } else if ([key isEqual:@"shift"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagShift;
    } else if ([key isEqual:@"fn"]) {
      self.keyEquivalentModifierMask |= NSEventModifierFlagFunction;
    } else if ([key isEqual:@""]) {
      self.keyEquivalent = @"+";
    } else {
      self.keyEquivalent = key;
    }
  }
}
@end

@implementation MenuContainer
- (instancetype)init {
  self = [super init];
  self.children = [[NSMutableArray alloc] init];
  return self;
}

- (void)addChild:(MenuItem *)child {
  [self.children addObject:child];

  if ([child isSeparator]) {
    [self addItem:child.separator];
    child.menu = self;
    return;
  }

  [self addItem:child];
}

- (void)insertChild:(MenuItem *)child atIndex:(NSInteger)index {
  [self.children insertObject:child atIndex:index];

  if ([child isSeparator]) {
    [self insertItem:child.separator atIndex:index];
    child.menu = self;
  } else {
    [self insertItem:child atIndex:index];
  }
}

- (void)removeChildAtIndex:(NSInteger)index {
  [self.children removeObjectAtIndex:index];
  [self removeItemAtIndex:index];
}

@end

@implementation Menu
+ (bridge_result)newMenu:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];

  dispatch_async(dispatch_get_main_queue(), ^{
    Menu *menu = [[Menu alloc] init];
    menu.ID = ID;

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

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Menu *menu = driver.elements[ID];
    NSString *err = nil;

    @try {
      menu.root = [menu newContainer:content];
      menu.root.delegate = menu;
    } @catch (NSException *exception) {
      err = exception.reason;
    }

    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
  });
  return make_bridge_result(nil, nil);
}

+ (bridge_result)render:(NSURLComponents *)url payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];
  NSDictionary *tag = [JSONDecoder decodeObject:payload];
  NSString *tagID = tag[@"id"];

  dispatch_async(dispatch_get_main_queue(), ^{
                     // Driver *driver = [Driver current];
                     // Menu *menu = driver.elements[ID];

                     // id elem = [menu elementByID:tagID];
                     // MenuContainer *container = nil;
                     // MenuItem *item = nil;

                     // // Menu container.
                     // if ([elem isKindOfClass:[MenuContainer class]]) {
                     //   container = (MenuContainer *)elem;
                     //   return;
                     // }

                     // // Menu item.
                     // item = (MenuItem *)elem;
                 });
  return make_bridge_result(nil, nil);
}

+ (bridge_result)renderAttributes:(NSURLComponents *)url
                          payload:(NSString *)payload {
  NSString *ID = [url queryValue:@"id"];
  NSString *returnID = [url queryValue:@"return-id"];

  NSDictionary *tag = [JSONDecoder decodeObject:payload];
  NSString *tagID = tag[@"id"];

  NSDictionary *attributes = tag[@"attributes"];
  NSString *label = @"";
  BOOL separator = NO;

  if (attributes != (id)[NSNull null]) {
    label = attributes[@"label"];
    label = label == nil ? @"" : label;

    separator = attributes[@"separator"] != nil ? YES : NO;
  }

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Menu *menu = driver.elements[ID];
    NSString *err = nil;

    id elem = [menu elementByID:tagID];
    if (elem == nil) {
      err = [NSString stringWithFormat:@"no element with id %@", tagID];
      [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
      return;
    }

    // Menu container.
    // Should occur only for the root menu container.
    if ([elem isKindOfClass:[MenuContainer class]]) {
      MenuContainer *container = (MenuContainer *)elem;
      container.title = label;

      [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
      return;
    }

    // Menu item.
    MenuItem *item = (MenuItem *)elem;

    if ([item isSeparator] != separator) {
      NSString *err = nil;

      @try {
        MenuItem *newItem = [menu newItem:tag];
        MenuContainer *container = (MenuContainer *)item.menu;
        NSInteger index = [container.children indexOfObject:item];

        [container removeChildAtIndex:index];
        [container insertChild:newItem atIndex:index];
      } @catch (NSException *exception) {
        err = exception.reason;
      }

      [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
      return;
    }

    item = [item initFromTag:tag];
    [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
  });
  return make_bridge_result(nil, nil);
}

- (MenuContainer *)newContainer:(NSDictionary *)tag {
  NSString *name = tag[@"name"];
  NSString *ID = tag[@"id"];
  NSString *compoID = tag[@"compo-id"];

  NSDictionary *attributes = tag[@"attributes"];
  NSString *label = @"";
  NSArray *children = tag[@"children"];

  if (attributes != (id)[NSNull null]) {
    label = attributes[@"label"];
    label = label == nil ? @"" : label;
  }

  if (![name isEqual:@"menu"]) {
    @throw [NSException
        exceptionWithName:@"ErrMenuContainer"
                   reason:[NSString
                              stringWithFormat:
                                  @"cannot create a MenuContainer from a %@",
                                  name]
                 userInfo:nil];
  }

  MenuContainer *container = [[MenuContainer alloc] init];
  container.ID = ID;
  container.compoID = compoID;
  container.title = label;

  if (children != (id)[NSNull null]) {
    for (NSDictionary *child in children) {
      NSString *childName = child[@"name"];
      MenuItem *item = [self newItem:child];

      if ([childName isEqual:@"menu"]) {
        item.submenu = [self newContainer:child];
      }

      [container addChild:item];
    }
  }
  return container;
}

- (MenuItem *)newItem:(NSDictionary *)tag {
  MenuItem *item = nil;
  NSString *name = tag[@"name"];

  if (![name isEqual:@"menuitem"] && ![name isEqual:@"menu"]) {
    @throw [NSException
        exceptionWithName:@"ErrMenuItem"
                   reason:[NSString
                              stringWithFormat:
                                  @"cannot create a MenuItem from a %@", name]
                 userInfo:nil];
  }

  item = [[MenuItem alloc] initFromTag:tag];
  item.elemID = self.ID;
  return item;
}

- (id)elementByID:(NSString *)ID {
  return [self elementFromContainer:self.root ID:ID];
}

- (id)elementFromContainer:(MenuContainer *)container ID:(NSString *)ID {
  if ([container.ID isEqual:ID]) {
    return container;
  }

  for (MenuItem *child in container.children) {
    id elem = [self elementFromItem:child ID:ID];

    if (elem != nil) {
      return elem;
    }
  }

  return nil;
}

- (id)elementFromItem:(MenuItem *)item ID:(NSString *)ID {
  if ([item.ID isEqual:ID]) {
    return item;
  }

  if (item.submenu == nil) {
    return nil;
  }

  return [self elementFromContainer:(MenuContainer *)item.submenu ID:ID];
}

- (void)menuDidClose:(NSMenu *)menu {
  Driver *driver = [Driver current];

  [driver.golang
      request:[NSString stringWithFormat:@"/menu/close?id=%@", self.ID]
      payload:nil];

  [driver.elements removeObjectForKey:self.ID];
}
@end