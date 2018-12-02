#include "bridge.hpp"
#include "_cgo_export.h"

void winCallReturn(char *retID, char *ret, char *err)
{
  onWinCallReturn(retID, ret, err);
}

char *goCall(char *in, char *ui)
{
  return onGoCall(in, ui);
}