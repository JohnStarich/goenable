package cutils

import "C"
import "unsafe"

var (
	pointerSize = unsafe.Sizeof(uintptr(0))
)

func stringArrayWithOffset(array unsafe.Pointer, offset int) **C.char {
	pointerAddress := uintptr(array) + uintptr(offset)*pointerSize
	return (**C.char)(unsafe.Pointer(pointerAddress))
}

// CStringArray converts a string slice to a **char, i.e. a C string array
// This creates a malloc'd array of malloc'd strings, so be sure to free all of them.
// Note: You will need to manually type cast this to (**C.char)(...)
func CStringArray(lines []string) unsafe.Pointer {
	arrayLen := int(pointerSize) * (len(lines) + 1)
	array := C.malloc(C.ulong(arrayLen))
	for i, line := range lines {
		arrayPtr := stringArrayWithOffset(array, i)
		*arrayPtr = C.CString(line)
	}
	lastLoc := stringArrayWithOffset(array, len(lines))
	*lastLoc = nil
	return array
}
