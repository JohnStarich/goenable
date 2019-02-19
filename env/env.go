package env

/*
#include <stdlib.h>
*/
import "C"

import (
	"os"
	"sync"
	"unsafe"
)

var envMutex = &sync.RWMutex{}

// Getenv is similar to os.Getenv, but *always* requests the variable from C's getenv
// Note: If you rely on this value being consistent, you must use this package's other env functions as well
func Getenv(key string) string {
	cKey := C.CString(key)
	envMutex.RLock()
	defer envMutex.RUnlock()
	cValue := C.getenv(cKey)
	value := C.GoString(cValue)
	C.free(unsafe.Pointer(cKey))
	return value
}

// Setenv is similar to os.Setenv, but ensures consistent results for env.Getenv
// Note: If you rely on this value being consistent, you must use this package's other env functions as well
func Setenv(key, value string) error {
	envMutex.Lock()
	defer envMutex.Unlock()
	return os.Setenv(key, value)
}

// Unsetenv is similar to os.Unsetenv, but ensures consistent results for env.Getenv
// Note: If you rely on this value being consistent, you must use this package's other env functions as well.
func Unsetenv(key string) error {
	cKey := C.CString(key)
	envMutex.Lock()
	defer envMutex.Unlock()
	_ = C.setenv(cKey, (*C.char)(C.NULL), 1)
	C.free(unsafe.Pointer(cKey))
	return os.Unsetenv(key)
}
