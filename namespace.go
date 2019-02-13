package main

/*
#include "builtins.h"

extern struct builtin namespace_struct;
*/
import "C"

import (
	"strings"
	"unsafe"

	"github.com/johnstarich/bash-go-loader/namespace"
)

var (
	pointerSize = unsafe.Sizeof(uintptr(0))
)

func init() {
	C.namespace_struct.name = C.CString(namespace.Name())
	longDoc := strings.Split(namespace.Usage(), "\n")
	C.namespace_struct.long_doc = cStringArray(longDoc)
	C.namespace_struct.short_doc = C.CString(namespace.UsageShort())
}

func stringArrayWithOffset(array unsafe.Pointer, offset int) **C.char {
	pointerAddress := uintptr(array) + uintptr(offset)*pointerSize
	return (**C.char)(unsafe.Pointer(pointerAddress))
}

// cStringArray converts a string slice to a **char, i.e. a C string array
// This creates a malloc'd array of malloc'd strings, so be sure to free all of them.
func cStringArray(lines []string) **C.char {
	arrayLen := int(pointerSize) * (len(lines) + 1)
	array := C.malloc(C.ulong(arrayLen))
	for i, line := range lines {
		arrayPtr := stringArrayWithOffset(array, i)
		*arrayPtr = C.CString(line)
	}
	lastLoc := stringArrayWithOffset(array, len(lines))
	*lastLoc = nil
	return (**C.char)(array)
}

//export namespace_builtin
func namespace_builtin(list *C.WORD_LIST) C.int {
	args := make([]string, 0)
	for list != nil {
		args = append(args, C.GoString(list.word.word))
		list = list.next
	}
	return C.int(namespace.Run(args))
}

//export namespace_builtin_load
func namespace_builtin_load(cName *C.char) C.int {
	name := C.GoString(cName)
	return C.int(namespace.Load(name))
}

//export namespace_builtin_unload
func namespace_builtin_unload() {
	namespace.Unload()
}
