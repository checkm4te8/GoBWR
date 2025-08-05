package reactor

import "math"

// --- CONSTANT DECLARATIONS ---
const MaxNeutrons int64 = 100000000000 // The amount of neutrons in the reactor core at 100% thermal power.

// --- STRUCT DECLARATIONS ---
type Reactor struct {
	RodsPulled      float64
	CurrentNeutrons int64
	OldNeutrons     int64 // The amount of neutrons in the previous event loop iteration.
	IdleNeutrons    int64 // Reactor needs some starting neutrons to start.
}

// --- STRUCT INITIALIZATIONS ---
func Init() Reactor {
	return Reactor{
		RodsPulled: 0.5,
	}
}

func SimulateFission(reactor *Reactor) {
	var FissionFactors float64 = 2 * (0.2 + reactor.RodsPulled*0.8)
	reactor.OldNeutrons = reactor.CurrentNeutrons
	reactor.CurrentNeutrons = int64(math.Round(FissionFactors * (float64(reactor.OldNeutrons) + float64(reactor.IdleNeutrons))))
	if reactor.CurrentNeutrons > MaxNeutrons {
		reactor.RodsPulled = 0 //scram
	}
}

func CalculateThermalPower(currentNeutrons int64) float64 {
	return float64(currentNeutrons) / float64(MaxNeutrons)
}
