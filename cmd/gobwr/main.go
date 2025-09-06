package main

import (
	"GoBWR/fluid"
	"GoBWR/reactor"
	"fmt"
	"time"
)

// --- CONSTANT DECLARATIONS ---
const deltaTime time.Duration = 1 * time.Second / 10 // How often the physics are recalculated.

// --- MAIN EVENT LOOP ---
func main() {
	fluid.InitializeFluidNodes()
	reactor.SetupReactor()
	for {
		fluid.SimulateFlow(deltaTime)
		reactor.SimulateFission()
		fmt.Println(fluid.FluidNodes)
		time.Sleep(deltaTime) // Execute the main event loop every deltaTime.
	}
}
