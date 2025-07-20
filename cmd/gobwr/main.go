package main

import (
	"GoBWR/reactor"
	"fmt"
	"time"
)

// --- CONSTANT DECLARATIONS ---
const eventLoopFrequency time.Duration = 1 * time.Second // How often the physics are recalculated. By default, 1 second.

// --- MAIN EVENT LOOP ---
func main() {
	for {
		fmt.Println(reactor.CalculateThermalPower())
		time.Sleep(eventLoopFrequency) // Execute the main event loop every eventLoopFrequency seconds.
	}
}
