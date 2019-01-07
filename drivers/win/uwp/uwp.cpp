#include "uwp.hpp"
#include "_cgo_export.h"

void winCallReturn(char *retID, char *ret, char *err)
{
  onWinCallReturn(retID, ret, err);
}

void goCall(char *in) { return onGoCall(in); }