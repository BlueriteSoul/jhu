package main

/*
#cgo LDFLAGS: ./libs/libtokei.a -ldl -lm

#include <stdint.h>
#include <stdlib.h>

// Declare the Rust functions you want to call
uint64_t count_loc(const char* path);
*/
import "C"
import "unsafe"

func CountLOCTokei(path string) uint64 {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath)) // free memory after use

	return uint64(C.count_loc(cPath))
}
