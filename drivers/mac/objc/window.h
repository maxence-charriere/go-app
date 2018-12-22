#ifndef window_h
#define window_h

#import "bridge.h"
#import "retro.h"
#import "webview.h"

@interface Window : NSWindowController <NSWindowDelegate, WKNavigationDelegate,
                                        WKUIDelegate, WKScriptMessageHandler>
@property NSString *ID;
@property(weak) AppWebView *webview;
@property NSURL *loadURL;
@property NSURL *baseURL;
@property NSString *loadReturnID;

+ (void) new:(NSDictionary *)in return:(NSString *)returnID;
- (void)configBackgroundColor:(NSString *)color frosted:(BOOL)frosted;
- (void)configWebview;
- (void)configTitlebar:(NSString *)title hidden:(BOOL)isHidden;
+ (void)load:(NSDictionary *)in return:(NSString *)returnID;
+ (void)render:(NSDictionary *)in return:(NSString *)returnID;
+ (void)position:(NSDictionary *)in return:(NSString *)returnID;
+ (void)move:(NSDictionary *)in return:(NSString *)returnID;
+ (void)center:(NSDictionary *)in return:(NSString *)returnID;
+ (void)size:(NSDictionary *)in return:(NSString *)returnID;
+ (void)resize:(NSDictionary *)in return:(NSString *)returnID;
+ (void)focus:(NSDictionary *)in return:(NSString *)returnID;
+ (void)toggleFullScreen:(NSDictionary *)in return:(NSString *)returnID;
+ (void)toggleMinimize:(NSDictionary *)in return:(NSString *)returnID;
+ (void)close:(NSDictionary *)in return:(NSString *)returnID;
@end

@interface WindowTitleBar : NSView
@end

#endif /* window_h */
