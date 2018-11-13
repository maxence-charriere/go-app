#include "driver.h"
#include "controller.h"

@implementation Controller

// Create a new Controller that will listen for and
// dispatch MFi controller events.
+ (void) new:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];

    // Create blank controller, track it for future connection
    Controller *controller = [[Controller alloc] init];
    controller.ID = ID;
    driver.elements[ID] = controller;

    // Register for controller connection notifications
    [[NSNotificationCenter defaultCenter] addObserver:controller selector:@selector(connected) name:GCControllerDidConnectNotification object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:controller selector:@selector(disconnected) name:GCControllerDidDisconnectNotification object:nil];

    // Check if controller was connected before application started
    if ([[GCController controllers] count] > 0) {
        [controller connected];
    }

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}

// Called when a controller becomes connected
- (void) connected {
    Driver *driver = [Driver current];
    self.context = [GCController controllers][0];
    self.profile = self.context.extendedGamepad;
    NSDictionary *in = @{
        @"ID": self.ID,
    };

    [driver.goRPC call:@"controller.OnConnected" withInput:in onUI:YES];
}

// Called when a controller becomes disconnected
- (void) disconnected {
    Driver *driver = [Driver current];
    NSDictionary *in = @{
        @"ID": self.ID,
    };

    [driver.goRPC call:@"controller.OnDisconnected" withInput:in onUI:YES];
}

// Close removes the controller from the elements
+ (void) close:(NSDictionary *)in return:(NSString *)returnID {
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];
    Controller *controller = driver.elements[ID];

    if (controller == nil) {
      [NSException raise:@"ErrNoController" format:@"no controller with id %@", ID];
    }
    
    [driver.goRPC call:@"controller.OnClose" withInput:in onUI:NO];

    controller.context = nil;
    controller.profile = nil;
    [driver.elements removeObjectForKey:ID];

    [driver.macRPC return:returnID withOutput:nil andError:nil];
}

// General, shorthand function for emitting controller.onButtonChange
+ (void) emitButton:(NSString *)ID Input:(int)input Value:(float)value Pressed:(BOOL)pressed {
    Driver *driver = [Driver current];
    NSDictionary *in = @{
        @"ID": ID,
        @"Input": @(input),
        @"Value": @(value),
        @"Pressed": @(pressed),
    };

    [driver.goRPC call:@"controller.OnButtonChange" withInput:in onUI:YES];
}

// General, shorthand function for emitting controller.onDpadChange
+ (void) emitDpad:(NSString *)ID Input:(int)input X:(float)x Y:(float)y {
    Driver *driver = [Driver current];
    NSDictionary *in = @{
        @"ID": ID,
        @"Input": @(input),
        @"X": @(x),
        @"Y": @(y),
    };

    [driver.goRPC call:@"controller.OnDpadChange" withInput:in onUI:YES];
}

// Sets up handler functions for the controller object
+ (void) listen:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    Driver *driver = [Driver current];
    NSString *ID = in[@"ID"];

    Controller *controller = driver.elements[ID];

    controller.context.controllerPausedHandler = ^(GCController *controller)
    {
        [driver.goRPC call:@"controller.OnPause" withInput:in onUI:YES];
    };

    controller.profile.dpad.valueChangedHandler = ^(GCControllerDirectionPad *dpad, float x, float y) {
        [Controller emitDpad:ID Input:DirectionalPad X:x Y:y];
    };
    controller.profile.leftThumbstick.valueChangedHandler = ^(GCControllerDirectionPad *dpad, float x, float y) {
        [Controller emitDpad:ID Input:LeftThumbstick X:x Y:y];
    };
    controller.profile.rightThumbstick.valueChangedHandler = ^(GCControllerDirectionPad *dpad, float x, float y) {
        [Controller emitDpad:ID Input:RightThumbstick X:x Y:y];
    };
    if (@available(macOS 10.14.1, *)) {
        controller.profile.leftThumbstickButton.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
            [Controller emitButton:ID Input:LeftThumbstickButton Value:value Pressed:pressed];
        };
        controller.profile.rightThumbstickButton.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
            [Controller emitButton:ID Input:RightThumbstickButton Value:value Pressed:pressed];
        };
    }
    controller.profile.buttonA.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
        [Controller emitButton:ID Input:A Value:value Pressed:pressed];
    };
    controller.profile.buttonB.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
        [Controller emitButton:ID Input:B Value:value Pressed:pressed];
    };
    controller.profile.buttonX.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
        [Controller emitButton:ID Input:X Value:value Pressed:pressed];
    };
    controller.profile.buttonY.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
        [Controller emitButton:ID Input:Y Value:value Pressed:pressed];
    };
    controller.profile.leftShoulder.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
        [Controller emitButton:ID Input:L1 Value:value Pressed:pressed];
    };
    controller.profile.leftTrigger.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
        [Controller emitButton:ID Input:L2 Value:value Pressed:pressed];
    };
    controller.profile.rightShoulder.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
        [Controller emitButton:ID Input:R1 Value:value Pressed:pressed];
    };
    controller.profile.rightTrigger.pressedChangedHandler = ^(GCControllerButtonInput *button, float value, BOOL pressed) {
        [Controller emitButton:ID Input:R2 Value:value Pressed:pressed];
    };

    [driver.macRPC return:returnID withOutput:nil andError:nil];
  });
}
@end