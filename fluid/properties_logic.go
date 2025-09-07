package fluid

// Property codes from SEUIF97
const (
	PRESSURE          = 0
	TEMPERATURE       = 1
	SPECIFIC_VOLUME   = 3
	ENTHALPY          = 4
	ENTROPY           = 5
	DYNAMIC_VISCOSITY = 24
)

func CalculatePressureHs(EnthalpyKJKG float64, EntropyKJ float64) (Pressure float64) {
	var pressure float64 = cgoHs(EnthalpyKJKG, EntropyKJ, PRESSURE) // MPa
	return pressure
}

func CalculateTemperatureHs(EnthalpyKJKG float64, EntropyKJ float64) (Temperature float64) {
	var temperature float64 = cgoHs(EnthalpyKJKG, EntropyKJ, TEMPERATURE) // Celsius
	return temperature
}

func CalculateTemperaturePh(PressureMPa float64, EnthalpyKJKG float64) float64 {
	var temperature float64 = cgoPh(PressureMPa, EnthalpyKJKG, TEMPERATURE) // Celsius
	return temperature
}

func CalculateTemperaturePs(PressureMPa float64, EntropyKJ float64) float64 {
	var temperature float64 = cgoPs(PressureMPa, EntropyKJ, TEMPERATURE) // Celsius
	return temperature
}

func CalculateSpecificVolumePh(PressureMPa float64, EnthalpyKJKG float64) float64 {
	var volume float64 = cgoPh(PressureMPa, EnthalpyKJKG, SPECIFIC_VOLUME) // m^3/kg
	return volume
}

func CalculateSpecificVolumePs(PressureMPa float64, EntropyKJ float64) float64 {
	var volume float64 = cgoPs(PressureMPa, EntropyKJ, SPECIFIC_VOLUME) // m^3/kg
	return volume
}

func CalculateSpecificVolumePt(PressureMPa float64, TemperatureC float64) float64 {
	var volume float64 = cgoPt(PressureMPa, TemperatureC, SPECIFIC_VOLUME) // m^3/kg
	return volume
}

func CalculateDensityPt(PressureMPa float64, TemperatureC float64) float64 {
	var specificVolumeM3kg = cgoPt(PressureMPa, TemperatureC, SPECIFIC_VOLUME)
	return 1.0 / float64(specificVolumeM3kg) // density = 1/specific_volume. kg/m^3
}

func CalculateDensityPh(PressureMPa float64, EnthalpyKJKG float64) float64 {
	var specificVolumeM3kg = cgoPh(PressureMPa, EnthalpyKJKG, SPECIFIC_VOLUME)
	return 1.0 / float64(specificVolumeM3kg)
}

func CalculateMass(DensityKGM3 float64, VolumeM3 float64) (Mass float64) { // custom function based on density and specific volume
	return DensityKGM3 * VolumeM3 //kg
}

func CalculateEnthalpyPt(PressureMPa float64, TemperatureC float64) (Enthalpy float64) {
	var EnthalpyKJKG = cgoPt(PressureMPa, TemperatureC, ENTHALPY)
	return EnthalpyKJKG // kJ/kg
}

func CalculateEnthalpyPs(PressureMPa float64, EntropyKJ float64) (Enthalpy float64) {
	var EnthalpyKJKG = cgoPs(PressureMPa, EntropyKJ, ENTHALPY)
	return EnthalpyKJKG // kJ/kg
}

func CalculateEntropyPt(PressureMPa float64, TemperatureC float64) float64 {
	var entropy float64 = cgoPt(PressureMPa, TemperatureC, ENTROPY)
	return entropy
}

func CalculateDynamicViscosityPt(PressureMPa float64, TemperatureC float64) float64 {
	var viscosity float64 = cgoPt(PressureMPa, TemperatureC, DYNAMIC_VISCOSITY) // kg/(m·s)
	return viscosity
}
