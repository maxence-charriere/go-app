#include "controller.h"
#include "driver.h"

@implementation Controller

+ (void) new:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];

    // Create blank controller, track it for future connection.
    Controller *controller = [[Controller alloc] init];
    controller.ID = ID;
    driver.elements[ID] = controller;

    // Register for controller connection notifications.
    [[NSNotificationCenter defaultCenter]
        addObserver:controller
           selector:@selector(connected)
               name:GCControllerDidConnectNotification
             object:nil];
    [[NSNotificationCenter defaultCenter]
        addObserver:controller
           selector:@selector(disconnected)
               name:GCControllerDidDisconnectNotification
             object:nil];

    // Check if controller was connected before application started.
    if ([[GCController controllers] count] > 0) {
      [controller connected];
    }

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)connected {
  Driver *driver = [Driver current];
  self.context = [GCController controllers][0];
  self.profile = self.context.extendedGamepad;
  NSDictionary *in = @{
    @"ID" : self.ID,
  };

  [driver.goRPC call:@"controller.OnConnected" withInput:in onUI:YES];
}

- (void)disconnected {
  Driver *driver = [Driver current];
  NSDictionary *in = @{
    @"ID" : self.ID,
  };

  [driver.goRPC call:@"controller.OnDisconnected" withInput:in onUI:YES];
}

+ (void)close:(NSDictionary *)in return:(NSString *)returnID {
  Driver *driver = [Driver current];
  NSString *ID = in[@"ID"];
  Controller *controller = driver.elements[ID];

  if (controller == nil) {
    [NSException raise:@"ErrNoController"
                format:@"no controller with id %@", ID];
  }

  [driver.goRPC call:@"controller.OnClose" withInput:in onUI:NO];

  controller.context = nil;
  controller.profile = nil;
  [driver.elements removeObjectForKey:ID];

  [driver.macRPC return:returnID withOutput:nil andError:nil];
}

+ (void)emitDirection:(NSString *)elemID
                input:(ControllerInput)input
                    x:(float)x
                    y:(float)y {
  Driver *driver = [Driver current];
  NSDictionary *in = @{
    @"ID" : elemID,
    @"Input" : @(input),
    @"X" : @(x),
    @"Y" : @(y),
  };

  [driver.goRPC call:@"controller.OnDirectionChange" withInput:in onUI:YES];
}

+ (void)emitButton:(NSString *)elemID
             input:(ControllerInput)input
             value:(float)value
           pressed:(BOOL)pressed {
  Driver *driver = [Driver current];

  NSDictionary *in = @{
    @"ID" : elemID,
    @"Input" : @(input),
    @"Value" : @(value),
    @"Pressed" : @(pressed),
  };

  [driver.goRPC call:@"controller.OnButtonPressed" withInput:in onUI:YES];
}

+ (void)listen:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];

    Controller *controller = driver.elements[ID];

    controller.context.controllerPausedHandler = ^(GCController *controller) {
      [driver.goRPC call:@"controller.OnPause" withInput:in onUI:YES];
    };

    controller.profile.dpad.valueChangedHandler =
        ^(GCControllerDirectionPad *dpad, float x, float y) {
          [Controller emitDirection:ID input:DirectionalPad x:x y:y];
        };
    controller.profile.leftThumbstick.valueChangedHandler =
        ^(GCControllerDirectionPad *dpad, float x, float y) {
          [Controller emitDirection:ID input:LeftThumbstick x:x y:y];
        };
    controller.profile.rightThumbstick.valueChangedHandler =
        ^(GCControllerDirectionPad *dpad, float x, float y) {
          [Controller emitDirection:ID input:RightThumbstick x:x y:y];
        };

#if MAC_OS_X_VERSION_MIN_ALLOWED > 101400
    controller.profile.leftThumbstickButton.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID
                           input:LeftThumbstick
                           value:value
                         pressed:pressed];
        };
    controller.profile.rightThumbstickButton.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID
                           input:RightThumbstick
                           value:value
                         pressed:pressed];
        };
#endif

    controller.profile.buttonA.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID input:A value:value pressed:pressed];
        };
    controller.profile.buttonB.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID input:B value:value pressed:pressed];
        };
    controller.profile.buttonX.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID input:X value:value pressed:pressed];
        };
    controller.profile.buttonY.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID input:Y value:value pressed:pressed];
        };
    controller.profile.leftShoulder.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID input:L1 value:value pressed:pressed];
        };
    controller.profile.leftTrigger.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID input:L2 value:value pressed:pressed];
        };
    controller.profile.rightShoulder.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID input:R1 value:value pressed:pressed];
        };
    controller.profile.rightTrigger.pressedChangedHandler =
        ^(GCControllerButtonInput *button, float value, BOOL pressed) {
          [Controller emitButton:ID input:R2 value:value pressed:pressed];
        };

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}
@end