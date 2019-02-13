package goenable

/*
#include "builtins.h"

extern struct builtin goenable_struct;
*/
import "C"

import (
	"strings"

	"github.com/johnstarich/bash-go-loader/cutils"
)

func init() {
	C.goenable_struct.name = C.CString(Name())
	longDoc := strings.Split(Usage(), "\n")
	C.goenable_struct.long_doc = (**C.char)(cutils.CStringArray(longDoc))
	C.goenable_struct.short_doc = C.CString(UsageShort())
}

//export goenable_builtin
func goenable_builtin(list *C.WORD_LIST) C.int {
	args := make([]string, 0)
	for list != nil {
		args = append(args, C.GoString(list.word.word))
		list = list.next
	}
	return C.int(Run(args))
}

//export goenable_builtin_load
func goenable_builtin_load(cName *C.char) C.int {
	name := C.GoString(cName)
	return C.int(Load(name))
}

//export goenable_builtin_unload
func goenable_builtin_unload() {
	Unload()
}
