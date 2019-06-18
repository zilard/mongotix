
package structs

type NodeMetricsMap map[string]NodeData

type ProcessMetricsMap map[string][]ProcessMeasurement


type NodeData struct {
    NodeMeasurementArray []NodeMeasurement
    ProcessMeasurementMap ProcessMetricsMap
}



type NodeMeasurement struct {
    TimeSlice float64     `json:"timeslice,omitempty" bson:"timeslice,omitempty"`
    Cpu       float64     `json:"cpu,omitempty" bson:"cpu,omitempty"`
    Mem       float64     `json:"mem,omitempty" bson:"mem,omitempty"`
}


type ProcessMeasurement struct {
    TimeSlice float64     `json:"timeslice,omitempty" bson:"timeslice,omitempty"`
    CpuUsed   float64     `json:"cpu_used,omitempty" bson:"cpu_used,omitempty"`
    MemUsed   float64     `json:"mem_used,omitempty" bson:"mem_used,omitempty"`
}


var nodeMeasurementArray []NodeMeasurement
var processMeasurementArray []ProcessMeasurement



type ProcessMetricsArray []ProcessMetricsByName

type ProcessMetricsByName struct {
    ProcessName string
    MetricsData ProcessMeasurement
}


