#ifndef sandbox_h
#define sandbox_h
#import <Foundation/Foundation.h>

typedef enum {
  OBCodeSignStateUnsigned = 1,
  OBCodeSignStateSignatureValid,
  OBCodeSignStateSignatureInvalid,
  OBCodeSignStateSignatureNotVerifiable,
  OBCodeSignStateSignatureUnsupported,
  OBCodeSignStateError
} OBCodeSignState;

@interface NSBundle (OBCodeSigningInfo)
- (BOOL)isSandboxed;
- (BOOL)ob_comesFromAppStore;
- (BOOL)ob_isSandboxed;
- (OBCodeSignState)ob_codeSignState;
@end

@interface NSBundle (OBCodeSigningInfoPrivateMethods)
- (SecStaticCodeRef)ob_createStaticCode;
- (SecRequirementRef)ob_sandboxRequirement;
@end

#endif /* sandbox_h */
