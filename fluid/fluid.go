package fluid

import (
	"errors"
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
	MinorKFactor float64           // the minor K-Factor caused by things like fittings, elbows... in the piping. The major K-Factor is calculated when simulating flow.
}

type FlowPath struct {
	SourceNodeID      string
	DestinationNodeID string
	JunctionIDs       []string
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
		450,
		30,
		3.5,
	},
	"HotwellToTest2": FluidPipe{
		FluidJunctionBase{
			"Junction",
			"HotwellToTest",
			"Junction",
			"HotwellToTest3",
		},
		400,
		10,
		4,
	},
	"HotwellToTest3": FluidPipe{
		FluidJunctionBase{
			"Junction",
			"HotwellToTest2",
			"Node",
			"ReactorVessel",
		},
		550,
		80,
		2.5,
	},
}

var FlowPaths []FlowPath // Will be initialized automatically

func FindConnectionToJunction(junctionId string) (nextType string, nextId string, searchError error) {
	var junction FluidPipe
	var ok bool
	junction, ok = FluidPipes[junctionId]

	if !ok {
		return "", "", errors.New("junction not found")
	}

	if junction.JunctionBase.DestinationType == "Junction" {
		var ok bool
		_, ok = FluidPipes[junction.JunctionBase.DestinationID]
		if !ok {
			return "", "", errors.New("destination junction does not exist")
		}
	}

	if junction.JunctionBase.DestinationType == "Node" {
		var ok bool
		_, ok = FluidNodes[junction.JunctionBase.DestinationID]
		if !ok {
			return "", "", errors.New("destination node does not exist")
		}
	}

	return junction.JunctionBase.DestinationType, junction.JunctionBase.DestinationID, nil
}

func GetJunctionPathToDestination(startJunctionId string) (junctionPath []string, destinationNodeId string, err error) {
	var currentJunctionId string = startJunctionId
	junctionPath = append(junctionPath, currentJunctionId)
	for {
		nextStepType, nextStepId, searchError := FindConnectionToJunction(currentJunctionId)
		if searchError != nil {
			return junctionPath, currentJunctionId, searchError
		}
		if nextStepType == "Node" {
			return junctionPath, nextStepId, nil
		}
		junctionPath = append(junctionPath, nextStepId)
		currentJunctionId = nextStepId
	}
}

func InitializeFluidNodes() {
	for name, node := range FluidNodes { // Initialize Enthalpy and Mass of all nodes
		node.Enthalpy = CalculateEnthalpy(node.Temperature, node.Pressure)
		node.Mass = CalculateMass(node.Volume, CalculateDensity(node.Temperature, node.Pressure))
		FluidNodes[name] = node
	}

	for pipeName, pipe := range FluidPipes { // Initialize flow paths, fluid will only flow if connected to a junction directly. Never from one node to another.
		if pipe.JunctionBase.SourceType == "Node" && pipe.JunctionBase.DestinationType == "Junction" {
			var path, destination, err = GetJunctionPathToDestination(pipeName)
			var flowPath FlowPath = FlowPath{
				pipe.JunctionBase.SourceID,
				destination,
				path,
			}
			if err == nil {
				FlowPaths = append(FlowPaths, flowPath)
			}
		}
	}
}

func CalculateTotalPipeKAndVelocityMap(flowPath FlowPath, kPipeMap map[string]float64, pressureMagnitude float64, sourceNodeDensity float64) (normalizedTotalK float64, pipeVelocityMap map[string]float64) {
	pipeVelocityMap = make(map[string]float64)
	var firstIteration bool = true // get the first pipe of the path to use as a reference for total K-factor
	var firstPipe FluidPipe
	for _, pipeId := range flowPath.JunctionIDs { // calculate total normalized K-Factor
		if firstIteration {
			normalizedTotalK += kPipeMap[pipeId]
			firstPipe = FluidPipes[pipeId]
			firstIteration = false
			continue
		}
		var firstPipeA float64 = math.Pi * math.Pow((firstPipe.PipeDiameter/1000)/2, 2) // cross-sectional area of first and current pipe
		var currentPipeA float64 = math.Pi * math.Pow((FluidPipes[pipeId].PipeDiameter/1000)/2, 2)
		normalizedTotalK += kPipeMap[pipeId] * math.Pow(firstPipeA/currentPipeA, 2)
	}

	var firstPipeVelocity float64 = math.Sqrt(2 * pressureMagnitude / (normalizedTotalK * sourceNodeDensity))
	for _, pipeId := range flowPath.JunctionIDs { // get velocity in all pipes
		var firstPipeA float64 = math.Pi * math.Pow((firstPipe.PipeDiameter/1000)/2, 2) // cross-sectional area of first and current pipe
		var currentPipeA float64 = math.Pi * math.Pow((FluidPipes[pipeId].PipeDiameter/1000)/2, 2)
		var currentPipeVelocity float64 = firstPipeVelocity * (firstPipeA / currentPipeA) // principle of continuity to derive current pipe velocity from the first (reference) pipe
		pipeVelocityMap[pipeId] = currentPipeVelocity
	}
	return
}

func CalculateKPipeMapAndFrictionFactorMap(flowPath FlowPath, previousFrictionFactors map[string]float64, pressureMagnitude float64, sourceNode FluidNode) (kPipeMap map[string]float64, pipeFrictionFactorMap map[string]float64) {
	kPipeMap = make(map[string]float64)
	var sourceNodeDensity = CalculateDensity(sourceNode.Temperature, sourceNode.Pressure)
	for _, pipeId := range flowPath.JunctionIDs { // populate kPipeMap
		kPipeMap[pipeId] = (previousFrictionFactors[pipeId] * (FluidPipes[pipeId].PipeLength) / (FluidPipes[pipeId].PipeDiameter / 1000)) + FluidPipes[pipeId].MinorKFactor // diameter unit conversion mm->m
	}

	var _, pipeVelocityMap = CalculateTotalPipeKAndVelocityMap(flowPath, kPipeMap, pressureMagnitude, sourceNodeDensity)

	pipeFrictionFactorMap = make(map[string]float64)
	for _, pipeId := range flowPath.JunctionIDs { // get a better friction factor estimate
		var reynoldsNumber float64 = (sourceNodeDensity * pipeVelocityMap[pipeId] * (FluidPipes[pipeId].PipeDiameter / 1000)) / CalculateDynamicViscosity(sourceNode.Temperature, sourceNode.Pressure)
		var pipeAbsoluteRoughness float64 = 0.000045 // meters, this is an estimate
		var swameeJainFrictionFactor float64 = 0.25 / math.Pow(math.Log10((pipeAbsoluteRoughness/(3.7*(FluidPipes[pipeId].PipeDiameter/1000)))+(5.74/math.Pow(reynoldsNumber, 0.9))), 2)
		pipeFrictionFactorMap[pipeId] = swameeJainFrictionFactor
	}
	return
}

func SimulateFlow() {
	for _, flowPath := range FlowPaths {
		fmt.Println(flowPath)
		var sourceNode FluidNode = FluidNodes[flowPath.SourceNodeID]
		var destinationNode FluidNode = FluidNodes[flowPath.DestinationNodeID]
		var deltaP float64 = sourceNode.Pressure - destinationNode.Pressure // positive means flow from source to destination, 0 means no flow, negative means flow from destination to source
		var pressureMagnitude float64 = math.Abs(deltaP)
		var fGuessMap map[string]float64 = make(map[string]float64) // find the darcy friction factor using an interative loop to get the major K-Factor
		for _, pipeName := range flowPath.JunctionIDs {
			fGuessMap[pipeName] = 0.02
		}
		var pipeFrictionFactorMap map[string]float64 = make(map[string]float64)
		_, pipeFrictionFactorMap = CalculateKPipeMapAndFrictionFactorMap(flowPath, fGuessMap, pressureMagnitude, sourceNode)
		for i := 0; i < 2; i += 1 { // 3 iterative recalculations for higher accuracy
			_, pipeFrictionFactorMap = CalculateKPipeMapAndFrictionFactorMap(flowPath, pipeFrictionFactorMap, pressureMagnitude, sourceNode)
		}
		var kMap, _ = CalculateKPipeMapAndFrictionFactorMap(flowPath, pipeFrictionFactorMap, pressureMagnitude, sourceNode)
		var totalFinalK float64 = 0.0
		var firstIteration bool = true
		var firstPipe FluidPipe
		for _, pipeId := range flowPath.JunctionIDs {
			if firstIteration {
				totalFinalK += kMap[pipeId]
				firstPipe = FluidPipes[pipeId]
				firstIteration = false
				fmt.Println("first " + pipeId + " is")
				fmt.Println(kMap[pipeId])
				continue
			}
			var firstPipeA float64 = math.Pi * math.Pow((firstPipe.PipeDiameter/1000)/2, 2) // cross-sectional area of first and current pipe
			var currentPipeA float64 = math.Pi * math.Pow((FluidPipes[pipeId].PipeDiameter/1000)/2, 2)
			var normalizedK float64 = kMap[pipeId] * math.Pow(firstPipeA/currentPipeA, 2)
			totalFinalK += normalizedK
			fmt.Println(pipeId + " is")
			fmt.Println(kMap[pipeId])
		}
		fmt.Println(totalFinalK)
	}
}

func GetReactorWaterLevel() float64 {
	var RPVNode FluidNode = FluidNodes["ReactorVessel"]
	var density float64 = CalculateDensity(RPVNode.Temperature, RPVNode.Pressure)
	var currentWaterVolume float64 = RPVNode.Mass / density
	return (currentWaterVolume / RPVNode.MaxVolume) * float64(RPVHeight)
}
