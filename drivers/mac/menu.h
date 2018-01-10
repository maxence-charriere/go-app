#ifndef menu_h
#define menu_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface MenuContainer : NSMenu
@property NSString *ID;
@property NSString *compoID;
@end

@interface MenuItem : NSMenuItem
@property NSString *ID;
@property NSString *compoID;
@property NSString *elemID;
@property NSString *onClick;

- (void)setupOnClick:(NSString *)selector;
- (void)clicked:(id)sender;
- (void)setupKeys:(NSString *)keys;
@end

@interface Menu : NSObject <NSMenuDelegate>
@property NSString *ID;
@property MenuContainer *root;

+ (bridge_result)newMenu:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)load:(NSURLComponents *)url payload:(NSString *)payload;
- (MenuContainer *)newContainer:(NSDictionary *)map;
- (MenuItem *)newItem:(NSDictionary *)map;
@end

#endif /* menu_h */
