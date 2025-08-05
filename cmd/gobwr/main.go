package main

import (
	"fmt"
	"time"

	"github.com/checkm4te8/gobwr/fluid"
	"github.com/checkm4te8/gobwr/reactor"
)

// --- CONSTANT DECLARATIONS ---
const eventLoopFrequency time.Duration = 1 * time.Second // How often the physics are recalculated. By default, 1 second.

// --- MAIN EVENT LOOP ---
func main() {
	nodes := fluid.Init()
	state := reactor.Init()

	for {
		reactor.SimulateFission(&state)
		fmt.Println(nodes)
		fmt.Println(fluid.GetReactorWaterLevel(nodes))
		fmt.Println(reactor.CalculateThermalPower(state.CurrentNeutrons))
		time.Sleep(eventLoopFrequency) // Execute the main event loop every eventLoopFrequency seconds.
	}
}
