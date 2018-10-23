#include "menu.h"
#include "driver.h"
#include "image.h"
#include "json.h"

@implementation MenuItem
+ (instancetype)create:(NSString *)ID
               compoID:(NSString *)compoID
                inMenu:(NSString *)elemID {
  MenuItem *i =
      [[MenuItem alloc] initWithTitle:@"" action:NULL keyEquivalent:@""];
  i.ID = ID;
  i.compoID = compoID;
  i.elemID = elemID;
  return i;
}

- (void)setAttr:(NSString *)key value:(NSString *)value {
  if ([key isEqual:@"separator"] && self.separator == nil) {
    [self setSeparator];
    return;
  }

  if ([key isEqual:@"label"]) {
    self.title = value != nil ? value : @"";
    return;
  }

  if ([key isEqual:@"disabled"]) {
    self.enabled = NO;
    return;
  }

  if ([key isEqual:@"title"]) {
    self.toolTip = value;
    return;
  }

  if ([key isEqual:@"checked"]) {
    self.state = NSControlStateValueOn;
    return;
  }

  if ([key isEqual:@"keys"]) {
    self.keys = value;
    [self setupKeys];
    return;
  }

  if ([key isEqual:@"icon"]) {
    NSString *icon = value;
    icon = icon != nil ? icon : @"";

    if (icon.length != 0) {
      NSBundle *mainBundle = [NSBundle mainBundle];
      icon =
          [NSString stringWithFormat:@"%@/%@", mainBundle.resourcePath, icon];
    }

    if (![self.icon isEqual:icon]) {
      self.icon = icon;
      [self setIconWithPath:icon];
    }
    return;
  }

  if ([key isEqual:@"onclick"]) {
    self.onClick = value;
    [self setupOnClick];
    return;
  }

  if ([key isEqual:@"selector"]) {
    self.selector = value;
    [self setupOnClick];
    return;
  }
}

- (void)delAttr:(NSString *)key {
  if ([key isEqual:@"separator"] && self.separator != nil) {
    [self unsetSeparator];
    return;
  }

  if ([key isEqual:@"label"]) {
    self.title = @"";
    return;
  }

  if ([key isEqual:@"disabled"]) {
    self.enabled = YES;
    return;
  }

  if ([key isEqual:@"title"]) {
    self.toolTip = nil;
    return;
  }

  if ([key isEqual:@"checked"]) {
    self.state = NSControlStateValueOff;
    return;
  }

  if ([key isEqual:@"keys"]) {
    self.keys = nil;
    [self setupKeys];
    return;
  }

  if ([key isEqual:@"icon"]) {
    NSString *icon = @"";
    self.icon = icon;
    [self setIconWithPath:icon];
    return;
  }

  if ([key isEqual:@"onclick"]) {
    self.onClick = nil;
    [self setupOnClick];
    return;
  }

  if ([key isEqual:@"selector"]) {
    self.selector = nil;
    [self setupOnClick];
    return;
  }
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
+ (instancetype)create:(NSString *)ID
               compoID:(NSString *)compoID
                inMenu:(NSString *)elemID {
  MenuContainer *m = [[MenuContainer alloc] initWithTitle:@""];
  m.ID = ID;
  m.compoID = compoID;
  m.elemID = elemID;
  return m;
}

- (void)setAttr:(NSString *)key value:(NSString *)value {
  if ([key isEqual:@"label"]) {
    self.title = value != nil ? value : @"";
  } else if ([key isEqual:@"disabled"]) {
    self.disabled = true;
  }

  [self updateParentItem];
}

- (void)delAttr:(NSString *)key {
  if ([key isEqual:@"label"]) {
    self.title = @"";
  } else if ([key isEqual:@"disabled"]) {
    self.disabled = false;
  }

  [self updateParentItem];
}

- (void)updateParentItem {
  NSMenu *supermenu = self.supermenu;
  if (supermenu == nil) {
    return;
  }

  // Updating parent menuitem title.
  for (NSMenuItem *i in supermenu.itemArray) {
    if (i.submenu == self) {
      i.title = self.title;
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

  self.actions = @{
    @"setRoot" : @0,
    @"newNode" : @1,
    @"delNode" : @2,
    @"setAttr" : @3,
    @"delAttr" : @4,
    @"setText" : @5,
    @"appendChild" : @6,
    @"removeChild" : @7,
    @"replaceChild" : @8,
  };

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

    for (NSDictionary *c in changes) {
      NSString *action = c[@"Action"];

      switch (menu.actions[action].intValue) {
      case 0:
        [menu setRootNode:c];
        break;

      case 1:
        [menu newNode:c];
        break;

      case 2:
        [menu delNode:c];
        break;

      case 3:
        [menu setAttr:c];
        break;

      case 4:
        [menu delAttr:c];
        break;

      case 6:
        [menu appendChild:c];
        break;

      case 7:
        [menu removeChild:c];
        break;

      case 8:
        [menu replaceChild:c];
        break;

      default:
        [NSException raise:@"ErrChange"
                    format:@"%@ change is not supported", action];
      }
    }

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)setRootNode:(NSDictionary *)change {
}

- (void)newNode:(NSDictionary *)change {
  NSString *nodeID = change[@"NodeID"];
  NSString *compoID = change[@"CompoID"];
  NSString *type = change[@"Type"];
  BOOL isCompo = [change[@"IsCompo"] boolValue];

  if (isCompo) {
    MenuCompo *c = [[MenuCompo alloc] init];
    c.ID = nodeID;
    c.type = type;
    self.nodes[nodeID] = c;
    return;
  }

  if ([type isEqual:@"menu"]) {
    self.nodes[nodeID] =
        [MenuContainer create:nodeID compoID:compoID inMenu:self.ID];
    return;
  }

  if ([type isEqual:@"menuitem"]) {
    self.nodes[nodeID] =
        [MenuItem create:nodeID compoID:compoID inMenu:self.ID];
    return;
  }

  [NSException raise:@"ErrMenu" format:@"menu does not support %@ tag", type];
}

- (void)delNode:(NSDictionary *)change {
  [self.nodes removeObjectForKey:change[@"NodeID"]];
}

- (void)setAttr:(NSDictionary *)change {
  id node = self.nodes[change[@"NodeID"]];
  if (node == nil) {
    return;
  }

  NSString *key = change[@"Key"];
  NSString *value = change[@"Value"];

  if ([node isKindOfClass:[MenuContainer class]]) {
    MenuContainer *m = node;
    [m setAttr:key value:value];
    return;
  }

  if ([node isKindOfClass:[MenuItem class]]) {
    MenuItem *mi = node;
    [mi setAttr:key value:value];
    return;
  }

  [NSException raise:@"ErrMenu" format:@"unknown menu element"];
}

- (void)delAttr:(NSDictionary *)change {
  id node = self.nodes[change[@"NodeID"]];
  if (node == nil) {
    return;
  }

  NSString *key = change[@"Key"];

  if ([node isKindOfClass:[MenuContainer class]]) {
    MenuContainer *m = node;
    [m delAttr:key];
    return;
  }

  if ([node isKindOfClass:[MenuItem class]]) {
    MenuItem *mi = node;
    [mi delAttr:key];
    return;
  }

  [NSException raise:@"ErrMenu" format:@"unknown menu element"];
}

- (void)appendChild:(NSDictionary *)change {
}

- (void)removeChild:(NSDictionary *)change {
}

- (void)replaceChild:(NSDictionary *)change {
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