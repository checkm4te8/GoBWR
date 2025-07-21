package fluid

// --- STRUCT DECLARATIONS ---
type FluidNode struct {
	Temperature float64 // overall temperature of the fluid node in degrees Celsius
	Pressure    float64 // overall pressure of the fluid node in Pa
	Volume      float64 // volume of the fluid node in cubic meters
	Mass        float64 // total mass of the fluid in kg
	Enthalpy    float64 // J/kg
}

// --- CONSTANT DECLARATIONS ---

// --- VARIABLE DECLARATIONS ---
var FluidNodes map[string]FluidNode = map[string]FluidNode{
	"FeedwaterHeader": FluidNode{
		20, 101325, 10, 0, 0, // Enthalpy is calculated by IAPWS IF97 upon initialization.
	},
	"Downcomer": FluidNode{
		25, 200000, 50, 50000, 0,
	},
	"LowerPlenum": FluidNode{
		30, 250000, 20, 20000, 0,
	},
	"CoreBottom": FluidNode{
		35, 240000, 15, 15000, 0,
	},
	"CoreMiddle": FluidNode{
		40, 230000, 15, 15000, 0,
	},
	"CoreTop": FluidNode{
		45, 220000, 15, 15000, 0,
	},
	"SteamDome": FluidNode{
		50, 200000, 100, 5000, 0,
	},
}

var FlowRates map[string]float64 = map[string]float64{
	"FeedwaterHeader->Downcomer": 0.0,
	"Downcomer->LowerPlenum":     0.0,
	"LowerPlenum->CoreBottom":    0.0,
	"CoreBottom->CoreMiddle":     0.0,
	"CoreMiddle->CoreTop":        0.0,
	"CoreTop->SteamDome":         0.0,
}

func InitializeFluidNodes() {
	for name, node := range FluidNodes {
		node.Enthalpy = CalculateEnthalpy(node.Temperature, node.Pressure)
		FluidNodes[name] = node
	}
}

func GetNodeDensity(InputNode FluidNode) float64 {
	return CalculateDensity(InputNode.Temperature, InputNode.Pressure)
}
