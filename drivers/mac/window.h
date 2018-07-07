#ifndef window_h
#define window_h

#import "bridge.h"
#import "webview.h"
#if MAC_OS_X_VERSION_MAX_ALLOWED < 101200
#define NSWindowStyleMaskResizable NSResizableWindowMask
#define NSEventTypeLeftMouseDown NSLeftMouseDown
#define NSEventTypeLeftMouseUp NSLeftMouseUp
#define NSEventTypeRightMouseDown NSRightMouseDown
#define NSEventTypeRightMouseUp NSRightMouseUp
#define NSEventTypeOtherMouseDown NSOtherMouseDown
#define NSEventTypeOtherMouseUp NSOtherMouseUp
#define NSEventTypeScrollWheel NSScrollWheel
#define NSEventTypeMouseMoved NSMouseMoved
#define NSEventTypeLeftMouseDragged NSLeftMouseDragged
#define NSEventTypeRightMouseDragged NSRightMouseDragged
#define NSEventTypeOtherMouseDragged NSOtherMouseDragged
#define NSCompositingOperationCopy NSCompositeCopy
#define NSCompositingOperationSourceIn NSCompositeSourceIn
#define NSEventTypeFlagsChanged NSFlagsChanged
#define NSWindowStyleMaskTitled NSTitledWindowMask
#define NSWindowStyleMaskClosable NSClosableWindowMask
#define NSWindowStyleMaskMiniaturizable NSMiniaturizableWindowMask
#define NSWindowStyleMaskBorderless NSBorderlessWindowMask
#define NSWindowStyleMaskFullSizeContentView NSBorderlessWindowMask
#endif

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
