#include "menu.h"
#include "driver.h"
#include "json.h"

@implementation MenuItem
+ (instancetype)create:(NSString *)ID inMenu:(NSString *)menuID {
  return nil;
}

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
  Driver *driver = [Driver current];

  NSDictionary *mapping = @{
    @"CompoID" : self.compoID,
    @"Target" : self.onClick,
    @"JSONValue" : @"{}",
  };

  NSDictionary *in = @{
    @"ID" : self.elemID,
    @"Mapping" : [JSONEncoder encode:mapping],
  };

  [driver.goRPC call:@"menus.OnCallback" withInput:in onUI:YES];
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
+ (instancetype)create:(NSString *)ID inMenu:(NSString *)menuID {
  return nil;
}

// - (instancetype)initWithMenuID:(NSString *)menuID andTag:(NSDictionary *)tag
// {
//   NSString *name = tag[@"Name"];
//   if (![name isEqual:@"menu"]) {
//     NSString *err = [NSString
//         stringWithFormat:@"cannot create a MenuContainer from a %@", name];
//     @throw [NSException exceptionWithName:@"ErrMenuContainer"
//                                    reason:err
//                                  userInfo:nil];
//   }

//   self.elemID = menuID;
//   self.ID = tag[@"ID"];
//   self.compoID = tag[@"CompoID"];
//   self.title = @"";

//   [self removeAllItems];
//   self.children = [[NSMutableArray alloc] init];

//   NSDictionary *attributes = tag[@"Attributes"];
//   if (attributes != nil) {
//     NSString *label = attributes[@"label"];
//     if (label != nil) {
//       self.title = label;
//     }
//   }

//   NSArray *children = tag[@"Children"];
//   if (children == nil) {
//     return self;
//   }

//   for (NSDictionary *child in children) {
//     MenuItem *childItem = [[MenuItem alloc] initWithMenuID:menuID
//     andTag:child];

//     NSString *childName = child[@"Name"];
//     if ([childName isEqual:@"menu"]) {
//       childItem.submenu =
//           [[MenuContainer alloc] initWithMenuID:menuID andTag:child];
//     }

//     [self addChild:childItem];
//   }

//   return self;
// }

// - (instancetype)init {
//   self = [super init];
//   self.children = [[NSMutableArray alloc] init];
//   return self;
// }

// - (void)addChild:(MenuItem *)child {
//   [self.children addObject:child];

//   if ([child isSeparator]) {
//     [self addItem:child.separator];
//     child.menu = self;
//     return;
//   }

//   [self addItem:child];
// }

// - (void)insertChild:(MenuItem *)child atIndex:(NSInteger)index {
//   [self.children insertObject:child atIndex:index];

//   if ([child isSeparator]) {
//     [self insertItem:child.separator atIndex:index];
//     child.menu = self;
//   } else {
//     [self insertItem:child atIndex:index];
//   }
// }

// - (void)removeChildAtIndex:(NSInteger)index {
//   [self.children removeObjectAtIndex:index];
//   [self removeItemAtIndex:index];
// }

@end

@implementation Menu
+ (void) new:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    NSString *ID = in[@"ID"];

    Menu *menu = [[Menu alloc] init];
    menu.ID = ID;
    menu.nodes = [[NSMutableDictionary alloc] init];

    Driver *driver = [Driver current];
    driver.elements[ID] = menu;
    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)render:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];

    NSString *ID = in[@"ID"];
    NSArray *changes = [JSONDecoder decode:in[@"Changes"]];

    Menu *menu = driver.elements[ID];
    if (menu == nil) {
      [NSException raise:@"ErrNoMenu" format:@"no menu with id %@", ID];
    }

    NSDictionary<NSString *, NSNumber *> *typeMap = @{
      @"createText" : @0,
      @"setText" : @1,
      @"createElem" : @2,
      @"setAttrs" : @3,
      @"appendChild" : @4,
      @"removeChild" : @5,
      @"replaceChild" : @6,
      @"mountElem" : @7,
      @"createCompo" : @8,
      @"setCompoRoot" : @9,
      @"deleteNode" : @10
    };

    for (NSDictionary *c in changes) {
      NSString *type = c[@"Type"];

      switch (typeMap[type].intValue) {
      case 2:
        [menu createElem:c];
        break;

      case 3:
        [menu setAttrs:c];
        break;

      case 4:
        [menu appendChild:c];
        break;

      case 5:
        [menu removeChild:c];
        break;

      case 6:
        [menu replaceChild:c];
        break;

      case 7:
        [menu mountElem:c];
        break;

      case 8:
        [menu createCompo:c];
        break;

      case 9:
        [menu setCompoRoot:c];
        break;

      case 10:
        [menu deleteNode:c];
        break;

      default:
        [NSException raise:@"ErrChange"
                    format:@"%@ change is not supported", type];
      }
    }

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)createElem:(NSDictionary *)change {
  NSString *ID = change[@"Value"][@"ID"];
  NSString *TagName = change[@"Value"][@"TagName"];

  if ([TagName isEqual:@"menu"]) {
    [MenuContainer create:ID inMenu:self.ID];
    return;
  }

  if ([TagName isEqual:@"menuitem"]) {
    [MenuItem create:ID inMenu:self.ID];
    return;
  }

  [NSException raise:@"ErrMenu"
              format:@"menu does not support %@ tag", TagName];
}

- (void)setAttrs:(NSDictionary *)change {
}

- (void)appendChild:(NSDictionary *)change {
}

- (void)removeChild:(NSDictionary *)change {
}

- (void)replaceChild:(NSDictionary *)change {
}

- (void)mountElem:(NSDictionary *)change {
}

- (void)createCompo:(NSDictionary *)change {
}

- (void)setCompoRoot:(NSDictionary *)change {
}

- (void)deleteNode:(NSDictionary *)change {
}

- (void)menuDidClose:(NSMenu *)menu {
  NSDictionary *in = @{
    @"ID" : self.ID,
  };

  Driver *driver = [Driver current];
  [driver.goRPC call:@"menus.OnClose" withInput:in onUI:YES];
}

+ (void) delete:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    NSString *ID = in[@"ID"];

    Driver *driver = [Driver current];
    [driver.elements removeObjectForKey:ID];
    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}
@end