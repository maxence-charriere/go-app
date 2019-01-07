
#ifndef uwp_h
#define uwp_h

#include <stdlib.h>

#ifdef __cplusplus

extern "C" {

#endif

void winCallReturn(char *retID, char *ret, char *err);
void goCall(char *in);

#ifdef __cplusplus
}

#endif
#endif /* uwp_h */