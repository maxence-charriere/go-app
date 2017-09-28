#ifndef window_h
#define window_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface Window : NSWindowController <NSWindowDelegate, WKNavigationDelegate,
                                        WKUIDelegate, WKScriptMessageHandler>
@property NSString *ID;
@property WKWebView *webview;

+ (bridge_result)newWindow:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)position:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)move:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)center:(NSURLComponents *)url payload:(NSString *)payload;
@end

@interface WindowTitleBar : NSView
@end

#endif /* window_h */