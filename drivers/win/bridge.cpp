#include "bridge.hpp"
#include "_cgo_export.h"

void winReturn(char *retID, char *ret, char *err) {
  winCallReturn(retID, ret, err);
}
