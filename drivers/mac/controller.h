#ifndef controller_h
#define controller_h

#import <GameController/GCController.h>

typedef enum ControllerInput : NSUInteger {
	DirectionalPad,
	LeftThumbstick,
	RightThumbstick,
	LeftThumbstickButton,
	RightThumbstickButton,
	A,
	B,
	X,
	Y,
	L1,
	L2,
	R1,
	R2,
	Pause
} ControllerInput;

@interface Controller : GCEventViewController
@property NSString *ID;
@property GCController *context;
@property GCExtendedGamepad *profile;
+ (void) new:(NSDictionary *)in return:(NSString *)returnID;
+ (void) emitButton:(NSString *)ID Input:(int)input Value:(float)value Pressed:(BOOL)pressed;
+ (void) emitDpad:(NSString *)ID Input:(int)input X:(float)x Y:(float)y;
+ (void) close:(NSDictionary *)in return:(NSString *)returnID;
+ (void) listen:(NSDictionary *)in return:(NSString *)returnID;

- (void) connected;
- (void) disconnected;
@end

#endif /* controller_h */