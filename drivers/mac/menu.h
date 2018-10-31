#ifndef menu_h
#define menu_h

#import "retro.h"
#import <Cocoa/Cocoa.h>

@interface MenuItem : NSMenuItem
@property NSString *ID;
@property NSString *compoID;
@property NSString *elemID;
@property NSString *onClick;
@property NSString *selector;
@property NSString *keys;
@property NSString *icon;
@property NSMenuItem *separator;

+ (instancetype)create:(NSString *)ID
               compoID:(NSString *)compoID
                inMenu:(NSString *)elemID;
- (void)setAttr:(NSString *)key value:(NSString *)value;
- (void)delAttr:(NSString *)key;
- (void)setSeparator;
- (void)unsetSeparator;
- (void)setIconWithPath:(NSString *)icon;
- (void)setupOnClick;
- (void)clicked:(id)sender;
- (void)setupKeys;
@end

@interface MenuContainer : NSMenu
@property NSString *ID;
@property NSString *compoID;
@property NSString *elemID;
@property BOOL disabled;

+ (instancetype)create:(NSString *)ID
               compoID:(NSString *)compoID
                inMenu:(NSString *)elemID;
- (void)setAttr:(NSString *)key value:(NSString *)value;
- (void)delAttr:(NSString *)key;
- (void)updateParentItem;
- (void)insertChild:(id)child atIndex:(NSInteger)index;
- (void)appendChild:(id)child;
- (void)removeChild:(id)child;
- (void)replaceChild:(id)old with:(id) new;
@end

@interface Menu : NSObject <NSMenuDelegate>
@property NSString *ID;
@property NSMutableDictionary<NSString *, id> *nodes;
@property MenuContainer *root;

- (instancetype)initWithID:(NSString *)ID;
+ (void) new:(NSDictionary *)in return:(NSString *)returnID;
+ (void)load:(NSDictionary *)in return:(NSString *)returnID;
+ (void)render:(NSDictionary *)in return:(NSString *)returnID;
- (void)setRootNode:(NSDictionary *)change;
- (void)newNode:(NSDictionary *)change;
- (void)delNode:(NSDictionary *)change;
- (void)setAttr:(NSDictionary *)change;
- (void)delAttr:(NSDictionary *)change;
- (void)appendChild:(NSDictionary *)change;
- (void)removeChild:(NSDictionary *)change;
- (void)replaceChild:(NSDictionary *)change;
- (id)compoRoot:(id)node;
+ (void) delete:(NSDictionary *)in return:(NSString *)returnID;
@end

@interface MenuCompo : NSObject
@property NSString *ID;
@property NSString *rootID;
@property NSString *type;
@property BOOL isRootCompo;
@end

#endif /* menu_h */
