package main

import (
	"GoBWR/fluid"
	"GoBWR/reactor"
	"fmt"
	"time"
)

// --- CONSTANT DECLARATIONS ---
const eventLoopFrequency time.Duration = 1 * time.Second // How often the physics are recalculated. By default, 1 second.

// --- MAIN EVENT LOOP ---
func main() {
	fluid.InitializeFluidNodes()
	reactor.SetupReactor()
	for {
		reactor.SimulateFission()
		fmt.Println(fluid.FluidNodes)
		fmt.Println(fluid.GetReactorWaterLevel())
		fmt.Println(reactor.CalculateThermalPower())
		time.Sleep(eventLoopFrequency) // Execute the main event loop every eventLoopFrequency seconds.
	}
}
