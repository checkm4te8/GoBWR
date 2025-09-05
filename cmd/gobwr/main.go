package main

import (
	"GoBWR/fluid"
	"GoBWR/reactor"
	"time"
)

// --- CONSTANT DECLARATIONS ---
const deltaTime time.Duration = 1 * time.Second / 10 // How often the physics are recalculated. By default, 1 second.

// --- MAIN EVENT LOOP ---
func main() {
	fluid.InitializeFluidNodes()
	reactor.SetupReactor()
	for {
		fluid.SimulateFlow()
		reactor.SimulateFission()
		time.Sleep(deltaTime) // Execute the main event loop every deltaTime.
	}
}
