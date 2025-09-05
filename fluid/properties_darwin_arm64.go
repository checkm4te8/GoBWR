//go:build darwin && arm64

package fluid

/*
#cgo CFLAGS: -I${SRCDIR}/../lib/darwin_arm64
#cgo LDFLAGS: -L${SRCDIR}/../lib/darwin_arm64 -lseuif97 -lm
#include "seuif97.h"
*/
import "C"

// cgoPt is a macOS-specific wrapper that calls the C.pt function.
func cgoPt(pressure float64, temperature float64, propertyID int) float64 {
	return float64(C.pt(C.double(pressure), C.double(temperature), C.int(propertyID)))
}
