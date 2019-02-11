package main

/*
#cgo pkg-config: bash

#include "builtins.h"

extern int import_builtin(WORD_LIST *list);

struct builtin import_struct = {
  NULL,        // builtin name
  import_builtin,  // function implementing the builtin
  BUILTIN_ENABLED, // initial flags for builtin
  NULL,            // array of long documentation strings.
  NULL,            // usage synopsis; becomes short_doc
  0                // reserved for internal use
};
*/
import "C"

func main() {}
