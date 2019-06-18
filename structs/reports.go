
package structs


type NodeAverageAnalytics map[string]NodeAverageReport

type NodeAverageReport struct {
    TimeSlice float64     `json:"timeslice,omitempty"`
    CpuUsed   float64     `json:"cpu_used,omitempty"`
    MemUsed   float64     `json:"mem_used,omitempty"`
}


type ProcessAverageReport struct {
    TimeSlice    float64     `json:"timeslice,omitempty"`
    CpuUsed      float64     `json:"cpu_used,omitempty"`
    MemUsed      float64     `json:"mem_used,omitempty"`
    NumInstances float64     `json:"num_instances,omitempty"`
}


type ProcessHistoryReport struct {
    TimeSlice  float64        `json:"timeslice,omitempty"`
    Processes  []ProcessInfo  `json:"processes,omitempty"`
}


type ProcessInfo struct {
    Name string  `json:"name,omitempty"`
    Url  string  `json:"url,omitempty"`
}


