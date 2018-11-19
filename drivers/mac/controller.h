#ifndef controller_h
#define controller_h

#import <GameController/GCController.h>

typedef enum ControllerInput : NSUInteger {
  DirectionalPad,
  LeftThumbstick,
  RightThumbstick,
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
+ (void)emitDirection:(NSString *)elemID
                input:(ControllerInput)input
                    x:(float)x
                    y:(float)y;
+ (void)emitButton:(NSString *)elemID
             input:(ControllerInput)input
             value:(float)value
           pressed:(BOOL)pressed;
+ (void)close:(NSDictionary *)in return:(NSString *)returnID;
+ (void)listen:(NSDictionary *)in return:(NSString *)returnID;

- (void)connected;
- (void)disconnected;
@end

#endif /* controller_h */