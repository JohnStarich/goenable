package main

/*
#cgo pkg-config: bash

#include "builtins.h"

extern int goenable_builtin(WORD_LIST *list);

char *empty_doc[] = {
	NULL
};

struct builtin goenable_struct = {
  "goenable",       // builtin name
  goenable_builtin, // function implementing the builtin
  BUILTIN_ENABLED,  // initial flags for builtin
  empty_doc,        // array of long documentation strings.
  "Run 'goenable help' for help.",  // usage synopsis; becomes short_doc. Note: This constant is replaced at load time on Bash 4.4+.
  0                 // reserved for internal use
};
*/
import "C"
