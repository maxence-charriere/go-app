#ifndef window_h
#define window_h

#import "bridge.h"
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface Window : NSWindowController <NSWindowDelegate, WKNavigationDelegate,
                                        WKUIDelegate, WKScriptMessageHandler>
@property NSString *ID;
@property WKWebView *webview;

+ (bridge_result)newWindow:(NSURL *)url payload:(NSString *)payload;
@end

@interface WindowTitleBar : NSView
@end

#endif /* window_h */