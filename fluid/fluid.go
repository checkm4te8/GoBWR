package fluid

import (
	"fmt"
	"math"
)

// --- STRUCT DECLARATIONS ---
type FluidNode struct {
	Temperature float64 // overall temperature of the fluid node in degrees Celsius
	Pressure    float64 // overall pressure of the fluid node in Pa
	Volume      float64 // volume of the fluid node in cubic meters
	Mass        float64 // total mass of the fluid in kg
	Enthalpy    float64 // J/kg
	MaxVolume   float64 // max volume in cubic metres. Fluid will not flow into the node if it is full.
}

type FluidJunctionBase struct {
	SourceType      string // Node/Junction
	SourceID        string // e.g. FeedwaterHeader
	DestinationType string // Node/Junction
	DestinationID   string // e.g. ReactorVessel
}

type FluidPipe struct {
	JunctionBase FluidJunctionBase // the base junction info
	PipeDiameter float64           // diameter of the pipe in milimeters
	PipeLength   float64           // length of the pipe in meters
	KFactor      float64           // the total K-factor for pressure loss sim.
}

// --- CONSTANT DECLARATIONS ---
const RPVHeight float32 = 21.3 // In meters, how high the water can fill in the RPV

// --- VARIABLE DECLARATIONS ---
var FluidNodes map[string]FluidNode = map[string]FluidNode{
	"Hotwell": FluidNode{
		20, 101325, math.Inf(1), 0, 0, math.Inf(1), // Enthalpy is calculated by IAPWS IF97 upon initialization. Mass is calculated upon initialization as well.
	},
	"ReactorVessel": FluidNode{
		35, 230000, 623, 0, 0, 928,
	},
}

var FluidPipes map[string]FluidPipe = map[string]FluidPipe{
	"HotwellToTest": FluidPipe{
		FluidJunctionBase{
			"Node",
			"Hotwell",
			"Junction",
			"HotwellToTest2",
		},
		50,
		50,
		10,
	},
	"HotwellToTest2": FluidPipe{
		FluidJunctionBase{
			"Junction",
			"HotwellToTest",
			"Junction",
			"HotwellToTest3",
		},
		50,
		50,
		10,
	},
	"HotwellToTest3": FluidPipe{
		FluidJunctionBase{
			"Junction",
			"HotwellToTest2",
			"Node",
			"Hotwell",
		},
		50,
		50,
		10,
	},
}

func GetConnection(initialJunctionId string) {
	for pipeName, pipe := range FluidPipes {
		if pipe.JunctionBase.SourceType == "Junction" && pipe.JunctionBase.SourceID == initialJunctionId {
			fmt.Println(pipeName + " connects to " + pipe.JunctionBase.SourceID)
			if pipe.JunctionBase.DestinationType == "Junction" {
				// there's more connections, continue
				GetConnection(pipeName)
			} else {
				fmt.Println(pipeName + " ends connection chain at node " + pipe.JunctionBase.DestinationID)
			}
		}
	}
}

func BuildNodeFlowTree() { // Start from node, find the first junction that the node is a source to, and continue looping through all junctions until arrive at a node.
	for name, _ := range FluidNodes {
		for pipeName, pipe := range FluidPipes {
			if pipe.JunctionBase.SourceType == "Node" && pipe.JunctionBase.SourceID == name {
				GetConnection(pipeName)
			}
		}
	}
}

func InitializeFluidNodes() {
	for name, node := range FluidNodes {
		node.Enthalpy = CalculateEnthalpy(node.Temperature, node.Pressure)
		node.Mass = CalculateMass(node.Volume, CalculateDensity(node.Temperature, node.Pressure))
		FluidNodes[name] = node
	}
	BuildNodeFlowTree()
}

func GetReactorWaterLevel() float64 {
	var RPVNode FluidNode = FluidNodes["ReactorVessel"]
	var density float64 = CalculateDensity(RPVNode.Temperature, RPVNode.Pressure)
	var currentWaterVolume float64 = RPVNode.Mass / density
	return (currentWaterVolume / RPVNode.MaxVolume) * float64(RPVHeight)
}
