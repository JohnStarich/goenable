package namespace

/*
#include "builtins.h"

extern struct builtin namespace_struct;
*/
import "C"

import (
	"strings"

	"github.com/johnstarich/bash-go-loader/cutils"
)

func init() {
	C.namespace_struct.name = C.CString(Name())
	longDoc := strings.Split(Usage(), "\n")
	C.namespace_struct.long_doc = (**C.char)(cutils.CStringArray(longDoc))
	C.namespace_struct.short_doc = C.CString(UsageShort())
}

//export namespace_builtin
func namespace_builtin(list *C.WORD_LIST) C.int {
	args := make([]string, 0)
	for list != nil {
		args = append(args, C.GoString(list.word.word))
		list = list.next
	}
	return C.int(Run(args))
}

//export namespace_builtin_load
func namespace_builtin_load(cName *C.char) C.int {
	name := C.GoString(cName)
	return C.int(Load(name))
}

//export namespace_builtin_unload
func namespace_builtin_unload() {
	Unload()
}
