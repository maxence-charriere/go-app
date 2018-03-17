#include "menu.h"
#include "driver.h"
#include "json.h"

@implementation MenuItem
- (instancetype)initWithMenuID:(NSString *)menuID andTag:(NSDictionary *)tag {
  self = [super init];

  NSString *name = tag[@"Name"];
  if (![name isEqual:@"menuitem"] && ![name isEqual:@"menu"]) {
    NSString *err = [NSString
        stringWithFormat:@"cannot create a MenuItem from a %@ tag", name];
    @throw
        [NSException exceptionWithName:@"ErrMenuItem" reason:err userInfo:nil];
  }

  self.elemID = menuID;
  self.ID = tag[@"ID"];
  self.compoID = tag[@"CompoID"];
  self.separator = nil;
  self.title = @"";
  self.enabled = YES;
  self.toolTip = nil;
  self.onClick = nil;
  self.selector = nil;
  self.keys = nil;

  NSDictionary *attributes = tag[@"Attributes"];
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

    self.toolTip = attributes[@"title"];
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
  NSString *name = tag[@"Name"];
  if (![name isEqual:@"menu"]) {
    NSString *err = [NSString
        stringWithFormat:@"cannot create a MenuContainer from a %@", name];
    @throw [NSException exceptionWithName:@"ErrMenuContainer"
                                   reason:err
                                 userInfo:nil];
  }

  self.elemID = menuID;
  self.ID = tag[@"ID"];
  self.compoID = tag[@"CompoID"];
  self.title = @"";

  [self removeAllItems];
  self.children = [[NSMutableArray alloc] init];

  NSDictionary *attributes = tag[@"Attributes"];
  if (attributes != nil) {
    NSString *label = attributes[@"label"];
    if (label != nil) {
      self.title = label;
    }
  }

  NSArray *children = tag[@"Children"];
  if (children == nil) {
    return self;
  }

  for (NSDictionary *child in children) {
    MenuItem *childItem = [[MenuItem alloc] initWithMenuID:menuID andTag:child];

    NSString *childName = child[@"Name"];
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
+ (void) new:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];

    Menu *menu = [[Menu alloc] init];
    menu.ID = ID;
    driver.elements[ID] = menu;

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)load:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    NSString *ID = in[@"ID"];
    NSDictionary *tag = in[@"Tag"];

    Menu *menu = driver.elements[ID];
    if (menu == nil) {
      [NSException raise:@"ErrNoMenu" format:@"no menu with id %@", ID];
    }

    menu.root = [[MenuContainer alloc] initWithMenuID:ID andTag:tag];
    menu.root.delegate = menu;

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)render:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    NSString *ID = in[@"ID"];
    NSDictionary *tag = in[@"Tag"];

    NSString *tagID = tag[@"ID"];
    NSString *name = tag[@"Name"];

    Menu *menu = driver.elements[ID];
    if (menu == nil) {
      [NSException raise:@"ErrNoMenu" format:@"no menu with id %@", ID];
    }

    id elem = [menu elementByID:tagID];
    if (elem == nil) {
      [NSException raise:@"ErrElemNotFound"
                  format:@"no element with id %@", tagID];
    }

    // Menu container.
    // Should occur only for the root menu container.
    if ([elem isKindOfClass:[MenuContainer class]]) {
      if (![name isEqual:@"menu"]) {
        [NSException raise:@"ErrNoMenu"
                    format:@"root tag must be a menu: %@", name];
      }

      MenuContainer *container = (MenuContainer *)elem;
      container = [container initWithMenuID:menu.ID andTag:tag];

      [driver.macRPC return:returnID withOutput:nil andError:nil];
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

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)renderAttributes:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    NSString *ID = in[@"ID"];
    NSDictionary *tag = in[@"Tag"];
    NSString *tagID = tag[@"ID"];
    NSDictionary *attributes = tag[@"Attributes"];

    BOOL separator = NO;
    if (attributes != nil) {
      separator = attributes[@"separator"] != nil ? YES : NO;
    }

    Menu *menu = driver.elements[ID];
    if (menu == nil) {
      [NSException raise:@"ErrNoMenu" format:@"no menu with id %@", ID];
    }

    id elem = [menu elementByID:tagID];
    if (elem == nil) {
      [NSException raise:@"ErrElemNotFound"
                  format:@"no element with id %@", tagID];
    }

    // Menu container.
    // Should occur only for the root menu container.
    if ([elem isKindOfClass:[MenuContainer class]]) {
      MenuContainer *container = (MenuContainer *)elem;
      container = [container initWithMenuID:menu.ID andTag:tag];

      [driver.macRPC return:returnID withOutput:nil andError:nil];
      return;
    }

    // Menu item.
    MenuItem *item = (MenuItem *)elem;

    if ([item isSeparator] != separator) {
      MenuItem *newItem = [[MenuItem alloc] initWithMenuID:menu.ID andTag:tag];
      MenuContainer *container = (MenuContainer *)item.menu;
      NSInteger index = [container.children indexOfObject:item];

      [container removeChildAtIndex:index];
      [container insertChild:newItem atIndex:index];
    }

    item = [item initWithMenuID:menu.ID andTag:tag];

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
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
}

+ (void) delete:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    NSString *ID = in[@"ID"];
    [driver.elements removeObjectForKey:ID];

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}
@end