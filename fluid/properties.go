package fluid

/*
#cgo LDFLAGS: -L${SRCDIR}/../lib -lseuif97 -lm
#cgo CFLAGS: -I${SRCDIR}/../lib
#include "seuif97.h"
*/
import "C"

// Property codes from SEUIF97
const (
	SPECIFIC_VOLUME = 3
	ENTHALPY        = 4
	ENTROPY         = 5
)

func CalculateEnthalpy(Temperature float64, Pressure float64) (Enthalpy float64) {
	var PressureMPa float64 = Pressure / 1000000
	var EnthalpyKJ = C.pt(C.double(PressureMPa), C.double(Temperature), ENTHALPY)
	return float64(EnthalpyKJ * 1000) //return in J/kg
}

func CalculateDensity(Temperature float64, Pressure float64) float64 {
	var PressureMPa float64 = Pressure / 1000000
	var specificVolumeM3kg = C.pt(C.double(PressureMPa), C.double(Temperature), SPECIFIC_VOLUME)
	return 1.0 / float64(specificVolumeM3kg) // density = 1/specific_volume
}
