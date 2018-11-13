
#ifndef bridge_h
#define bridge_h

#ifdef __cplusplus

extern "C" {

#endif

void winCallReturn(char *retID, char *ret, char *err);
char *goCall(char *in, char *ui);

#ifdef __cplusplus
}

#endif
#endif /* bridge_h */