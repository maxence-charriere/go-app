#ifndef window_h
#define window_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface Window : NSWindowController <NSWindowDelegate, WKNavigationDelegate,
                                        WKUIDelegate, WKScriptMessageHandler>
@property NSString *ID;
@property(weak) WKWebView *webview;
@property NSURL *loadURL;
@property NSURL *baseURL;

+ (bridge_result)newWindow:(NSURLComponents *)url payload:(NSString *)payload;
- (void)configBackgroundColor:(NSString *)color
                     vibrancy:(NSVisualEffectMaterial)vibrancy;
- (void)configWebview;
- (void)configTitlebar:(NSString *)title hidden:(BOOL)isHidden;
+ (bridge_result)load:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)render:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)renderAttributes:(NSURLComponents *)url
                          payload:(NSString *)payload;
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
