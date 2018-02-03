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

+ (bridge_result)newMenu:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)load:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)render:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)renderAttributes:(NSURLComponents *)url
                          payload:(NSString *)payload;
- (id)elementByID:(NSString *)ID;
- (id)elementFromContainer:(MenuContainer *)container ID:(NSString *)ID;
- (id)elementFromItem:(MenuItem *)item ID:(NSString *)ID;
@end

#endif /* menu_h */
