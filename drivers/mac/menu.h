#ifndef menu_h
#define menu_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface Menu : NSObject <NSMenuDelegate>
@property NSString *ID;
@property MenuContainer *Root;
@property NSMutableDictionary<NSString *, id> *elements;

+ (bridge_result)newMenu:(NSURLComponents *)url payload:(NSString *)payload;
@end

@interface MenuContainer : NSMenu
@property NSString *ID;
@end

@interface MenuItem : NSMenuItem
@property NSString *ID;
@property NSString *onClick;
@property BOOL isSeparator;
@end

#endif /* menu_h */
