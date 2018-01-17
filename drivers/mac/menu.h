#ifndef menu_h
#define menu_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface MenuItem : NSMenuItem
@property NSString *ID;
@property NSString *compoID;
@property NSString *elemID;
@property NSString *onClick;
@property NSMenuItem *separator;

- (void)setupOnClick:(NSString *)selector;
- (void)clicked:(id)sender;
- (void)setupKeys:(NSString *)keys;
@end

@interface MenuContainer : NSMenu
@property NSString *ID;
@property NSString *compoID;
@property NSMutableArray<MenuItem *> *children;

- (instancetype)init;
- (void)addChild:(MenuItem *)child;
@end

@interface Menu : NSObject <NSMenuDelegate>
@property NSString *ID;
@property MenuContainer *root;

+ (bridge_result)newMenu:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)load:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)render:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)renderAttributes:(NSURLComponents *)url
                          payload:(NSString *)payload;
- (MenuContainer *)newContainer:(NSDictionary *)map;
- (MenuItem *)newItem:(NSDictionary *)map;
- (id)elementByID:(NSString *)ID;
- (id)elementFromContainer:(MenuContainer *)container ID:(NSString *)ID;
- (id)elementFromItem:(MenuItem *)item ID:(NSString *)ID;
@end

#endif /* menu_h */
