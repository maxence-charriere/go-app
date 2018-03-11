#ifndef window_h
#define window_h

#import "bridge.h"
#import "webview.h"

@interface Window : NSWindowController <NSWindowDelegate, WKNavigationDelegate,
                                        WKUIDelegate, WKScriptMessageHandler>
@property NSString *ID;
@property(weak) AppWebView *webview;
@property NSURL *loadURL;
@property NSURL *baseURL;
@property NSString *loadReturnID;

+ (void) new:(NSDictionary *)in return:(NSString *)returnID;
- (void)configBackgroundColor:(NSString *)color
                     vibrancy:(NSVisualEffectMaterial)vibrancy;
- (void)configWebview;
- (void)configTitlebar:(NSString *)title hidden:(BOOL)isHidden;
+ (void)load:(NSDictionary *)in return:(NSString *)returnID;
+ (void)render:(NSDictionary *)in return:(NSString *)returnID;
+ (void)renderAttributes:(NSDictionary *)in return:(NSString *)returnID;

+ (bridge_result)position:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)move:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)center:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)size:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)resize:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)focus:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)toggleFullScreen:(NSURLComponents *)url
                          payload:(NSString *)payload;
+ (bridge_result)toggleMinimize:(NSURLComponents *)url
                        payload:(NSString *)payload;
+ (bridge_result)close:(NSURLComponents *)url payload:(NSString *)payload;
@end

@interface WindowTitleBar : NSView
@end

#endif /* window_h */
