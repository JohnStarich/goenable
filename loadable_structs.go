package main

/*
#cgo pkg-config: bash

#include "builtins.h"

extern int goenable_builtin(WORD_LIST *list);

struct builtin goenable_struct = {
  NULL,             // builtin name
  goenable_builtin, // function implementing the builtin
  BUILTIN_ENABLED,  // initial flags for builtin
  NULL,             // array of long documentation strings.
  NULL,             // usage synopsis; becomes short_doc
  0                 // reserved for internal use
};
*/
import "C"
