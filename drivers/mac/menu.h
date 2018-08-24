#ifndef menu_h
#define menu_h

#import "retro.h"
#import <Cocoa/Cocoa.h>

@interface MenuItem : NSMenuItem
@property NSString *elemID;
@property NSString *ID;
@property NSString *compoID;
@property NSString *menuID;
@property NSString *onClick;
@property NSString *selector;
@property NSString *keys;
@property NSMenuItem *separator;

+ (instancetype)create:(NSString *)ID inMenu:(NSString *)menuID;

- (instancetype)initWithMenuID:(NSString *)menuID andTag:(NSDictionary *)tag;
- (void)setupOnClick;
- (void)clicked:(id)sender;
- (void)setupKeys;
@end

@interface MenuContainer : NSMenu
@property NSString *ID;
@property NSString *compoID;
@property NSString *menuID;

+ (instancetype)create:(NSString *)ID inMenu:(NSString *)menuID;
@end

@interface Menu : NSObject <NSMenuDelegate>
@property NSString *ID;
@property NSDictionary<NSString *, id> *nodes;
@property MenuContainer *root;

+ (void) new:(NSDictionary *)in return:(NSString *)returnID;
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
+ (void) delete:(NSDictionary *)in return:(NSString *)returnID;
@end

#endif /* menu_h */
