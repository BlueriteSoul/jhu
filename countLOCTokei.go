package main

/*
#cgo LDFLAGS: ./libs/libtokei.a -ldl -lm
#include <stdlib.h>

int run_tokei_with_args(int argc, char **argv);
*/
import "C"
import "unsafe"

// CountLOCTokei runs the embedded Tokei library with safe argument handling.
func CountLOCTokei(args []string) {
	// Always ensure there’s at least a program name arg.
	if len(args) == 0 {
		args = []string{"tokei", "."}
	} else {
		// Ensure tokei gets a directory arg if user didn’t specify one.
		hasDir := false
		for _, a := range args {
			if a == "." || a == ".." || a == "./" {
				hasDir = true
				break
			}
		}
		if !hasDir {
			args = append(args, ".")
		}
	}

	argc := C.int(len(args))
	cargv := make([]*C.char, len(args))
	for i, s := range args {
		cargv[i] = C.CString(s)
		defer C.free(unsafe.Pointer(cargv[i]))
	}

	C.run_tokei_with_args(argc, &cargv[0])
}
