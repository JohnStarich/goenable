package main

/*
#include "builtins.h"

extern struct builtin import_struct;
*/
import "C"

import (
	"strings"
	"unsafe"

	"./importer"
)

var (
	pointerSize = unsafe.Sizeof(uintptr(0))
)

func init() {
	C.import_struct.name = C.CString(importer.Name())
	longDoc := strings.Split(importer.Usage(), "\n")
	C.import_struct.long_doc = cStringArray(longDoc)
	C.import_struct.short_doc = C.CString(importer.UsageShort())
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

//export import_builtin
func import_builtin(list *C.WORD_LIST) C.int {
	args := make([]string, 0)
	for list != nil {
		args = append(args, C.GoString(list.word.word))
		list = list.next
	}
	return C.int(importer.Run(args))
}

//export import_builtin_load
func import_builtin_load(cName *C.char) C.int {
	name := C.GoString(cName)
	return C.int(importer.Load(name))
}

//export import_builtin_unload
func import_builtin_unload() {
	importer.Unload()
}
