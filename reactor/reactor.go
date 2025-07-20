package reactor

// --- CONSTANT DECLARATIONS ---
const MaxNeutrons int64 = 100000000000 // The amount of neutrons in the reactor core at 100% thermal power.

// --- VARIABLE DECLARATIONS ---
var Neutrons int64 = 1000 // The amount of neutrons currently in the reactor core. Reactor needs some starting neutrons to start.

func CalculateThermalPower() float64 {
	return float64(Neutrons) / float64(MaxNeutrons)
}
