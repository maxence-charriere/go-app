#ifndef retro_h
#define retro_h

#import <Cocoa/Cocoa.h>

#if MAC_OS_X_VERSION_MAX_ALLOWED < 101300

/* menu.h */
#define NSControlStateValueOn NSOnState
#define NSControlStateValueOff NSOffState

#endif

#if MAC_OS_X_VERSION_MAX_ALLOWED < 101200

/* image.h */
static const NSCompositingOperation NSCompositingOperationCopy =
    NSCompositeCopy;

/* window.h */
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

/* menu.h */
#define NSEventModifierFlagControl NSControlKeyMask
#define NSEventModifierFlagShift NSShiftKeyMask
#define NSEventModifierFlagOption NSAlternateKeyMask
#define NSEventModifierFlagCommand NSCommandKeyMask
#define NSEventModifierFlagFunction NSFunctionKeyMask
#define NSControlStateValueOn NSOnState
#define NSControlStateValueOff NSOffState

#endif

#endif /* retro_h */