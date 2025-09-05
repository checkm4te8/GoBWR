//go:build windows

package fluid

/*
#cgo CFLAGS: -I${SRCDIR}/../lib/windows_amd64
#cgo LDFLAGS: -L${SRCDIR}/../lib/windows_amd64 -lseuif97 -lm
#include "seuif97.h"
*/
import "C"

// cgoPt is the Windows-specific wrapper that calls the C.pt function.
func cgoPt(pressure float64, temperature float64, propertyID int) float64 {
	return float64(C.pt(C.double(pressure), C.double(temperature), C.int(propertyID)))
}
