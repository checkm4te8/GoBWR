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

// cgoPh is the macOS-specific wrapper that calls the C.ph function.
func cgoPh(pressure float64, enthalpy float64, propertyID int) float64 {
	return float64(C.ph(C.double(pressure), C.double(enthalpy), C.int(propertyID)))
}

// cgoPs is the macOS-specific wrapper that calls the C.ph function.
func cgoPs(pressure float64, entropy float64, propertyID int) float64 {
	return float64(C.ps(C.double(pressure), C.double(entropy), C.int(propertyID)))
}

// cgoHs is the macOS-specific wrapper that calls the C.ph function.
func cgoHs(enthalpy float64, entropy float64, propertyID int) float64 {
	return float64(C.hs(C.double(enthalpy), C.double(entropy), C.int(propertyID)))
}
