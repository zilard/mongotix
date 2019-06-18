
package utils

import (
    "strconv"
    "math/rand"
    "time"

    s "github.com/zilard/mongotix/structs"
)




// CreateDummyNodeMetrics - creates dummy Node metrics and stores it in nodeMetricsMap
// used by Golden testing
func CreateDummyNodeMetrics(nodeMetricsMap s.NodeMetricsMap) {
    for i := 1; i <= 2; i++ {
        nodeData := s.NodeData{}
        for j := 1; j <= 10; j++ {
            nodeData.NodeMeasurementArray = append(nodeData.NodeMeasurementArray,
                         s.NodeMeasurement{
                             TimeSlice: float64(j * 10),
                             Cpu: float64(j * 5),
                             Mem: float64(j * 6),
                         })
        }
        nodeMetricsMap["n" + strconv.Itoa(i)] = nodeData
     }
}



// CreateDummyProcessMetrics - creates dummy Process metrics and stores it in nodeMetricsMap
// used by Golden testing
func CreateDummyProcessMetrics(nodeMetricsMap s.NodeMetricsMap) {

    for i := 1; i <= 2; i++ {
        nodeData := s.NodeData{}

        processMetricsMap := make(s.ProcessMetricsMap)

        for j := 1; j <= 10; j++ {
            processMeasurementArray := []s.ProcessMeasurement{}
            for k := 1; k <= 2; k++ {
                processMeasurementArray = append(processMeasurementArray,
                             s.ProcessMeasurement{
                                 TimeSlice: float64(k * 10),
                                 CpuUsed: float64(k * 5),
                                 MemUsed: float64(k * 6),
                              })
            }
            processMetricsMap["proc" + strconv.Itoa(j)] = processMeasurementArray
        }
        nodeData.ProcessMeasurementMap = processMetricsMap
        nodeMetricsMap["node" + strconv.Itoa(i)] = nodeData
     }

}



// CreateDummyProcessMetricsHistory - creates dummy history of submitted Process metrics and 
// stores it in processMetricsArray
// used by Golden testing
func CreateDummyProcessMetricsHistory(processMetricsArray *[]s.ProcessMetricsByName) {

    var timeSliceArray = []float64{10, 30, 20, 50, 43, 22, 12, 34, 22, 8,
                                   100, 43, 56, 99, 87, 56, 21, 5, 12, 84}

    rand.Seed(time.Now().UnixNano())

    for i := 1; i <= 20; i++ {

        processMetricsData := s.ProcessMetricsByName{
                 ProcessName: "proc" + strconv.Itoa(i),
                 MetricsData: s.ProcessMeasurement{
                                  TimeSlice: timeSliceArray[i-1],
                                  CpuUsed: float64(rand.Intn(100)),
                                  MemUsed: float64(rand.Intn(100)),
                              },
        }

        *processMetricsArray = append(*processMetricsArray, processMetricsData)

    }

}



