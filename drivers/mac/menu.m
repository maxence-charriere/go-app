#include "menu.h"
#include "driver.h"
#include "json.h"

@implementation MenuItem
- (instancetype)initWithMenuID:(NSString *)menuID andTag:(NSDictionary *)tag {
  self = [super init];

  NSString *name = tag[@"name"];
  if (![name isEqual:@"menuitem"] && ![name isEqual:@"menu"]) {
    NSString *err =
        [NSString stringWithFormat:@"cannot create a MenuItem from a %@", name];
    @throw
        [NSException exceptionWithName:@"ErrMenuItem" reason:err userInfo:nil];
  }

  self.elemID = menuID;
  self.ID = tag[@"id"];
  self.compoID = tag[@"compo-id"];
  self.separator = nil;
  self.title = @"";
  self.enabled = YES;
  self.onClick = nil;
  self.selector = nil;
  self.keys = nil;

  NSDictionary *attributes = tag[@"attributes"];
  if (attributes != nil) {
    BOOL separator = attributes[@"separator"] != nil ? YES : NO;
    if (separator) {
      self.separator = [NSMenuItem separatorItem];
      return self;
    }

    NSString *label = attributes[@"label"];
    if (label != nil) {
      self.title = label;
      if (self.submenu != nil) {
        self.submenu.title = label;
      }
    }

    if (attributes[@"disabled"] != nil) {
      self.enabled = false;
    }

    self.onClick = attributes[@"onclick"];
    self.selector = attributes[@"selector"];
    self.keys = attributes[@"keys"];
  }

  [self setupOnClick];
  [self setupKeys];
  return self;
}

- (BOOL)isSeparator {
  return self.separator != nil;
}

- (void)setupOnClick {
  if (!self.enabled) {
    self.action = nil;
    return;
  }

  if (self.hasSubmenu) {
    self.action = @selector(submenuAction:);
    return;
  }

  if (self.selector != nil && self.selector.length > 0) {
    self.action = NSSelectorFromString(self.selector);
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

- (void)setupKeys {
  if (self.keys == nil || self.keys.length == 0) {
    return;
  }

  self.keyEquivalentModifierMask = 0;
  self.keys = [self.keys lowercaseString];

  NSArray *keys = [self.keys componentsSeparatedByString:@"+"];
  for (NSString *key in keys) {
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
- (instancetype)initWithMenuID:(NSString *)menuID andTag:(NSDictionary *)tag {
  self = [super init];

  NSString *name = tag[@"name"];
  if (![name isEqual:@"menu"]) {
    NSString *err = [NSString
        stringWithFormat:@"cannot create a MenuContainer from a %@", name];
    @throw [NSException exceptionWithName:@"ErrMenuContainer"
                                   reason:err
                                 userInfo:nil];
  }

  self.elemID = menuID;
  self.ID = tag[@"id"];
  self.compoID = tag[@"compo-id"];
  self.title = @"";
  self.children = [[NSMutableArray alloc] init];
  [self removeAllItems];

  NSDictionary *attributes = tag[@"attributes"];
  if (attributes != nil) {
    NSString *label = attributes[@"label"];
    if (label != nil) {
      self.title = label;
    }
  }

  NSArray *children = tag[@"children"];
  if (children == nil) {
    return self;
  }

  for (NSDictionary *child in children) {
    MenuItem *childItem = [[MenuItem alloc] initWithMenuID:menuID andTag:child];

    NSString *childName = child[@"name"];
    if ([childName isEqual:@"menu"]) {
      childItem.submenu =
          [[MenuContainer alloc] initWithMenuID:menuID andTag:child];
    }

    [self addChild:childItem];
  }
  return self;
}

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
  NSDictionary *tag = [JSONDecoder decodeObject:payload];

  dispatch_async(dispatch_get_main_queue(), ^{
    Driver *driver = [Driver current];
    Menu *menu = driver.elements[ID];
    NSString *err = nil;

    @try {
      menu.root = [[MenuContainer alloc] initWithMenuID:ID andTag:tag];
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
  NSString *returnID = [url queryValue:@"return-id"];

  NSDictionary *tag = [JSONDecoder decodeObject:payload];
  NSString *name = tag[@"name"];
  NSString *tagID = tag[@"id"];

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

    @try {
      // Menu container.
      // Should occur only for the root menu container.
      if ([elem isKindOfClass:[MenuContainer class]]) {
        if (![name isEqual:@"menu"]) {
          err =
              [NSString stringWithFormat:@"root tag must be a menu: %@", name];
          [driver.objc asyncReturn:returnID
                            result:make_bridge_result(nil, err)];
          return;
        }

        MenuContainer *container = (MenuContainer *)elem;
        container = [container initWithMenuID:menu.ID andTag:tag];
        [driver.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
        return;
      }

      MenuItem *item = (MenuItem *)elem;
      MenuContainer *container = (MenuContainer *)item.menu;
      NSInteger index = [container.children indexOfObject:item];
      [container removeChildAtIndex:index];

      MenuItem *newItem = [[MenuItem alloc] initWithMenuID:menu.ID andTag:tag];

      if ([name isEqual:@"menu"]) {
        newItem.submenu =
            [[MenuContainer alloc] initWithMenuID:menu.ID andTag:tag];
      }

      [container insertChild:newItem atIndex:index];

      [driver.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
    } @catch (NSException *exception) {
      err = exception.reason;
      [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
    }
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

  BOOL separator = NO;
  if (attributes != nil) {
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

    @try {
      // Menu container.
      // Should occur only for the root menu container.
      if ([elem isKindOfClass:[MenuContainer class]]) {
        MenuContainer *container = (MenuContainer *)elem;
        container = [container initWithMenuID:menu.ID andTag:tag];
        [driver.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
        return;
      }

      // Menu item.
      MenuItem *item = (MenuItem *)elem;

      if ([item isSeparator] != separator) {
        MenuItem *newItem =
            [[MenuItem alloc] initWithMenuID:menu.ID andTag:tag];
        MenuContainer *container = (MenuContainer *)item.menu;
        NSInteger index = [container.children indexOfObject:item];

        [container removeChildAtIndex:index];
        [container insertChild:newItem atIndex:index];
      }

      item = [item initWithMenuID:menu.ID andTag:tag];
      [driver.objc asyncReturn:returnID result:make_bridge_result(nil, nil)];
    } @catch (NSException *exception) {
      err = exception.reason;
      [driver.objc asyncReturn:returnID result:make_bridge_result(nil, err)];
    }
  });
  return make_bridge_result(nil, nil);
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