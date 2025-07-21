package fluid

// --- STRUCT DECLARATIONS ---
type FluidNode struct {
	Temperature float64 // overall temperature of the fluid node in degrees Celsius
	Pressure    float64 // overall pressure of the fluid node in Pa
	Volume      float64 // volume of the fluid node in cubic meters
	Mass        float64 // total mass of the fluid in kg
	Enthalpy    float64 // J/kg
	MaxVolume   float64 // max volume in cubic metres. Fluid will not flow into the node if it is full.
}

// --- CONSTANT DECLARATIONS ---
const RPVHeight float32 = 21.3 // In meters, how high the water can fill in the RPV

// --- VARIABLE DECLARATIONS ---
var FluidNodes map[string]FluidNode = map[string]FluidNode{
	"FeedwaterHeader": FluidNode{
		20, 101325, 10, 0, 0, 999999, // Enthalpy is calculated by IAPWS IF97 upon initialization. Mass is calculated upon initialization as well.
	},
	"ReactorVessel": FluidNode{
		35, 230000, 623, 0, 0, 928,
	},
	"SteamDome": FluidNode{
		50, 200000, 100, 0, 0, 999999,
	},
}

var FlowRates map[string]float64 = map[string]float64{
	"FeedwaterHeader->ReactorVessel": 0.0,
}

func InitializeFluidNodes() {
	for name, node := range FluidNodes {
		node.Enthalpy = CalculateEnthalpy(node.Temperature, node.Pressure)
		node.Mass = CalculateMass(node.Volume, CalculateDensity(node.Temperature, node.Pressure))
		FluidNodes[name] = node
	}
}

func GetReactorWaterLevel() float64 {
	var RPVNode FluidNode = FluidNodes["ReactorVessel"]
	var density float64 = CalculateDensity(RPVNode.Temperature, RPVNode.Pressure)
	var currentWaterVolume float64 = RPVNode.Mass / density
	return (currentWaterVolume / RPVNode.MaxVolume) * float64(RPVHeight)
}
