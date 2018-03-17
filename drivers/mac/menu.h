#ifndef menu_h
#define menu_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface MenuItem : NSMenuItem
@property NSString *elemID;
@property NSString *ID;
@property NSString *compoID;
@property NSString *onClick;
@property NSString *selector;
@property NSString *keys;
@property NSMenuItem *separator;

- (instancetype)initWithMenuID:(NSString *)menuID andTag:(NSDictionary *)tag;
- (void)setupOnClick;
- (void)clicked:(id)sender;
- (void)setupKeys;
@end

@interface MenuContainer : NSMenu
@property NSString *elemID;
@property NSString *ID;
@property NSString *compoID;
@property NSMutableArray<MenuItem *> *children;

- (instancetype)initWithMenuID:(NSString *)menuID andTag:(NSDictionary *)tag;
- (void)addChild:(MenuItem *)child;
- (void)insertChild:(MenuItem *)child atIndex:(NSInteger)index;
- (void)removeChildAtIndex:(NSInteger)index;
@end

@interface Menu : NSObject <NSMenuDelegate>
@property NSString *ID;
@property MenuContainer *root;

+ (void) new:(NSDictionary *)in return:(NSString *)returnID;
+ (void)load:(NSDictionary *)in return:(NSString *)returnID;
+ (void)render:(NSDictionary *)in return:(NSString *)returnID;
+ (void)renderAttributes:(NSDictionary *)in return:(NSString *)returnID;
+ (void) delete:(NSDictionary *)in return:(NSString *)returnID;
- (id)elementByID:(NSString *)ID;
- (id)elementFromContainer:(MenuContainer *)container ID:(NSString *)ID;
- (id)elementFromItem:(MenuItem *)item ID:(NSString *)ID;
@end

#endif /* menu_h */
