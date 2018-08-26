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

+ (instancetype)create:(NSString *)ID inMenu:(NSString *)elemID;
- (void)setAttrs:(NSDictionary<NSString *, NSString *> *)attrs;
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

+ (instancetype)create:(NSString *)ID inMenu:(NSString *)elemID;
- (void)setAttrs:(NSDictionary<NSString *, NSString *> *)attrs;
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
- (void)createElem:(NSDictionary *)change;
- (void)setAttrs:(NSDictionary *)change;
- (void)appendChild:(NSDictionary *)change;
- (void)removeChild:(NSDictionary *)change;
- (void)replaceChild:(NSDictionary *)change;
- (void)mountElem:(NSDictionary *)change;
- (void)createCompo:(NSDictionary *)change;
- (void)setCompoRoot:(NSDictionary *)change;
- (void)deleteNode:(NSDictionary *)change;
- (id)childElem:(id)node;
+ (void) delete:(NSDictionary *)in return:(NSString *)returnID;
@end

@interface MenuCompo : NSObject
@property NSString *ID;
@property NSString *rootID;
@property NSString *name;
@end

#endif /* menu_h */
