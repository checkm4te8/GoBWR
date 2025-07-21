package reactor

import "math"

// --- CONSTANT DECLARATIONS ---
const MaxNeutrons int64 = 100000000000 // The amount of neutrons in the reactor core at 100% thermal power.

// --- VARIABLE DECLARATIONS ---
var CurrentNeutrons int64 = 0
var IdleNeutrons int64 = 1000 // Reactor needs some starting neutrons to start.
var OldNeutrons int64 = 0     // The amount of neutrons in the previous event loop iteration.

// --- STRUCT DECLARATIONS ---
type Reactor struct {
	RodsPulled float64
}

// --- STRUCT INITIALIZATIONS ---
var ReactorState Reactor

func SetupReactor() {
	ReactorState.RodsPulled = 0.5
}

func SimulateFission() {
	var FissionFactors float64 = 2 * (0.2 + ReactorState.RodsPulled*0.8)
	OldNeutrons = CurrentNeutrons
	CurrentNeutrons = int64(math.Round(FissionFactors * (float64(OldNeutrons) + float64(IdleNeutrons))))
	if CurrentNeutrons > MaxNeutrons {
		ReactorState.RodsPulled = 0 //scram
	}
}

func CalculateThermalPower() float64 {
	return float64(CurrentNeutrons) / float64(MaxNeutrons)
}
