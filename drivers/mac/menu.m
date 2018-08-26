#include "menu.h"
#include "driver.h"
#include "image.h"
#include "json.h"

@implementation MenuItem
+ (instancetype)create:(NSString *)ID inMenu:(NSString *)elemID {
  MenuItem *i =
      [[MenuItem alloc] initWithTitle:@"" action:NULL keyEquivalent:@""];
  i.ID = ID;
  i.elemID = elemID;
  return i;
}

- (void)setAttrs:(NSDictionary<NSString *, NSString *> *)attrs {
  BOOL separator = attrs[@"separator"] != nil ? YES : NO;
  if (separator && self.separator == nil) {
    [self setSeparator];
  } else if (!separator && self.separator != nil) {
    [self unsetSeparator];
  }

  NSString *label = attrs[@"label"];
  label = label != nil ? label : @"";
  self.title = label;

  self.enabled = attrs[@"disabled"] == nil ? true : false;
  self.toolTip = attrs[@"title"];
  self.onClick = attrs[@"onclick"];
  self.selector = attrs[@"selector"];
  self.keys = attrs[@"keys"];

  if (attrs[@"checked"] != nil) {
    self.state = NSControlStateValueOn;
  } else {
    self.state = NSControlStateValueOff;
  }

  NSString *icon = attrs[@"icon"];
  icon = icon != nil ? icon : @"";

  if (icon.length != 0) {
    NSBundle *mainBundle = [NSBundle mainBundle];
    icon = [NSString stringWithFormat:@"%@/%@", mainBundle.resourcePath, icon];
  }

  if (![self.icon isEqual:icon]) {
    self.icon = icon;
    [self setIconWithPath:icon];
  }

  [self setupOnClick];
  [self setupKeys];
}

- (void)setSeparator {
  NSMenuItem *sep = [NSMenuItem separatorItem];
  self.separator = sep;

  MenuContainer *parent = (MenuContainer *)self.menu;
  if (parent == nil) {
    return;
  }

  NSInteger index = [parent indexOfItem:self];
  [parent removeItemAtIndex:index];
  [parent insertItem:sep atIndex:index];
}

- (void)unsetSeparator {
  NSMenuItem *sep = self.separator;
  self.separator = nil;

  MenuContainer *parent = (MenuContainer *)sep.menu;
  if (parent == nil) {
    return;
  }

  NSInteger index = [parent indexOfItem:sep];
  [parent removeItemAtIndex:index];
  [parent insertItem:self atIndex:index];
}

- (void)setIconWithPath:(NSString *)icon {
  if (icon.length == 0) {
    self.image = nil;
    return;
  }

  CGFloat menuBarHeight = [[NSApp mainMenu] menuBarHeight];

  NSImage *img = [[NSImage alloc] initByReferencingFile:icon];
  self.image = [NSImage resizeImage:img
                  toPixelDimensions:NSMakeSize(menuBarHeight, menuBarHeight)];
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
    @"FieldOrMethod" : self.onClick,
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
+ (instancetype)create:(NSString *)ID inMenu:(NSString *)elemID {
  MenuContainer *m = [[MenuContainer alloc] initWithTitle:@""];
  m.ID = ID;
  m.elemID = elemID;
  return m;
}

- (void)setAttrs:(NSDictionary<NSString *, NSString *> *)attrs {
  NSString *label = attrs[@"label"];
  label = label != nil ? label : @"";
  self.title = label;

  self.disabled = attrs[@"disabled"] != nil ? true : false;

  MenuContainer *supermenu = (MenuContainer *)self.supermenu;
  if (supermenu == nil) {
    return;
  }

  // Updating parent menuitem title.
  for (NSMenuItem *i in supermenu.itemArray) {
    if (i.submenu == self) {
      i.title = label;
      i.enabled = !self.disabled;
      return;
    }
  }
}

- (void)insertChild:(id)child atIndex:(NSInteger)index {
  if ([child isKindOfClass:[MenuContainer class]]) {
    MenuContainer *c = child;
    NSMenuItem *item = [[NSMenuItem alloc] initWithTitle:c.title
                                                  action:NULL
                                           keyEquivalent:@""];

    item.submenu = c;
    item.enabled = !c.disabled;
    [self insertItem:item atIndex:index];
    return;
  }

  MenuItem *item = child;

  if (item.separator != nil) {
    [self insertItem:item.separator atIndex:index];
    return;
  }

  [self insertItem:item atIndex:index];
}

- (void)appendChild:(id)child {
  [self insertChild:child atIndex:self.numberOfItems];
}

- (void)removeChild:(id)child {
  if ([child isKindOfClass:[MenuContainer class]]) {
    for (NSMenuItem *c in self.itemArray) {
      if (c.submenu == child) {
        [self removeItem:c];
        return;
      }
    }

    return;
  }

  MenuItem *item = child;

  if (item.separator != nil) {
    [self removeItem:item.separator];
    return;
  }

  [self removeItem:item];
}

- (void)replaceChild:(id)old with:(id) new {
  NSInteger index = -1;

  if ([old isKindOfClass:[MenuContainer class]]) {
    NSArray<NSMenuItem *> *children = self.itemArray;

    for (int i = 0; i < children.count; ++i) {
      if (children[i].submenu == old) {
        index = i;
        break;
      }
    }
  } else {
    MenuItem *item = old;

    if (item.separator != nil) {
      index = [self indexOfItem:item.separator];
    } else {
      index = [self indexOfItem:item];
    }
  }

  if (index < 0) {
    return;
  }

  [self removeItemAtIndex:index];
  [self insertChild:new atIndex:index];
}
@end

@implementation Menu
- (instancetype)initWithID:(NSString *)ID {
  self = [super init];

  self.ID = ID;
  self.nodes = [[NSMutableDictionary alloc] init];

  return self;
}

+ (void) new:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    NSString *ID = in[@"ID"];
    Menu *menu = [[Menu alloc] initWithID:ID];

    Driver *driver = [Driver current];
    driver.elements[ID] = menu;
    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

+ (void)load:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];

    Menu *menu = driver.elements[ID];
    if (menu == nil) {
      [NSException raise:@"ErrNoMenu" format:@"no menu with id %@", ID];
    }

    menu.root = nil;
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
        [menu createElem:c[@"Value"]];
        break;

      case 3:
        [menu setAttrs:c[@"Value"]];
        break;

      case 4:
        [menu appendChild:c[@"Value"]];
        break;

      case 5:
        [menu removeChild:c[@"Value"]];
        break;

      case 6:
        [menu replaceChild:c[@"Value"]];
        break;

      case 7:
        [menu mountElem:c[@"Value"]];
        break;

      case 8:
        [menu createCompo:c[@"Value"]];
        break;

      case 9:
        [menu setCompoRoot:c[@"Value"]];
        break;

      case 10:
        [menu deleteNode:c[@"Value"]];
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
  NSString *ID = change[@"ID"];
  NSString *TagName = change[@"TagName"];

  if ([TagName isEqual:@"menu"]) {
    self.nodes[ID] = [MenuContainer create:ID inMenu:self.ID];
    return;
  }

  if ([TagName isEqual:@"menuitem"]) {
    self.nodes[ID] = [MenuItem create:ID inMenu:self.ID];
    return;
  }

  [NSException raise:@"ErrMenu"
              format:@"menu does not support %@ tag", TagName];
}

- (void)setAttrs:(NSDictionary *)change {
  NSDictionary<NSString *, NSString *> *attrs = change[@"Attrs"];
  if (attrs == nil) {
    return;
  }

  id node = self.nodes[change[@"ID"]];
  if (node == nil) {
    return;
  }

  if ([node isKindOfClass:[MenuContainer class]]) {
    MenuContainer *m = node;
    [m setAttrs:attrs];
    return;
  }

  if ([node isKindOfClass:[MenuItem class]]) {
    MenuItem *mi = node;
    [mi setAttrs:attrs];
    return;
  }

  [NSException raise:@"ErrMenu" format:@"unknown menu element"];
}

- (void)appendChild:(NSDictionary *)change {
  id child = self.nodes[change[@"ChildID"]];
  child = [self childElem:child];
  if (child == nil) {
    return;
  }

  NSString *parentID = change[@"ParentID"];

  if ([parentID isEqual:@"root:"]) {
    if ([child isKindOfClass:[MenuItem class]]) {
      [NSException raise:@"ErrMenu" format:@"menu root is a menuitem"];
    }

    MenuContainer *m = child;
    m.delegate = self;
    self.root = m;
    return;
  }

  MenuContainer *parent = self.nodes[parentID];
  if (parent == nil) {
    return;
  }

  [parent appendChild:child];
}

- (void)removeChild:(NSDictionary *)change {
  NSString *parentID = change[@"ParentID"];

  if ([parentID isEqual:@"root:"]) {
    self.root = nil;
    return;
  }

  MenuContainer *parent = self.nodes[parentID];
  if (parent == nil) {
    return;
  }

  id child = self.nodes[change[@"ChildID"]];
  child = [self childElem:child];
  if (child == nil) {
    return;
  }

  [parent removeChild:child];
}

- (void)replaceChild:(NSDictionary *)change {
  NSString *parentID = change[@"ParentID"];

  if ([parentID isEqual:@"root:"]) {
    [NSException raise:@"ErrMenu" format:@"root menu can't be replaced"];
  }

  MenuContainer *parent = self.nodes[parentID];
  if (parent == nil) {
    return;
  }

  id newChild = self.nodes[change[@"ChildID"]];
  newChild = [self childElem:newChild];
  if (newChild == nil) {
    return;
  }

  id oldChild = self.nodes[change[@"OldID"]];
  oldChild = [self childElem:oldChild];
  if (oldChild == nil) {
    return;
  }

  [parent replaceChild:oldChild with:newChild];
}

- (void)mountElem:(NSDictionary *)change {
  id node = self.nodes[change[@"ID"]];
  if (node == nil) {
    return;
  }

  NSString *compoID = change[@"CompoID"];

  if ([node isKindOfClass:[MenuContainer class]]) {
    MenuContainer *m = node;
    m.compoID = compoID;
    return;
  }

  MenuItem *i = node;
  i.compoID = compoID;
}

- (void)createCompo:(NSDictionary *)change {
  MenuCompo *c = [[MenuCompo alloc] init];
  c.ID = change[@"ID"];
  c.name = change[@"Name"];
  self.nodes[c.ID] = c;
}

- (void)setCompoRoot:(NSDictionary *)change {
  MenuCompo *c = self.nodes[change[@"ID"]];
  if (c == nil) {
    return;
  }

  c.rootID = change[@"RootID"];
}

- (void)deleteNode:(NSDictionary *)change {
  [self.nodes removeObjectForKey:change[@"ID"]];
}

- (id)childElem:(id)node {
  if (![node isKindOfClass:[MenuCompo class]]) {
    return node;
  }

  MenuCompo *c = node;
  return self.nodes[c.rootID];
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

@implementation MenuCompo
@end